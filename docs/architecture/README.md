# Архітектура проекту

## 🏗️ Загальна архітектура

Pets Search REST API побудований на основі чистої архітектури (Clean Architecture) з розділенням на шари.

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Mobile App    │    │   Public Pages  │
│   (React/Vue)   │    │   (React Native)│    │   (SSR)         │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   REST API      │
                    │   (Fiber v3)    │
                    └─────────┬───────┘
                              │
                ┌─────────────┼─────────────┐
                │             │             │
        ┌───────▼──────┐ ┌───▼────┐ ┌──────▼──────┐
        │  PostgreSQL  │ │ Redis  │ │   MinIO     │
        │   (Database) │ │(Cache) │ │ (Files/S3)  │
        └──────────────┘ └────────┘ └─────────────┘
```

## 📂 Структура проекту

### Основні папки

```
cmd/api/                    # Вхідна точка програми
├── main.go                # HTTP сервер та маршрутизація

internal/                   # Приватний код програми
├── config/                # Конфігурація
├── handlers/              # HTTP обробники
├── services/              # Бізнес-логіка сервісів
├── database/              # Моделі та репозиторії
├── auth/                  # Аутентифікація та JWT
├── oauth/                 # OAuth провайдери (Google)
├── fx/                    # Dependency Injection
└── storage/               # Робота з файлами (S3/MinIO)

pkg/                       # Публічні пакети
├── helper/                # Загальні утиліти
└── middleware/            # HTTP middleware

migrations/                # Міграції бази даних
├── *.up.sql              # Міграції вгору
└── *.down.sql            # Міграції назад

web/                       # Веб ресурси (якщо потрібно)
├── templates/             # HTML шаблони для SSR
└── static/                # CSS, JS, зображення
```

### Принципи архітектури

1. **Separation of Concerns** - кожен пакет має одну відповідальність
2. **Dependency Injection** - залежності передаються через конструктори
3. **Interface Segregation** - маленькі, специфічні інтерфейси
4. **Clean Code** - читабельний та підтримуваний код

## 🗄️ Схема бази даних (PostgreSQL)

### Таблиця `users`
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20),
    name VARCHAR(100),
    avatar_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);
```

### Таблиця `listings`
```sql
CREATE TABLE listings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(10) NOT NULL CHECK (type IN ('lost', 'found', 'adopt')),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    city VARCHAR(100),
    location VARCHAR(255),
    contact_phone VARCHAR(20),
    contact_tg VARCHAR(100),
    status VARCHAR(10) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'archived')),
    slug VARCHAR(255) UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);
```

### Таблиця `events`
```sql
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    listing_id INTEGER NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('view', 'qr_scan', 'contact_click', 'phone_click')),
    payload JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### Таблиця `images` (поліморфна)
```sql
CREATE TABLE images (
    id SERIAL PRIMARY KEY,
    imageable_type VARCHAR(50) NOT NULL,
    imageable_id INTEGER NOT NULL,
    url TEXT NOT NULL,
    filename VARCHAR(255),
    size_bytes INTEGER,
    mime_type VARCHAR(100),
    alt_text TEXT,
    sort_order INTEGER DEFAULT 0,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);
```

## 🔄 Потік даних

### 1. Створення оголошення
```
User → API → Validation → Database → File Upload → QR/PDF Generation
```

### 2. Перегляд оголошення
```
Request → API → Database → Cache Check → Response + Event Logging
```

### 3. Пошук оголошень
```
Search Query → API → Database Query → Cache → Filtered Results
```

## 🔧 Технологічний стек

### Backend
- **Go 1.22+** - мова програмування
- **Fiber v3** - HTTP framework
- **PostgreSQL** - реляційна база даних
- **Redis** - кешування та сесії
- **MinIO** - S3-сумісне файлове сховище
- **golang-migrate** - міграції бази даних

### Інфраструктура
- **Docker** - контейнеризація
- **Docker Compose** - оркестрація сервісів
- **Make** - автоматизація збірки

### Зовнішні сервіси
- **SMTP** - відправка email (magic links)
- **S3/MinIO** - зберігання фото та PDF
- **QR код генератор** - для друкованих матеріалів

## 🚦 API Routes

```
# Аутентифікація
GET    /auth/google             # Google OAuth вхід
GET    /auth/google/callback    # Google OAuth callback

# Оголошення (потребують авторизації)
GET    /api/v1/listings         # Список оголошень
POST   /api/v1/listings         # Створення
GET    /api/v1/listings/:id     # Деталі
PUT    /api/v1/listings/:id     # Оновлення
DELETE /api/v1/listings/:id     # Видалення

# Зображення
GET    /api/v1/listings/:id/images    # Отримання зображень оголошення
POST   /api/v1/listings/:id/images    # Додавання зображення
DELETE /api/v1/images/:id            # Видалення зображення
PUT    /api/v1/images/:id/primary     # Встановити як основне

# Профіль
GET    /api/v1/profile              # Профіль користувача
GET    /api/v1/profile/images       # Зображення користувача

# Файли
POST   /api/v1/listings/:id/pdf       # Генерація PDF

# Публічні сторінки
GET    /p/:slug                 # Публічне оголошення

# Система
GET    /healthz                 # Health check
```

## 🔐 Безпека

1. **Аутентифікація**: Google OAuth 2.0
2. **Авторизація**: JWT токени
3. **Валідація**: Всі вхідні дані
4. **Rate Limiting**: Обмеження запитів
5. **CORS**: Налаштований для frontend
6. **File Upload**: Обмеження типів та розмірів файлів
