# Pets Search API (Golang + Fiber)

## 🎯 Призначення
REST API для сервісу пошуку та оголошень про тварин. 
API відповідає за:
- створення та управління оголошеннями (знайдені, загублені, віддам у добрі руки);
- зберігання фото (через S3 або інше сумісне сховище);
- генерацію QR-кодів і PDF для друку;
- публічні сторінки для оголошень;
- аналітику (перегляди, скани QR).

---

## ⚙️ Стек
- Go (v1.25+)
- Fiber v2
- MongoDB
- Redis (сесії, кеш)
- S3 (MinIO / AWS S3) для зберігання фото та PDF
- Docker Compose для локального запуску

---

## 📂 Архітектура
- `cmd/api` — вхідна точка сервера
- `internal/listings` — бізнес-логіка оголошень
- `internal/users` — авторизація/ідентифікація (magic link, email/phone)
- `internal/storage` — інтеграції з S3
- `internal/pdf` — генерація PDF
- `internal/qrcode` — генерація QR

---

## 🔑 Функціонал

### 1. Аутентифікація
- Magic link через email (на майбутнє — через SMS/Telegram).
- Мінімальна таблиця користувачів.

### 2. Оголошення (Listings)
- CRUD для оголошень.
- Поля:  
  - `id`, `user_id`  
  - `type` (lost | found | adopt)  
  - `title`, `description`, `city`, `location`  
  - `contact_phone`, `contact_tg`  
  - `status` (draft | active | archived)  
  - `created_at`, `updated_at`
- Завантаження фото (pre-signed URLs → S3).
- Генерація QR → `https://domain/p/{slug}`
- Генерація PDF для друку (шаблони: A4, A5, візитка).

### 3. Публічні сторінки
- `GET /p/{slug}` → віддає JSON для фронта або SSR-сторінку.
- Відкриті без авторизації.
- Мають OG-теги для гарного відображення в Telegram/FB.

### 4. Аналітика
- Події: перегляд сторінки, скан QR.
- Зберігаються в таблиці `events` (user_id?, listing_id, type, payload, created_at).
- Використовуються для статистики (наприклад, скільки людей відсканували QR).

---

## 📌 Основні endpoints (чернетка)
- `POST /auth/magic-link`
- `POST /auth/magic-link/verify`
- `GET /api/v1/listings`
- `POST /api/v1/listings`
- `GET /api/v1/listings/{id}`
- `PUT /api/v1/listings/{id}`
- `DELETE /api/v1/listings/{id}`
- `POST /api/v1/listings/{id}/images`
- `POST /api/v1/listings/{id}/generate-pdf`
- `GET /healthz`

---
