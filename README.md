# Pets Search REST API

REST API для сервісу пошуку та оголошень про тварин.

## 🚀 Швидкий старт

### Передумови
- Docker та Docker Compose
- Go 1.22+ (для локальної розробки)

### Запуск з Docker Compose

1. Клонуйте репозиторій:
```bash
git clone <repository-url>
cd pets_search/rest
```

2. Скопіюйте файл конфігурації:
```bash
cp example.env .env
```

3. Запустіть всі сервіси:
```bash
docker-compose up -d
```

4. Перевірте роботу API:
```bash
curl http://localhost:8080/healthz
```

### Сервіси

Після запуску доступні наступні сервіси:

- **API**: http://localhost:8080
- **MinIO Console**: http://localhost:9090 (minioadmin/minioadmin)
- **MongoDB**: localhost:27017 (admin/password)
- **Redis**: localhost:6379

## 📂 Структура проекту

```
.
├── cmd/api/                 # Вхідна точка програми
├── internal/                # Внутрішня бізнес-логіка
│   ├── listings/           # Управління оголошеннями
│   ├── users/              # Аутентифікація користувачів
│   ├── storage/            # Інтеграція з S3/MinIO
│   ├── pdf/                # Генерація PDF
│   ├── qrcode/             # Генерація QR-кодів
│   ├── config/             # Конфігурація
│   └── database/           # Робота з базою даних
├── pkg/                    # Публічні пакети
│   ├── middleware/         # HTTP middleware
│   └── utils/              # Утиліти
├── web/                    # Веб-ресурси
│   ├── templates/          # HTML шаблони
│   └── static/             # Статичні файли
├── docs/                   # Документація
├── deployments/            # Конфігурації деплою
└── scripts/                # Скрипти
```

## 🔧 Розробка

### Локальний запуск без Docker

1. Запустіть залежності:
```bash
docker-compose up -d mongodb redis minio minio-init
```

2. Встановіть залежності Go:
```bash
go mod tidy
```

3. Запустіть API:
```bash
go run ./cmd/api
```

### Оновлення залежностей

```bash
go get -u ./...
go mod tidy
```

## 📋 API Endpoints

### Аутентифікація
- `POST /auth/magic-link` - Відправка magic link
- `POST /auth/magic-link/verify` - Верифікація magic link

### Оголошення
- `GET /api/v1/listings` - Список оголошень
- `POST /api/v1/listings` - Створення оголошення
- `GET /api/v1/listings/{id}` - Отримання оголошення
- `PUT /api/v1/listings/{id}` - Оновлення оголошення
- `DELETE /api/v1/listings/{id}` - Видалення оголошення
- `POST /api/v1/listings/{id}/images` - Завантаження фото
- `POST /api/v1/listings/{id}/generate-pdf` - Генерація PDF

### Публічні сторінки
- `GET /p/{slug}` - Публічна сторінка оголошення

### Система
- `GET /healthz` - Перевірка здоров'я сервісу

## 🗃️ База даних

Проект використовує MongoDB з наступними колекціями:

- **users** - Користувачі системи
- **listings** - Оголошення про тварин
- **events** - Аналітичні події

## 📦 Залежності

- [Fiber v2](https://github.com/gofiber/fiber) - HTTP framework
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver) - MongoDB драйвер
- [Redis Go Client](https://github.com/go-redis/redis) - Redis клієнт
- [MinIO Go SDK](https://github.com/minio/minio-go) - S3-сумісний клієнт

## 🔒 Безпека

- Всі паролі та ключі мають бути змінені у production
- JWT токени для автентифікації
- Валідація всіх вхідних даних
- Rate limiting для API endpoints

## 📝 Ліцензія

[Вкажіть ліцензію]
