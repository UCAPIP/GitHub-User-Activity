package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// GitHubEvent describes a single element in an array of events
type GitHubEvent struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Actor     Actor     `json:"actor"`
	Repo      Repo      `json:"repo"`
	Payload   Payload   `json:"payload"`
	Public    bool      `json:"public"`
	CreatedAt time.Time `json:"created_at"`
}

// “Actor” describes a user who has performed an action
type Actor struct {
	ID           int    `json:"id"`
	Login        string `json:"login"`
	DisplayLogin string `json:"display_login"`
	GravatarID   string `json:"gravatar_id"`
	URL          string `json:"url"`
	AvatarURL    string `json:"avatar_url"`
}

// “Repo” refers to the repository in which the event occurred
type Repo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// The payload contains event-specific data for each event type.
// I've combined the fields from WatchEvent and PushEvent. Any fields that aren't present in a specific event will simply default to their default values.
type Payload struct {
	Action       string `json:"action,omitempty"`        // For WatchEvent ("started")
	RepositoryID int    `json:"repository_id,omitempty"` // For PushEvent
	PushID       int64  `json:"push_id,omitempty"`       // For PushEvent (I use int64 because IDs can be large)
	Ref          string `json:"ref,omitempty"`           // For PushEvent ("refs/heads/master")
	Head         string `json:"head,omitempty"`          // For PushEvent
	Before       string `json:"before,omitempty"`        // For PushEvent
	Size         int    `json:"size"`                    // For PushEvent
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Error: Exactly 1 argument is required")
		fmt.Println("Usage: github-activity <username>")
		os.Exit(1)
	}
	username := os.Args[1]
	fmt.Println("Fetching activity for", username)

	resp, err := http.Get("https://api.github.com/users/" + username + "/events")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		fmt.Fprintf(os.Stderr, "Error: User '%s' not found\n", username)
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "Error: API returned %s\n", resp.Status)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var events []GitHubEvent
	err = json.Unmarshal(body, &events)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	for _, event := range events {
		switch event.Type {
		case "IssuesEvent":
			fmt.Println("Opened a new issue in " + event.Repo.Name)
		case "PushEvent":
			if event.Payload.Size == 0 {
				event.Payload.Size = 1
			}
			fmt.Println("Pushed " + strconv.Itoa(event.Payload.Size) + " commits to " + event.Repo.Name)
		case "WatchEvent":
			fmt.Println("Starred " + event.Repo.Name)
		default:
		}
	}
}
