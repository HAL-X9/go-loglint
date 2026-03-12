# go-loglint

Линтер для Go, проверяющий сообщения в логах. Работает как отдельная утилита (CLI) и как плагин для golangci-lint.

## Установка в свой проект

### Вариант 1: Плагин для golangci-lint

Если используешь golangci-lint, подключи loglint как плагин.

**Шаг 1.** Создай `.custom-gcl.yml` в корне проекта:

```yaml
version: v2.11.3   # Должно совпадать с: golangci-lint version

plugins:
  - module: github.com/HAL-X9/go-loglint
    import: github.com/HAL-X9/go-loglint/plugin/module
    version: v1.0.6
```

**Шаг 2.** Собери кастомный golangci-lint:

```bash
golangci-lint custom
```

Появится бинарник `./custom-gcl`.

**Шаг 3.** Добавь в `.golangci.yml`:

```yaml
version: "2"
linters:
  enable:
    - loglint
  settings:
    custom:
      loglint:
        type: module
        description: Log message linter
```

**Шаг 4.** Запуск:

```bash
# Без конфига — все правила включены по умолчанию
./custom-gcl run

# С конфигом — через переменную окружения
LOGLINT_CONFIG_PATH=./loglint.yaml ./custom-gcl run
```

**Важно:** golangci-lint не передаёт `settings.config` в плагины. Используй `LOGLINT_CONFIG_PATH` для конфига.

---

### Вариант 2: Собрать из исходников

```bash
git clone https://github.com/HAL-X9/go-loglint.git
cd go-loglint
make build
```

Бинарник будет в `./bin/go-loglint`:

```bash
./bin/go-loglint ./...
```

**С конфигом:**

```bash
./bin/go-loglint -config=loglint.yaml ./...
```

---

## Конфигурация (loglint.yaml)

Создай `loglint.yaml` в корне проекта:

```yaml
version: 1.0

rules:
  - name: lowercase_start
    enable: true
  - name: english_only
    enable: true
  - name: no_special_chars
    enable: true
  - name: no_sensitive_data
    enable: true

loggers:
  include:
    - log
    - slog
```

| Правило | Описание |
|---------|----------|
| `lowercase_start` | Сообщения должны начинаться с маленькой буквы |
| `english_only` | Только английский язык в логах |
| `no_special_chars` | Только буквы, цифры и пробелы |
| `no_sensitive_data` | Поиск чувствительных данных (пароли, токены) |

---

## Примеры

**Неправильно:**
```go
log.Print("Starting server")
slog.Error("Ошибка подключения")
slog.Info("user password: " + password)
```

**Правильно:**
```go
log.Print("starting server")
slog.Error("connection failed")
slog.Info("user authenticated")
```
