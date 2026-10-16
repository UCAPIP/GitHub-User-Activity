# GitHub Activity CLI

🇬🇧 [English](#english) | 🇷🇺 [Русский](#русский)

---

## English

A small command-line tool written in Go that fetches the recent public activity of any GitHub user and prints it in the terminal. It uses the public [GitHub Events API](https://docs.github.com/en/rest/activity/events) and only the Go standard library, so there are no external dependencies.

### Features

- Takes a GitHub username as a single argument
- Fetches the user's recent public events from the GitHub API
- Prints a short human-readable line for each supported event
- Clear error messages for a missing argument, an unknown user and API failures

### Supported events

| GitHub event | Output |
|---|---|
| `PushEvent` | `Pushed N commits to <owner>/<repo>` |
| `IssuesEvent` | `Opened a new issue in <owner>/<repo>` |
| `WatchEvent` | `Starred <owner>/<repo>` |

Other event types are currently ignored.

### Requirements

- Go 1.21 or newer (recommended)
- Internet access (the tool calls `api.github.com`)

### Installation and build

```bash
git clone <your-repo-url>
cd <your-repo-folder>

# if the project has no go.mod yet
go mod init github-activity

go build -o github-activity
```

On Windows the result is `github-activity.exe`.

### Usage

```bash
./github-activity <username>
```

Without building, you can run it directly:

```bash
go run . <username>
```

Example:

```text
$ ./github-activity octocat
Fetching activity for octocat
Pushed 3 commits to octocat/Hello-World
Opened a new issue in octocat/Spoon-Knife
Starred octocat/Hello-World
```

On Windows PowerShell:

```powershell
.\github-activity.exe octocat
```

### Errors and exit codes

The program exits with code `1` and prints a message to stderr in these cases:

| Situation | Message |
|---|---|
| Wrong number of arguments | `Error: Exactly 1 argument is required` + usage hint |
| User does not exist (HTTP 404) | `Error: User '<name>' not found` |
| Any other non-200 response | `Error: API returned <status>` |
| Network or JSON decoding problem | `Error: <details>` |

### How it works

1. Reads the username from the command-line arguments.
2. Sends `GET https://api.github.com/users/<username>/events`.
3. Checks the HTTP status code.
4. Decodes the JSON array of events into Go structs.
5. Loops over the events and prints a line for each supported type.

### Known limitations

- Only public events are available, and the API returns just a limited window of recent events.
- Unauthenticated requests are subject to GitHub's low rate limit (60 requests per hour per IP address).
- `IssuesEvent` is always printed as "Opened a new issue", even if the issue was closed or reopened.
- For `PushEvent` the commit count is taken from the payload `size` field, and a value of 0 is shown as 1.
- Events are not grouped, so several pushes to one repository are printed as separate lines.
- No request timeout is configured.

### Roadmap

- [ ] Check the `action` field of `IssuesEvent` (opened / closed / reopened)
- [ ] Support more event types (pull requests, forks, comments, branch/tag creation)
- [ ] Group events by repository and event type
- [ ] HTTP client with timeout, `User-Agent` header and optional `GITHUB_TOKEN`
- [ ] Friendly message when the user has no recent activity
- [ ] Flags such as `--limit` and `--type`
- [ ] Unit tests with `net/http/httptest`

---

## Русский

Небольшая консольная утилита на Go, которая получает недавнюю публичную активность любого пользователя GitHub и выводит её в терминал. Использует публичный [GitHub Events API](https://docs.github.com/en/rest/activity/events) и только стандартную библиотеку Go, поэтому внешних зависимостей нет.

### Возможности

- Принимает имя пользователя GitHub единственным аргументом
- Загружает недавние публичные события пользователя через API GitHub
- Печатает короткую понятную строку для каждого поддерживаемого события
- Понятные сообщения об ошибках: нет аргумента, пользователь не найден, сбой API

### Поддерживаемые события

| Событие GitHub | Вывод |
|---|---|
| `PushEvent` | `Pushed N commits to <owner>/<repo>` |
| `IssuesEvent` | `Opened a new issue in <owner>/<repo>` |
| `WatchEvent` | `Starred <owner>/<repo>` |

Остальные типы событий пока игнорируются.

### Требования

- Go 1.21 или новее (рекомендуется)
- Доступ в интернет (утилита обращается к `api.github.com`)

### Установка и сборка

```bash
git clone <ссылка-на-репозиторий>
cd <папка-репозитория>

# если в проекте ещё нет go.mod
go mod init github-activity

go build -o github-activity
```

На Windows получится файл `github-activity.exe`.

### Использование

```bash
./github-activity <username>
```

Можно запустить без сборки:

```bash
go run . <username>
```

Пример:

```text
$ ./github-activity octocat
Fetching activity for octocat
Pushed 3 commits to octocat/Hello-World
Opened a new issue in octocat/Spoon-Knife
Starred octocat/Hello-World
```

В Windows PowerShell:

```powershell
.\github-activity.exe octocat
```

### Ошибки и коды завершения

Программа завершается с кодом `1` и печатает сообщение в stderr в таких случаях:

| Ситуация | Сообщение |
|---|---|
| Неверное количество аргументов | `Error: Exactly 1 argument is required` + подсказка по использованию |
| Пользователя не существует (HTTP 404) | `Error: User '<name>' not found` |
| Любой другой ответ кроме 200 | `Error: API returned <status>` |
| Проблема с сетью или разбором JSON | `Error: <подробности>` |

### Как это работает

1. Читает имя пользователя из аргументов командной строки.
2. Отправляет `GET https://api.github.com/users/<username>/events`.
3. Проверяет код статуса HTTP.
4. Раскладывает JSON-массив событий по структурам Go.
5. Проходит по событиям и печатает строку для каждого поддерживаемого типа.

### Известные ограничения

- Доступны только публичные события, и API отдаёт ограниченное окно недавних событий.
- Для неавторизованных запросов действует низкий лимит GitHub (60 запросов в час с одного IP-адреса).
- `IssuesEvent` всегда выводится как «Opened a new issue», даже если задачу закрыли или переоткрыли.
- Для `PushEvent` число коммитов берётся из поля `size`, а значение 0 показывается как 1.
- События не группируются: несколько пушей в один репозиторий выводятся отдельными строками.
- Таймаут HTTP-запроса не задан.

### Планы

- [ ] Учитывать поле `action` у `IssuesEvent` (opened / closed / reopened)
- [ ] Поддержать больше типов событий (pull request, fork, комментарии, создание веток и тегов)
- [ ] Группировать события по репозиторию и типу
- [ ] HTTP-клиент с таймаутом, заголовком `User-Agent` и необязательным `GITHUB_TOKEN`
- [ ] Дружелюбное сообщение, если у пользователя нет недавней активности
- [ ] Флаги вроде `--limit` и `--type`
- [ ] Юнит-тесты на `net/http/httptest`
