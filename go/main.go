package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	baseUrl                       = "https://api.github.com/users/%s/events"
	prEvent                       = "PullRequestEvent"
	pushEvent                     = "PushEvent"
	prReviewCommentEvent          = "PulllRequestReviewCommentEvent"
	issueCommentEvent             = "IssueCommentEvent"
	pullRequestReviewCommentEvent = "PullRequestReviewCommentEvent"
	watchEvent                    = "WatchEvent"
)

type Event struct {
	// Id   string `json:"id"`
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	}
	// Payload struct {
	// 	Commits []struct {
	// 		Message string `json:"message"`
	// 	}
	// }
}

func main() {
	if len(os.Args) == 1 {
		help()
		os.Exit(0)
	}
	username := os.Args[1]
	url := fmt.Sprintf(baseUrl, username)
	// fmt.Printf("url: %s\n", url)
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("error creating a new HTTP request: %w", err))
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	// fmt.Println("Requesting user events to GitHub API...")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(fmt.Errorf("error fetching data from GitHub API: %w", err))
	}
	defer resp.Body.Close()
	// fmt.Printf("Status Code: %d\n", resp.StatusCode)
	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		log.Fatal(fmt.Errorf("error decoding HTTP payload: %w", err))
	}
	commitMap := make(map[string]int)
	commentMap := make(map[string]int)
	prCommentMap := make(map[string]int)
	watchMap := make(map[string]int)
	createPrMap := make(map[string]int)
	for _, ev := range events {
		// fmt.Printf("%+v\n", ev)
		switch ev.Type {
		case pushEvent:
			commitMap[ev.Repo.Name]++
		case issueCommentEvent:
			commentMap[ev.Repo.Name]++
		case pullRequestReviewCommentEvent:
			prCommentMap[ev.Repo.Name]++
		case watchEvent:
			watchMap[ev.Repo.Name]++
		case prEvent:
			createPrMap[ev.Repo.Name]++
		default:
		}
	}
	for repo, count := range commitMap {
		fmt.Printf("- Pushed %d commits to %s\n", count, repo)
	}
	for repo, count := range commentMap {
		fmt.Printf("- Left %d comments on a issue in %s\n", count, repo)
	}
	for repo, count := range prCommentMap {
		fmt.Printf("- Left %d comments in %s\n", count, repo)
	}
	for repo := range watchMap {
		fmt.Printf("- Started watching %s\n", repo)
	}
	for repo, count := range createPrMap {
		fmt.Printf("- Opened new %d PRs in %s\n", count, repo)
	}
	os.Exit(0)
}

func help() {
	fmt.Println("USAGE: githubactivity <username>")
}
