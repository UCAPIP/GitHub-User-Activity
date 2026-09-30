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

// GitHubEvent описывает один элемент из массива событий
type GitHubEvent struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Actor     Actor     `json:"actor"`
	Repo      Repo      `json:"repo"`
	Payload   Payload   `json:"payload"`
	Public    bool      `json:"public"`
	CreatedAt time.Time `json:"created_at"`
}

// Actor описывает пользователя, совершившего действие
type Actor struct {
	ID           int    `json:"id"`
	Login        string `json:"login"`
	DisplayLogin string `json:"display_login"`
	GravatarID   string `json:"gravatar_id"`
	URL          string `json:"url"`
	AvatarURL    string `json:"avatar_url"`
}

// Repo описывает репозиторий, в котором произошло событие
type Repo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Payload содержит специфичные данные для каждого типа события.
// Мы объединили поля из WatchEvent и PushEvent. Те, которых нет в конкретном событии, просто останутся дефолтными.
type Payload struct {
	Action       string `json:"action,omitempty"`        // Для WatchEvent ("started")
	RepositoryID int    `json:"repository_id,omitempty"` // Для PushEvent
	PushID       int64  `json:"push_id,omitempty"`       // Для PushEvent (используем int64, так как ID бывают большими)
	Ref          string `json:"ref,omitempty"`           // Для PushEvent ("refs/heads/master")
	Head         string `json:"head,omitempty"`          // Для PushEvent
	Before       string `json:"before,omitempty"`        // Для PushEvent
	Size         int    `json:"size"`                    //Для PushEvent
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Ошибка: нужно ровно 1 аргумент")
		fmt.Println("Использование: github-activity <username>")
		os.Exit(1)
	}
	username := os.Args[1]
	fmt.Println("Fetching activity for", username)

	resp, err := http.Get("https://api.github.com/users/" + username + "/events")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
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
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}

	var events []GitHubEvent
	err = json.Unmarshal(body, &events)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
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
