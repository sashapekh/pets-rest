# Інструкції налаштування

## 🔧 Передумови

### Системні вимоги
- **Docker** та **Docker Compose** v2.0+
- **Go** 1.22+ (для локальної розробки)
- **Git** для клонування репозиторію

### Перевірка налаштування Docker в WSL2
```bash
# Перевірити чи працює Docker
docker --version
docker compose --version

# Якщо команди не працюють:
# 1. Відкрийте Docker Desktop
# 2. Settings → Resources → WSL Integration
# 3. Увімкніть інтеграцію для вашого дистрибутива
```

## 🚀 Швидкий запуск

### 1. Клонування та підготовка
```bash
cd /home/sashapekh/projects/pets_search/rest

# Створити файл конфігурації
cp example.env .env

# Редагувати конфігурацію (опціонально)
nano .env
```

### 2. Запуск всіх сервісів
```bash
# Через Docker Compose
docker compose up -d

# Або через Makefile
make docker-up

# Перевірка статусу
docker compose ps
```

### 3. Перевірка роботи
```bash
# API health check
curl http://localhost:8080/healthz

# Або відкрити в браузері
open http://localhost:8080/healthz
```

## 🐳 Docker сервіси

### Порти та доступ
| Сервіс | Порт | Доступ | Логін/Пароль |
|--------|------|--------|--------------|
| API | 8080 | http://localhost:8080 | - |
| MongoDB | 27017 | localhost:27017 | admin/password |
| Redis | 6379 | localhost:6379 | - |
| MinIO API | 9001 | http://localhost:9001 | minioadmin/minioadmin |
| MinIO Console | 9091 | http://localhost:9091 | minioadmin/minioadmin |

### Управління сервісами
```bash
# Запуск
docker compose up -d

# Зупинка
docker compose down

# Перезапуск
docker compose restart

# Логи всіх сервісів
docker compose logs -f

# Логи конкретного сервісу
docker compose logs -f api
docker compose logs -f mongodb
```

## 🛠️ Локальна розробка

### Запуск без Docker (рекомендовано для розробки)

1. **Запустити тільки залежності:**
```bash
docker compose up -d mongodb redis minio minio-init
```

2. **Встановити Go залежності:**
```bash
go mod download
go mod tidy
```

3. **Запустити API локально:**
```bash
# Через Makefile
make run

# Або напряму
go run ./cmd/api

# Або з hot reload (потрібен air)
go install github.com/cosmtrek/air@latest
air
```

### Корисні команди для розробки
```bash
# Збірка проекту
make build

# Запуск тестів
make test

# Очищення
make clean

# Оновлення залежностей
make deps
```

## 📊 MongoDB налаштування

### Підключення до бази
```bash
# Через Docker
docker compose exec mongodb mongosh -u admin -p password

# Локально (якщо встановлений mongosh)
mongosh "mongodb://admin:password@localhost:27017/pets_search?authSource=admin"
```

### Початкові дані
База автоматично ініціалізується з схемою через `scripts/mongo-init.js`:
- Створюються колекції: `users`, `listings`, `events`
- Налаштовуються індекси для оптимізації
- Додається валідація схем

### Скидання бази даних
```bash
# УВАГА: Видаляє всі дані!
make db-reset

# Або вручну
docker compose down mongodb
docker volume rm pets_search_rest_mongodb_data
docker compose up -d mongodb
```

## 🗂️ MinIO налаштування

### Доступ до консолі
- URL: http://localhost:9091
- Логін: `minioadmin`
- Пароль: `minioadmin`

### Bucket автоматично створюється
Контейнер `minio-init` автоматично:
- Створює bucket `pets-photos`
- Налаштовує публічний доступ
- Готовий до завантаження файлів

### Тестування завантаження
```bash
# Через MinIO Client (mc)
docker compose exec minio-init mc ls minio/pets-photos
```

## ⚙️ Конфігурація (.env)

### Основні параметри
```env
# Сервер
ENV=development
PORT=8080

# База даних
MONGODB_URI=mongodb://admin:password@localhost:27017/pets_search?authSource=admin

# Redis
REDIS_URL=redis://localhost:6379

# MinIO
MINIO_ENDPOINT=localhost:9001
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_USE_SSL=false
MINIO_BUCKET=pets-photos
```

### Налаштування для production
```env
# Змініть ці значення для production:
JWT_SECRET=your-super-secret-jwt-key-here
MINIO_ACCESS_KEY=your-secure-access-key
MINIO_SECRET_KEY=your-secure-secret-key

# Email конфігурація
SMTP_HOST=smtp.gmail.com
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

## 🔍 Troubleshooting

### Проблема: Docker команди не працюють
```bash
# Рішення
sudo service docker start
# або перезапустити Docker Desktop
```

### Проблема: Порт вже зайнятий
```bash
# Знайти процес на порту
sudo lsof -i :8080

# Зупинити процес
sudo kill -9 <PID>
```

### Проблема: MongoDB не запускається
```bash
# Перевірити логи
docker compose logs mongodb

# Очистити volume та перезапустити
docker compose down
docker volume rm pets_search_rest_mongodb_data
docker compose up -d mongodb
```

### Проблема: MinIO bucket не створюється
```bash
# Перезапустити ініціалізацію
docker compose up minio-init
```

## 📝 Логи та моніторинг

### Перегляд логів
```bash
# Всі сервіси
docker compose logs -f

# Конкретний сервіс
docker compose logs -f api

# Останні N рядків
docker compose logs --tail=50 api
```

### Моніторинг ресурсів
```bash
# Статистика Docker
docker stats

# Використання місця
docker system df
```

## 🔄 Оновлення проекту

### Оновлення коду
```bash
git pull origin main
docker compose build --no-cache
docker compose up -d
```

### Оновлення залежностей Go
```bash
go get -u ./...
go mod tidy
```
