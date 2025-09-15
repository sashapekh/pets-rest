# CI/CD Pipeline Documentation

## 🚀 Огляд

Цей проект використовує GitHub Actions для автоматизації тестування, збірки та деплойменту.

## 📋 Workflows

### 1. **CI Pipeline** (`.github/workflows/ci.yml`)

Запускається при:
- Push до `main` або `develop` гілок
- Pull Request до `main` або `develop` гілок

**Етапи:**

#### 🔍 **Lint & Format Check**
- Перевірка форматування коду (`gofmt`)
- Статичний аналіз коду (`golangci-lint`) 
- Перевірка `go vet`
- Кешування Go модулів

#### 🧪 **Tests**
- Запуск PostgreSQL та Redis сервісів
- Виконання міграцій бази даних
- Запуск тестів з coverage звітом
- Генерація HTML звіту покриття

#### 🔨 **Build**
- Збірка для Linux та Windows
- Оптимізовані бінарні файли з `ldflags`
- Завантаження артефактів

#### 🐳 **Docker**
- Збірка Docker образу
- Тестування образу
- Кешування для прискорення

#### 🔒 **Security Scan**
- Сканування на вразливості (`gosec`)
- SARIF звіти для GitHub Security

### 2. **Release Pipeline** (`.github/workflows/release.yml`)

Запускається при створенні тегу (`v*`)

**Можливості:**
- Збірка для множини платформ (Linux, Windows, macOS)
- ARM64 та AMD64 архітектури
- Публікація Docker образів в GitHub Container Registry
- Створення GitHub Release з артефактами
- Автоматичний changelog

### 3. **Dependabot** (`.github/dependabot.yml`)

Автоматичні оновлення:
- Go модулів (щотижня)
- GitHub Actions (щотижня)
- Docker базових образів (щотижня)

## 🛠️ Локальні команди

```bash
# Встановити інструменти CI
make ci-setup

# Запустити лінтер
make lint

# Виправити автоматично
make lint-fix

# Форматувати код
make format

# Тести з покриттям
make test-cover

# Перевірка перед commit
make check
```

## 📊 Покриття коду

- Генерується HTML звіт: `coverage.html`
- Завантажується як артефакт в GitHub Actions
- Мінімальне покриття: рекомендується 80%+

## 🏷️ Створення релізу

1. Створіть тег:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. GitHub автоматично:
   - Запустить тести
   - Зібере бінарні файли
   - Створить Docker образ
   - Опублікує релін

## 🐳 Docker образи

Образи публікуються в `ghcr.io/[username]/pets_search_rest`:
- `latest` - останній релін
- `v1.0.0` - конкретна версія

## 🔒 Security Features

- **Gosec** - статичний аналіз безпеки
- **SARIF upload** - інтеграція з GitHub Security
- **Dependabot** - автоматичні оновлення залежностей
- **Docker scan** - сканування образів

## ⚡ Оптимізації

- **Go modules cache** - прискорення збірки
- **Docker buildx cache** - швидші Docker збірки
- **Parallel jobs** - паралельне виконання етапів
- **Оптимізовані бінарні файли** - мінімальний розмір

## 🎯 Best Practices

1. **Завжди запускайте тести локально** перед push
2. **Використовуйте `make check`** перед commit
3. **Пишіть тести** для нового коду
4. **Дотримуйтесь** форматування коду
5. **Оновлюйте** залежності регулярно

## 🚨 Troubleshooting

### Помилки лінтингу
```bash
# Виправити автоматично
make lint-fix

# Ручне форматування
make format
```

### Тести не проходять
```bash
# Локальний запуск з деталями
go test -v ./...

# З race detection
go test -race ./...
```

### Docker збірка падає
```bash
# Локальна збірка
docker build -t pets-api .

# Перевірка образу
docker run --rm pets-api --help
```

## 📈 Моніторинг

GitHub надає метрики:
- Час виконання workflows
- Успішність збірок
- Покриття коду
- Security alerts
