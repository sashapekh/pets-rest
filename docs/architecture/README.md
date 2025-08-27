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
        │   MongoDB    │ │ Redis  │ │   MinIO     │
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
├── listings/              # Бізнес-логіка оголошень
├── users/                 # Управління користувачами
├── storage/               # Робота з файлами (S3/MinIO)
├── pdf/                   # Генерація PDF
├── qrcode/                # Генерація QR кодів
└── database/              # Робота з базою даних

pkg/                       # Публічні пакети
├── middleware/            # HTTP middleware
└── utils/                 # Загальні утиліти

web/                       # Веб ресурси
├── templates/             # HTML шаблони для SSR
└── static/                # CSS, JS, зображення
```

### Принципи архітектури

1. **Separation of Concerns** - кожен пакет має одну відповідальність
2. **Dependency Injection** - залежності передаються через конструктори
3. **Interface Segregation** - маленькі, специфічні інтерфейси
4. **Clean Code** - читабельний та підтримуваний код

## 🗄️ Схема бази даних (MongoDB)

### Колекція `users`
```javascript
{
  _id: ObjectId,
  email: String,           // required, unique
  phone: String,           // optional
  name: String,            // optional
  created_at: Date,        // required
  updated_at: Date         // optional
}
```

### Колекція `listings`
```javascript
{
  _id: ObjectId,
  user_id: ObjectId,       // required, ref: users
  type: String,            // required: "lost" | "found" | "adopt"
  title: String,           // required
  description: String,     // optional
  city: String,            // optional
  location: String,        // optional
  contact_phone: String,   // optional
  contact_tg: String,      // optional
  status: String,          // required: "draft" | "active" | "archived"
  slug: String,            // unique, for public URLs
  images: [String],        // array of image URLs
  created_at: Date,        // required
  updated_at: Date         // optional
}
```

### Колекція `events`
```javascript
{
  _id: ObjectId,
  user_id: ObjectId,       // optional, ref: users
  listing_id: ObjectId,    // required, ref: listings
  type: String,            // required: "view" | "qr_scan" | "contact_click"
  payload: Object,         // optional, additional data
  ip_address: String,      // optional
  user_agent: String,      // optional
  created_at: Date         // required
}
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
- **MongoDB** - NoSQL база даних
- **Redis** - кешування та сесії
- **MinIO** - S3-сумісне файлове сховище

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
POST   /auth/magic-link         # Відправка magic link
POST   /auth/magic-link/verify  # Верифікація

# Оголошення (потребують авторизації)
GET    /api/v1/listings         # Список оголошень
POST   /api/v1/listings         # Створення
GET    /api/v1/listings/:id     # Деталі
PUT    /api/v1/listings/:id     # Оновлення
DELETE /api/v1/listings/:id     # Видалення

# Файли
POST   /api/v1/listings/:id/images    # Завантаження фото
POST   /api/v1/listings/:id/pdf       # Генерація PDF

# Публічні сторінки
GET    /p/:slug                 # Публічне оголошення

# Система
GET    /healthz                 # Health check
```

## 🔐 Безпека

1. **Аутентифікація**: Magic links через email
2. **Авторизація**: JWT токени
3. **Валідація**: Всі вхідні дані
4. **Rate Limiting**: Обмеження запитів
5. **CORS**: Налаштований для frontend
6. **File Upload**: Обмеження типів та розмірів файлів
