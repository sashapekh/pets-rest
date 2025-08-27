# Docker сервіси

## 🐳 Огляд контейнерів

Проект використовує Docker Compose для управління наступними сервісами:

```
pets-network
├── api (Golang REST API)
├── mongodb (База даних)
├── redis (Кеш та сесії)
├── minio (Файлове сховище)
└── minio-init (Ініціалізація MinIO)
```

## 🔧 Сервіс API

### Конфігурація
```yaml
api:
  build:
    context: .
    dockerfile: Dockerfile
  ports:
    - "8080:8080"
  environment:
    - ENV=development
    - PORT=8080
    - MONGODB_URI=mongodb://mongodb:27017/pets_search
    - REDIS_URL=redis://redis:6379
    - MINIO_ENDPOINT=minio:9000
  depends_on:
    - mongodb
    - redis
    - minio
```

### Особливості
- **Порт**: 8080
- **Volume**: Весь проект монтується для розробки
- **Restart policy**: unless-stopped
- **Залежності**: Чекає запуску всіх інших сервісів

### Dockerfile етапи
```dockerfile
# Build stage - компіляція Go коду
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/api

# Runtime stage - мінімальний образ
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

## 🗄️ Сервіс MongoDB

### Конфігурація
```yaml
mongodb:
  image: mongo:7.0
  container_name: pets_mongodb
  ports:
    - "27017:27017"
  environment:
    - MONGO_INITDB_ROOT_USERNAME=admin
    - MONGO_INITDB_ROOT_PASSWORD=password
    - MONGO_INITDB_DATABASE=pets_search
  volumes:
    - mongodb_data:/data/db
    - ./scripts/mongo-init.js:/docker-entrypoint-initdb.d/mongo-init.js:ro
```

### Особливості
- **Образ**: mongo:7.0 (офіційний)
- **Порт**: 27017
- **Автентифікація**: admin/password
- **Ініціалізація**: Автоматично виконується скрипт `mongo-init.js`
- **Persistence**: Дані зберігаються у volume `mongodb_data`

### Підключення
```bash
# З контейнера
docker compose exec mongodb mongosh -u admin -p password

# Ззовні
mongosh "mongodb://admin:password@localhost:27017/pets_search?authSource=admin"
```

### Ініціалізаційний скрипт
Скрипт `scripts/mongo-init.js`:
- Створює колекції з валідацією схем
- Налаштовує індекси для оптимізації
- Готує базу до роботи

## ⚡ Сервіс Redis

### Конфігурація
```yaml
redis:
  image: redis:7.2-alpine
  container_name: pets_redis
  ports:
    - "6379:6379"
  command: redis-server --appendonly yes
  volumes:
    - redis_data:/data
```

### Особливості
- **Образ**: redis:7.2-alpine (мінімальний)
- **Порт**: 6379
- **Persistence**: AOF (Append Only File) для збереження даних
- **Volume**: redis_data для збереження між перезапусками

### Використання
```bash
# Підключення через CLI
docker compose exec redis redis-cli

# Перевірка роботи
docker compose exec redis redis-cli ping
```

### Приклади команд
```redis
# Встановити ключ
SET user:session:abc123 "user_data"

# Отримати ключ
GET user:session:abc123

# Встановити з TTL (1 година)
SETEX temp:magic_link:xyz789 3600 "email@example.com"
```

## 📁 Сервіс MinIO

### Конфігурація
```yaml
minio:
  image: minio/minio:latest
  container_name: pets_minio
  ports:
    - "9000:9000"  # API
    - "9090:9090"  # Console
  environment:
    - MINIO_ROOT_USER=minioadmin
    - MINIO_ROOT_PASSWORD=minioadmin
  command: server /data --console-address ":9090"
  volumes:
    - minio_data:/data
```

### Особливості
- **API порт**: 9000 (S3-сумісний API)
- **Console порт**: 9090 (веб-інтерфейс)
- **Credentials**: minioadmin/minioadmin
- **Persistence**: minio_data volume

### Доступ
- **API**: http://localhost:9000
- **Console**: http://localhost:9090
- **Login**: minioadmin / minioadmin

### MinIO Client (MC)
```bash
# Налаштування alias
mc alias set local http://localhost:9000 minioadmin minioadmin

# Список buckets
mc ls local/

# Завантажити файл
mc cp photo.jpg local/pets-photos/
```

## 🔧 Сервіс MinIO Init

### Конфігурація
```yaml
minio-init:
  image: minio/mc:latest
  container_name: pets_minio_init
  depends_on:
    - minio
  entrypoint: >
    /bin/sh -c "
    sleep 10;
    /usr/bin/mc alias set minio http://minio:9000 minioadmin minioadmin;
    /usr/bin/mc mb minio/pets-photos --ignore-existing;
    /usr/bin/mc policy set public minio/pets-photos;
    exit 0;
    "
```

### Призначення
Одноразовий контейнер для ініціалізації MinIO:
1. Чекає запуску MinIO (10 секунд)
2. Налаштовує підключення
3. Створює bucket `pets-photos`
4. Встановлює публічну політику доступу
5. Завершує роботу

## 🌐 Мережа

### Конфігурація
```yaml
networks:
  pets-network:
    driver: bridge
```

### Внутрішні hostname
Усі контейнери можуть звертатися один до одного за іменами:
- `mongodb` - база даних
- `redis` - кеш
- `minio` - файлове сховище
- `api` - REST API

## 💾 Volumes

### Постійне зберігання
```yaml
volumes:
  mongodb_data:    # Дані MongoDB
  redis_data:      # Дані Redis
  minio_data:      # Файли MinIO
```

### Управління volumes
```bash
# Список volumes
docker volume ls

# Видалити volume (УВАГА: втрата даних!)
docker volume rm pets_search_rest_mongodb_data

# Інспекція volume
docker volume inspect pets_search_rest_mongodb_data

# Очистити всі невикористовувані volumes
docker volume prune
```

## 📊 Моніторинг сервісів

### Статус контейнерів
```bash
# Список запущених контейнерів
docker compose ps

# Статистика ресурсів
docker stats

# Логи всіх сервісів
docker compose logs -f

# Логи конкретного сервісу
docker compose logs -f api
```

### Health checks
```bash
# API
curl http://localhost:8080/healthz

# MongoDB
docker compose exec mongodb mongosh --eval "db.adminCommand('ping')"

# Redis
docker compose exec redis redis-cli ping

# MinIO
curl http://localhost:9000/minio/health/live
```

## 🔄 Lifecycle управління

### Запуск
```bash
# Запуск всіх сервісів
docker compose up -d

# Запуск конкретного сервісу
docker compose up -d mongodb

# Запуск з rebuild
docker compose up -d --build
```

### Зупинка
```bash
# Зупинка всіх сервісів
docker compose down

# Зупинка з видаленням volumes
docker compose down -v

# Зупинка конкретного сервісу
docker compose stop api
```

### Оновлення
```bash
# Перезапуск сервісу
docker compose restart api

# Rebuild та перезапуск
docker compose up -d --build api

# Оновлення образів
docker compose pull
docker compose up -d
```

## 🐛 Troubleshooting

### Проблеми з підключенням
```bash
# Перевірити мережу
docker network ls
docker network inspect pets_search_rest_pets-network

# Перевірити DNS
docker compose exec api nslookup mongodb
```

### Проблеми з volumes
```bash
# Права доступу
docker compose exec mongodb ls -la /data/db

# Місце на диску
docker system df
```

### Проблеми з портами
```bash
# Перевірити зайняті порти
sudo lsof -i :8080
sudo lsof -i :27017

# Знайти процес та завершити
sudo kill -9 <PID>
```
