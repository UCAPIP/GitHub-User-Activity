# GitHub Activity CLI

Консольное приложение для просмотра последней активности пользователя GitHub.

## Использование

\`\`\`bash
github-activity <username>
\`\`\`

## Пример

\`\`\`bash
$ github-activity kamranahmedse
- Pushed 3 commits to kamranahmedse/developer-roadmap
- Opened a new issue in kamranahmedse/developer-roadmap
- Starred kamranahmedse/developer-roadmap
\`\`\`

## Обработка ошибок

- Пользователь не найден: `Error: User 'xxx' not found`
- Ошибка API: `Error: API returned 500 Internal Server Error`

## Технологии

- Go (стандартная библиотека)
- GitHub REST API
\`\`\`

---