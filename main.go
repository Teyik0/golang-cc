package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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
	jsonData, err := fetchAllRepoFromUser("MTthoas")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(jsonData)

	// Unmarshal JSON data
	var repos []Repo
	if err := json.Unmarshal([]byte(jsonData), &repos); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Print out the results
	for _, repo := range repos {
		fmt.Printf("Repository Name: %s\n", repo.FullName)
		fmt.Printf("Owner Login: %s\n", repo.Owner.Login)
	}
}

func fetchAllRepoFromUser(username string) (string, error) {
	const githubAPIEndpoint = "https://api.github.com/users/%s/repos"

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
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
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body), err
}
