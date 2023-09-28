package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Teyik0/golang-cc/utils"
	"github.com/joho/godotenv"
)

var wg sync.WaitGroup

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
	jsonData, err := utils.FetchAllRepoFromUser(user)
	if err != nil {
		log.Fatal("test", err)
	}

	// Unmarshal JSON data
	var repos []utils.Repo
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
	utils.CreateCSVFile(repos)

	// Clone all repos + fetch + pull
	for _, repo := range repos {
		wg.Add(1)
		go processRepo(repo, folderPath)
	}
	// Wait for all goroutines to finish
	wg.Wait()

	// Create a ZIP file
	if err := utils.ZipSource("repos", "repos.zip"); err != nil {
		log.Fatal(err)
	}
}

func processRepo(repo utils.Repo, folderPath string) {
	defer wg.Done()

	repoPath := folderPath + "/" + repo.Name
	utils.CloneRepo(repo, repoPath)
	utils.FetchRepository(repo, repoPath)
	utils.SwitchAndPull(repo, repoPath)
}
