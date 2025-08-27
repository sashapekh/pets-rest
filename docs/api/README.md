# API Документація

## 🔗 Base URL
```
http://localhost:8080
```

## 🚦 Статус відповідей

| Код | Статус | Опис |
|-----|--------|------|
| 200 | OK | Успішна операція |
| 201 | Created | Ресурс створений |
| 400 | Bad Request | Помилка валідації |
| 401 | Unauthorized | Не авторизований |
| 403 | Forbidden | Немає прав доступу |
| 404 | Not Found | Ресурс не знайдений |
| 500 | Internal Server Error | Серверна помилка |

## 🔐 Аутентифікація

### Magic Link Flow

#### 1. Запит magic link
```http
POST /auth/magic-link
Content-Type: application/json

{
  "email": "user@example.com"
}
```

**Відповідь:**
```json
{
  "message": "Magic link sent to your email",
  "success": true
}
```

#### 2. Верифікація magic link
```http
POST /auth/magic-link/verify
Content-Type: application/json

{
  "token": "magic_link_token_from_email"
}
```

**Відповідь:**
```json
{
  "access_token": "jwt_token_here",
  "user": {
    "id": "user_id",
    "email": "user@example.com",
    "name": "User Name"
  }
}
```

### Використання JWT токену
```http
Authorization: Bearer <jwt_token>
```

## 📋 Оголошення (Listings)

### Отримання списку оголошень
```http
GET /api/v1/listings?type=lost&city=Kyiv&page=1&limit=10
```

**Query параметри:**
- `type` - тип оголошення: `lost`, `found`, `adopt`
- `city` - місто
- `status` - статус: `draft`, `active`, `archived`
- `page` - номер сторінки (default: 1)
- `limit` - кількість на сторінці (default: 10)

**Відповідь:**
```json
{
  "listings": [
    {
      "id": "listing_id",
      "title": "Загубився кіт Мурзик",
      "type": "lost",
      "description": "Сірий кіт з білими лапками...",
      "city": "Київ",
      "location": "Район Печерський",
      "contact_phone": "+380501234567",
      "contact_tg": "@username",
      "status": "active",
      "images": [
        "https://minio.local/pets-photos/img1.jpg"
      ],
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 25,
    "pages": 3
  }
}
```

### Створення оголошення
```http
POST /api/v1/listings
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Загубився кіт Мурзик",
  "type": "lost",
  "description": "Сірий кіт з білими лапками, дуже ласкавий",
  "city": "Київ",
  "location": "Район Печерський, поблизу метро Арсенальна",
  "contact_phone": "+380501234567",
  "contact_tg": "@username",
  "status": "draft"
}
```

**Відповідь:**
```json
{
  "id": "new_listing_id",
  "slug": "lost-cat-murzik-kyiv-abc123",
  "message": "Listing created successfully"
}
```

### Отримання оголошення
```http
GET /api/v1/listings/{id}
```

**Відповідь:** Повна інформація про оголошення

### Оновлення оголошення
```http
PUT /api/v1/listings/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Оновлена назва",
  "status": "active"
}
```

### Видалення оголошення
```http
DELETE /api/v1/listings/{id}
Authorization: Bearer <token>
```

## 📸 Завантаження зображень

### Завантаження фото до оголошення
```http
POST /api/v1/listings/{id}/images
Authorization: Bearer <token>
Content-Type: multipart/form-data

file=@/path/to/image.jpg
```

**Відповідь:**
```json
{
  "image_url": "https://minio.local/pets-photos/listing_id/image_123.jpg",
  "message": "Image uploaded successfully"
}
```

**Обмеження:**
- Максимальний розмір: 10MB
- Дозволені формати: JPG, JPEG, PNG, GIF
- Максимум 5 фото на оголошення

## 📄 Генерація PDF

### Створення PDF для друку
```http
POST /api/v1/listings/{id}/generate-pdf
Authorization: Bearer <token>
Content-Type: application/json

{
  "template": "a4",  // "a4", "a5", "business_card"
  "include_qr": true
}
```

**Відповідь:**
```json
{
  "pdf_url": "https://minio.local/pets-photos/listing_id/poster_a4.pdf",
  "qr_code_url": "https://minio.local/pets-photos/listing_id/qr_code.png",
  "public_url": "https://yourdomain.com/p/lost-cat-murzik-kyiv-abc123"
}
```

## 🌐 Публічні сторінки

### Перегляд публічного оголошення
```http
GET /p/{slug}
```

**Відповідь:** JSON або HTML сторінка залежно від заголовка `Accept`

**JSON відповідь:**
```json
{
  "listing": {
    "id": "listing_id",
    "title": "Загубився кіт Мурзик",
    "type": "lost",
    "description": "...",
    "images": ["..."],
    "contact_phone": "+380501234567",
    "contact_tg": "@username",
    "created_at": "2024-01-01T12:00:00Z"
  },
  "qr_code": "https://minio.local/pets-photos/listing_id/qr_code.png"
}
```

## 📊 Аналітика

### Події відстеження
API автоматично відстежує наступні події:

- `view` - перегляд публічної сторінки
- `qr_scan` - скан QR коду
- `contact_click` - клік по контактам
- `phone_click` - клік по телефону

**Приклад події:**
```json
{
  "listing_id": "listing_id",
  "type": "view",
  "ip_address": "192.168.1.1",
  "user_agent": "Mozilla/5.0...",
  "payload": {
    "referrer": "https://google.com",
    "device": "mobile"
  },
  "created_at": "2024-01-01T12:00:00Z"
}
```

## 🔧 Система

### Health Check
```http
GET /healthz
```

**Відповідь:**
```json
{
  "status": "ok",
  "message": "Pets Search API is running",
  "services": {
    "database": "connected",
    "redis": "connected",
    "storage": "connected"
  },
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## ❌ Обробка помилок

### Стандартний формат помилки
```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": {
    "field": "validation error description"
  }
}
```

### Приклади помилок

**400 Bad Request:**
```json
{
  "error": "Validation failed",
  "code": "VALIDATION_ERROR",
  "details": {
    "title": "Title is required",
    "type": "Type must be one of: lost, found, adopt"
  }
}
```

**401 Unauthorized:**
```json
{
  "error": "Authentication required",
  "code": "UNAUTHORIZED"
}
```

**404 Not Found:**
```json
{
  "error": "Listing not found",
  "code": "NOT_FOUND"
}
```

## 🧪 Приклади використання

### cURL приклади

**Створення оголошення:**
```bash
curl -X POST http://localhost:8080/api/v1/listings \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Загубився кіт",
    "type": "lost",
    "description": "Сірий кіт з білими лапками",
    "city": "Київ",
    "contact_phone": "+380501234567"
  }'
```

**Завантаження фото:**
```bash
curl -X POST http://localhost:8080/api/v1/listings/LISTING_ID/images \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@cat_photo.jpg"
```

### JavaScript приклади

**Fetch API:**
```javascript
// Створення оголошення
const response = await fetch('/api/v1/listings', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    title: 'Загубився кіт',
    type: 'lost',
    description: 'Сірий кіт з білими лапками',
    city: 'Київ'
  })
});

const result = await response.json();
```

**Завантаження файлу:**
```javascript
const formData = new FormData();
formData.append('file', fileInput.files[0]);

const response = await fetch(`/api/v1/listings/${listingId}/images`, {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`
  },
  body: formData
});
```
