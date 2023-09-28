package main

import (
	"archive/zip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

type Repo struct {
	ID                       int        `json:"id"`
	NodeID                   string     `json:"node_id"`
	Name                     string     `json:"name"`
	FullName                 string     `json:"full_name"`
	Private                  bool       `json:"private"`
	Owner                    Owner      `json:"owner"`
	HTML_URL                 string     `json:"html_url"`
	Description              string     `json:"description"`
	Fork                     bool       `json:"fork"`
	Url                      string     `json:"url"`
	ForksURL                 string     `json:"forks_url"`
	KeysURL                  string     `json:"keys_url"`
	CollaboratorsURL         string     `json:"collaborators_url"`
	TeamsURL                 string     `json:"teams_url"`
	HooksURL                 string     `json:"hooks_url"`
	IssueEventsUrl           string     `json:"issue_events_url"`
	EventsURL                string     `json:"events_url"`
	AssigneesURL             string     `json:"assignees_url"`
	BranchesURL              string     `json:"branches_url"`
	TagsURL                  string     `json:"tags_url"`
	BlobsURL                 string     `json:"blobs_url"`
	GitTagsURL               string     `json:"git_tags_url"`
	GitRefsURL               string     `json:"git_refs_url"`
	TreesURL                 string     `json:"trees_url"`
	StatusesURL              string     `json:"statuses_url"`
	LanguagesURL             string     `json:"languages_url"`
	StargazersURL            string     `json:"stargazers_url"`
	ContributorsURL          string     `json:"contributors_url"`
	SubscribersURL           string     `json:"subscribers_url"`
	SubscriptionURL          string     `json:"subscription_url"`
	CommitsURL               string     `json:"commits_url"`
	GitCommitsURL            string     `json:"git_commits_url"`
	CommentsURL              string     `json:"comments_url"`
	IssueCommentURL          string     `json:"issue_comment_url"`
	ContentsURL              string     `json:"contents_url"`
	CompareURL               string     `json:"compare_url"`
	MergesURL                string     `json:"merges_url"`
	ArchiveURL               string     `json:"archive_url"`
	DownloadsURL             string     `json:"downloads_url"`
	IssuesURL                string     `json:"issues_url"`
	PullsURL                 string     `json:"pulls_url"`
	MilestonesURL            string     `json:"milestones_url"`
	NotificationsURL         string     `json:"notifications_url"`
	LabelsURL                string     `json:"labels_url"`
	ReleasesURL              string     `json:"releases_url"`
	DeploymentsURL           string     `json:"deployments_url"`
	CreatedAt                string     `json:"created_at"`
	UpdatedAt                string     `json:"updated_at"`
	PushedAt                 string     `json:"pushed_at"`
	GitURL                   string     `json:"git_url"`
	SshURL                   string     `json:"ssh_url"`
	CloneURL                 string     `json:"clone_url"`
	SvnURL                   string     `json:"svn_url"`
	Homepage                 string     `json:"homepage"`
	Size                     int        `json:"size"`
	StargazersCount          int        `json:"stargazers_count"`
	WatchersCount            int        `json:"watchers_count"`
	Language                 string     `json:"language"`
	HasIssues                bool       `json:"has_issues"`
	HasProjects              bool       `json:"has_projects"`
	HasDownloads             bool       `json:"has_downloads"`
	HasWiki                  bool       `json:"has_wiki"`
	HasPages                 bool       `json:"has_pages"`
	ForksCount               int        `json:"forks_count"`
	MirrorURL                string     `json:"mirror_url"`
	Archived                 bool       `json:"archived"`
	Disabled                 bool       `json:"disabled"`
	OpenIssuesCount          int        `json:"open_issues_count"`
	License                  License    `json:"license"`
	AllowForking             bool       `json:"allow_forking"`
	IsTemplate               bool       `json:"is_template"`
	WebCommitSignoffRequired bool       `json:"web_commit_signoff_required"`
	Permission               Permission `json:"permission"`
}

type License struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	SpdxID string `json:"spdx_id"`
	URL    string `json:"url"`
	NodeID string `json:"node_id"`
}

type Permission struct {
	Admin    bool `json:"admin"`
	Push     bool `json:"push"`
	Pull     bool `json:"pull"`
	Maintain bool `json:"maintain"`
	Triage   bool `json:"triage"`
}

type Owner struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HtmlURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	// Fetch all repos from user
	user, ok := os.LookupEnv("GITHUB_USER")
	if !ok {
		log.Fatal("GITHUB_USER is not set in .env file")
	}
	jsonData, err := fetchAllRepoFromUser(user)
	if err != nil {
		log.Fatal("test", err)
	}

	// Unmarshal JSON data
	var repos []Repo
	if err := json.Unmarshal([]byte(jsonData), &repos); err != nil {
		log.Fatal("Error decoding JSON:", err)
	}

	// Create/Verifie the folder for the cloned repos and the csv
	folderPath := "repos"
	_, err1 := os.Stat(folderPath)
	if err1 != nil {
		fmt.Println("Creating folder repos...")
		os.Mkdir(folderPath, os.ModePerm)
	}
	err2 := os.Chmod(folderPath, 0777)
	if err2 != nil {
		fmt.Println(err2)
	}

	// Create a CSV file
	createCSVFile(repos)

	// Clone all repos + fetch + pull
	for _, repo := range repos {
		repoPath := folderPath + "/" + repo.Name
		cloneRepo(repo, repoPath)
		fetchRepository(repo, repoPath)
		switchAndPull(repo, repoPath)
	}

	// Create a ZIP file
	if err := zipSource("repos", "repos.zip"); err != nil {
		log.Fatal(err)
	}
}

func zipSource(source, target string) error {
	// 1. Create a ZIP file and zip.Writer
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// 2. Go through all the files of the source
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}

func _getLastCommitBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func pullRepository(repo Repo, repoPath string) {
	fmt.Printf("Pulling to the latest changes for repository: %s ...\n\n", repo.Name)
	cmd := exec.Command("git", "-C", repoPath, "pull")
	if err := cmd.Run(); err != nil {
		log.Printf("Error pulling repository %s: %s\n\n", repo.Name, err)
	}
}

func switchAndPull(repo Repo, repoPath string) {
	lastCommitBranch, err := _getLastCommitBranch(repoPath)
	if err != nil {
		log.Printf("Error getting last commit branch for repository %s: %s\n", repo.Name, err)
		return
	}
	fmt.Println("Last commit branch : ", lastCommitBranch)

	if lastCommitBranch != "main" {
		cmd := exec.Command("git", "-C", repoPath, "checkout", lastCommitBranch)
		if err := cmd.Run(); err != nil {
			log.Printf("Error switching to branch %s for repository %s: %s", lastCommitBranch, repo.Name, err)
		}
	}

	// Pull the latest changes
	pullRepository(repo, repoPath)
}

func fetchRepository(repo Repo, repoPath string) {
	fmt.Printf("| Fetching latest changes for repository: %s ...\n", repo.Name)
	cmd := exec.Command("git", "fetch", repoPath)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func cloneRepo(repo Repo, repoPath string) {
	_, err := os.Stat(repoPath)
	if err != nil {
		fmt.Printf("| Cloning repository : %s ...\n", repo.FullName)
		os.Mkdir(repoPath, os.ModePerm)
		cmd := exec.Command("git", "clone", repo.CloneURL, repoPath)
		if err := cmd.Run(); err != nil {
			log.Fatal("Failed to clone repo", err)
		}
	} else {
		fmt.Printf("| Directory %s already exists \n", repoPath)
	}
}

func fetchAllRepoFromUser(username string) (string, error) {
	const githubAPIEndpoint = "https://api.github.com/users/%s/repos"

	tokenUrl, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		log.Fatal("GITHUB_TOKEN is not set in .env file")
	}
	var bearer = "Bearer " + tokenUrl

	url := fmt.Sprintf(githubAPIEndpoint, username)

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("test1", err)
	}

	return string(body), err
}

func createCSVFile(repos []Repo) error {
	// Create a CSV file
	file, err := os.Create("repos/repos.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "NodeID", "Name", "FullName", "Private", "Owner", "HTML_URL", "Description", "Fork", "Url", "ForksURL", "KeysURL", "CollaboratorsURL", "TeamsURL", "HooksURL", "IssueEventsUrl", "EventsURL", "AssigneesURL", "BranchesURL", "TagsURL", "BlobsURL", "GitTagsURL", "GitRefsURL", "TreesURL", "StatusesURL", "LanguagesURL", "StargazersURL", "ContributorsURL", "SubscribersURL", "SubscriptionURL", "CommitsURL", "GitCommitsURL", "CommentsURL", "IssueCommentURL", "ContentsURL", "CompareURL", "MergesURL", "ArchiveURL", "DownloadsURL", "IssuesURL", "PullsURL", "MilestonesURL", "NotificationsURL", "LabelsURL", "ReleasesURL", "DeploymentsURL", "CreatedAt", "UpdatedAt", "PushedAt", "GitURL", "SshURL", "CloneURL", "SvnURL", "Homepage", "Size", "StargazersCount", "WatchersCount", "Language", "HasIssues", "HasProjects", "HasDownloads", "HasWiki", "HasPages", "ForksCount", "MirrorURL", "Archived", "Disabled", "OpenIssuesCount", "License", "AllowForking", "IsTemplate", "WebCommitSignoffRequired", "Permission"}
	if err := writer.Write(header); err != nil {
		panic(err)
	}

	// Write data
	for _, repo := range repos {
		record := []string{
			// Convert each field to a string
			fmt.Sprint(repo.ID),
			repo.NodeID,
			repo.Name,
			repo.FullName,
			fmt.Sprint(repo.Private),
			fmt.Sprint(repo.Owner),
			repo.HTML_URL,
			repo.Description,
			fmt.Sprint(repo.Fork),
			repo.Url,
			repo.ForksURL,
			repo.KeysURL,
			repo.CollaboratorsURL,
			repo.TeamsURL,
			repo.HooksURL,
			repo.IssueEventsUrl,
			repo.EventsURL,
			repo.AssigneesURL,
			repo.BranchesURL,
			repo.TagsURL,
			repo.BlobsURL,
			repo.GitTagsURL,
			repo.GitRefsURL,
			repo.TreesURL,
			repo.StatusesURL,
			repo.LanguagesURL,
			repo.StargazersURL,
			repo.ContributorsURL,
			repo.SubscribersURL,
			repo.SubscriptionURL,
			repo.CommitsURL,
			repo.GitCommitsURL,
			repo.CommentsURL,
			repo.IssueCommentURL,
			repo.ContentsURL,
			repo.CompareURL,
			repo.MergesURL,
			repo.ArchiveURL,
			repo.DownloadsURL,
			repo.IssuesURL,
			repo.PullsURL,
			repo.MilestonesURL,
			repo.NotificationsURL,
			repo.LabelsURL,
			repo.ReleasesURL,
			repo.DeploymentsURL,
			repo.CreatedAt,
			repo.UpdatedAt,
			repo.PushedAt,
			repo.GitURL,
			repo.SshURL,
			repo.CloneURL,
			repo.SvnURL,
			repo.Homepage,
			fmt.Sprint(repo.Size),
			fmt.Sprint(repo.StargazersCount),
			fmt.Sprint(repo.WatchersCount),
			repo.Language,
			fmt.Sprint(repo.HasIssues),
			fmt.Sprint(repo.HasProjects),
			fmt.Sprint(repo.HasDownloads),
			fmt.Sprint(repo.HasWiki),
			fmt.Sprint(repo.HasPages),
			fmt.Sprint(repo.ForksCount),
			repo.MirrorURL,
			fmt.Sprint(repo.Archived),
			fmt.Sprint(repo.Disabled),
			fmt.Sprint(repo.OpenIssuesCount),
			repo.License.Name,
			fmt.Sprint(repo.AllowForking),
			fmt.Sprint(repo.IsTemplate),
			fmt.Sprint(repo.WebCommitSignoffRequired),
			fmt.Sprint(repo.Permission),
			// Add more fields as needed
		}
		if err := writer.Write(record); err != nil {
			panic(err)
		}
	}

	return nil
}
