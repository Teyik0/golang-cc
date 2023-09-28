package utils

import (
	"encoding/csv"
	"fmt"
	"os"
)

func CreateCSVFile(repos []Repo) error {
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
