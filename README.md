# GoMicroservicesCourse

<!-- Mahno9/3c66bcabd00aec91a86c0ba4c468f01c -->
![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/Mahno9/3c66bcabd00aec91a86c0ba4c468f01c/raw/coverage.json)

Этот репозиторий содержит проект из курса [Микросервисы, как в BigTech 2.0](https://olezhek28.courses/microservices) от [Олега Козырева](http://t.me/olezhek28go).

Для того чтобы вызывать команды из Taskfile, необходимо установить Taskfile CLI:

```bash
brew install go-task
```

## CI/CD

Проект использует GitHub Actions для непрерывной интеграции и доставки. Основные workflow:

- **CI** (`.github/workflows/ci.yml`) - проверяет код при каждом push и pull request
  - Линтинг кода
  - Проверка безопасности
  - Выполняется автоматическое извлечение версий из Taskfile.yml
