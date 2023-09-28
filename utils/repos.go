package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func _getLastCommitBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func PullRepository(repo Repo, repoPath string) {
	fmt.Printf("Pulling to the latest changes for repository: %s ...\n\n", repo.Name)
	cmd := exec.Command("git", "-C", repoPath, "pull")
	if err := cmd.Run(); err != nil {
		log.Printf("Error pulling repository %s: %s\n\n", repo.Name, err)
	}
}

func SwitchAndPull(repo Repo, repoPath string) {
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
	PullRepository(repo, repoPath)
}

func FetchRepository(repo Repo, repoPath string) {
	fmt.Printf("| Fetching latest changes for repository: %s ...\n", repo.Name)
	cmd := exec.Command("git", "fetch", repoPath)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func CloneRepo(repo Repo, repoPath string) {
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

func FetchAllRepoFromUser(username string) (string, error) {
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
