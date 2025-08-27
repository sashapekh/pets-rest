# Інструкція по налаштуванню проекту

## ✅ Що вже створено

### 📂 Структура проекту
```
pets_search/rest/
├── cmd/api/                 # Вхідна точка API сервера
├── internal/                # Внутрішня бізнес-логіка
│   ├── config/             # Конфігурація (готово)
│   ├── database/           # Робота з БД
│   ├── listings/           # Оголошення
│   ├── users/              # Користувачі
│   ├── storage/            # S3/MinIO
│   ├── pdf/                # PDF генерація
│   └── qrcode/             # QR коди
├── pkg/                    # Публічні пакети
├── web/                    # Веб ресурси
├── docker-compose.yml      # ✅ Готово
├── Dockerfile             # ✅ Готово
├── Makefile               # ✅ Готово
└── README.md              # ✅ Готово
```

### 🐳 Docker Compose сервіси
- **API** (Golang + Fiber v2) - порт 8080
- **MongoDB** - порт 27017 (admin/password)
- **Redis** - порт 6379
- **MinIO** - порти 9000 (API), 9090 (Console)

### 📦 Залежності Go
- ✅ Fiber v3 - HTTP framework
- ✅ MongoDB Driver - для роботи з базою
- ✅ Redis Client - для кешування  
- ✅ MinIO SDK - для файлів
- ✅ GoDotEnv - для конфігурації

## 🚀 Наступні кроки

### 1. Налаштування Docker
Якщо у вас WSL2, переконайтеся що Docker Desktop інтегрований з WSL:
- Відкрийте Docker Desktop
- Settings → Resources → WSL Integration
- Увімкніть інтеграцію для вашого дистрибутива

### 2. Запуск проекту
```bash
# Перейти в папку проекту
cd /home/sashapekh/projects/pets_search/rest

# Створити .env файл
cp example.env .env

# Запустити всі сервіси
docker compose up -d

# Або використати Makefile
make docker-up

# Перевірити роботу
curl http://localhost:8080/healthz
```

### 3. Доступ до сервісів
- **API**: http://localhost:8080
- **MinIO Console**: http://localhost:9091 (minioadmin/minioadmin)  
- **MongoDB**: mongodb://admin:password@localhost:27017/pets_search

### 4. Розробка
```bash
# Запуск тільки залежностей (без API)
docker compose up -d mongodb redis minio minio-init

# Локальний запуск API
make run
# або
go run ./cmd/api
```

## ❓ Що додати далі

1. **Redis** - вже включений в docker-compose, готовий до використання
2. **Nginx** - якщо потрібен reverse proxy
3. **Prometheus/Grafana** - для моніторингу
4. **Elasticsearch** - для пошуку оголошень

**Поточна конфігурація:** Golang + MongoDB + Redis + MinIO

Чи потрібно щось додати або змінити?
