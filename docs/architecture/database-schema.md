# Схема бази даних PostgreSQL

## 📊 Огляд таблиць

```
pets_search
├── users           # Користувачі системи
├── listings        # Оголошення про тварин
├── images          # Поліморфні зображення
└── events          # Аналітичні події
```

## 👤 Таблиця `users`

### Структура таблиці
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20),
    name VARCHAR(100),
    avatar_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE
);
```

### Поля
| Поле | Тип | Обов'язкове | Опис |
|------|-----|-------------|------|
| `id` | SERIAL | ✅ | Унікальний ідентифікатор (автоінкремент) |
| `email` | VARCHAR(255) | ✅ | Email адреса (унікальна) |
| `phone` | VARCHAR(20) | ❌ | Номер телефону |
| `name` | VARCHAR(100) | ❌ | Ім'я користувача |
| `avatar_url` | TEXT | ❌ | URL аватарки (синхронізується з images) |
| `created_at` | TIMESTAMPTZ | ✅ | Дата створення (автоматично) |
| `updated_at` | TIMESTAMPTZ | ❌ | Дата останнього оновлення |

### Індекси
```sql
-- Унікальний індекс на email (автоматично створюється)
CREATE INDEX idx_users_email ON users(email);

-- Індекс на телефон
CREATE INDEX idx_users_phone ON users(phone);
```

### Обмеження (Constraints)
```sql
-- Валідація email формату
ALTER TABLE users ADD CONSTRAINT email_format 
CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

-- Валідація телефону (міжнародний формат)
ALTER TABLE users ADD CONSTRAINT phone_format 
CHECK (phone IS NULL OR phone ~* '^\+[1-9]\d{1,14}$');
```

## 📋 Таблиця `listings`

### Структура таблиці
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
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE
);
```

### Поля
| Поле | Тип | Обов'язкове | Опис |
|------|-----|-------------|------|
| `id` | SERIAL | ✅ | Унікальний ідентифікатор (автоінкремент) |
| `user_id` | INTEGER | ✅ | Посилання на користувача (FK) |
| `type` | VARCHAR(10) | ✅ | Тип: `lost`, `found`, `adopt` |
| `title` | VARCHAR(255) | ✅ | Заголовок оголошення |
| `description` | TEXT | ❌ | Детальний опис |
| `city` | VARCHAR(100) | ❌ | Місто |
| `location` | VARCHAR(255) | ❌ | Конкретне місце |
| `contact_phone` | VARCHAR(20) | ❌ | Контактний телефон |
| `contact_tg` | VARCHAR(100) | ❌ | Telegram контакт |
| `status` | VARCHAR(10) | ✅ | Статус: `draft`, `active`, `archived` |
| `slug` | VARCHAR(255) | ❌ | URL slug (унікальний) |
| `created_at` | TIMESTAMPTZ | ✅ | Дата створення (автоматично) |
| `updated_at` | TIMESTAMPTZ | ❌ | Дата останнього оновлення |

### Індекси
```sql
-- Індекс на власника
CREATE INDEX idx_listings_user_id ON listings(user_id);

-- Індекс на тип оголошення
CREATE INDEX idx_listings_type ON listings(type);

-- Індекс на статус
CREATE INDEX idx_listings_status ON listings(status);

-- Унікальний індекс на slug
CREATE INDEX idx_listings_slug ON listings(slug);

-- Індекс на місто
CREATE INDEX idx_listings_city ON listings(city);

-- Індекс на дату створення (сортування)
CREATE INDEX idx_listings_created_at ON listings(created_at DESC);
```

## 🖼️ Таблиця `images` (Поліморфна)

### Структура таблиці
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
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE
);
```

### Поля
| Поле | Тип | Обов'язкове | Опис |
|------|-----|-------------|------|
| `id` | SERIAL | ✅ | Унікальний ідентифікатор (автоінкремент) |
| `imageable_type` | VARCHAR(50) | ✅ | Тип сутності: `user`, `listing` |
| `imageable_id` | INTEGER | ✅ | ID сутності |
| `url` | TEXT | ✅ | URL зображення |
| `filename` | VARCHAR(255) | ❌ | Оригінальне ім'я файлу |
| `size_bytes` | INTEGER | ❌ | Розмір файлу в байтах |
| `mime_type` | VARCHAR(100) | ❌ | MIME тип файлу |
| `alt_text` | TEXT | ❌ | Альтернативний текст |
| `sort_order` | INTEGER | ❌ | Порядок сортування (0, 1, 2...) |
| `is_primary` | BOOLEAN | ❌ | Основне зображення сутності |
| `created_at` | TIMESTAMPTZ | ✅ | Дата створення (автоматично) |
| `updated_at` | TIMESTAMPTZ | ❌ | Дата останнього оновлення |

### Типи сутностей
| Тип | Опис |
|-----|------|
| `user` | Аватарки користувачів |
| `listing` | Фото оголошень |

### Індекси
```sql
-- Поліморфний індекс
CREATE INDEX idx_images_polymorphic ON images(imageable_type, imageable_id);

-- Індекс для сортування
CREATE INDEX idx_images_sort ON images(imageable_type, imageable_id, sort_order);

-- Унікальний індекс для основного зображення
CREATE UNIQUE INDEX idx_images_unique_primary 
ON images(imageable_type, imageable_id) 
WHERE is_primary = TRUE;

-- Індекс на дату створення
CREATE INDEX idx_images_created_at ON images(created_at DESC);
```

### Обмеження (Constraints)
```sql
-- Валідація типу сутності
ALTER TABLE images ADD CONSTRAINT check_imageable_type 
CHECK (imageable_type IN ('user', 'listing'));

-- Тригер для updated_at
CREATE TRIGGER update_images_updated_at BEFORE UPDATE ON images
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

## 📈 Таблиця `events`

### Структура таблиці
```sql
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    listing_id INTEGER NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('view', 'qr_scan', 'contact_click', 'phone_click')),
    payload JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
```

### Поля
| Поле | Тип | Обов'язкове | Опис |
|------|-----|-------------|------|
| `id` | SERIAL | ✅ | Унікальний ідентифікатор (автоінкремент) |
| `user_id` | INTEGER | ❌ | Користувач (якщо авторизований) |
| `listing_id` | INTEGER | ✅ | Посилання на оголошення (FK) |
| `type` | VARCHAR(20) | ✅ | Тип події |
| `payload` | JSONB | ❌ | Додаткові дані |
| `ip_address` | INET | ❌ | IP адреса клієнта |
| `user_agent` | TEXT | ❌ | User Agent браузера |
| `created_at` | TIMESTAMPTZ | ✅ | Час події |

### Типи подій
| Тип | Опис |
|-----|------|
| `view` | Перегляд публічної сторінки |
| `qr_scan` | Скан QR коду |
| `contact_click` | Клік по контактним даним |
| `phone_click` | Клік по номеру телефону |

### Індекси
```sql
-- Індекс на оголошення
CREATE INDEX idx_events_listing_id ON events(listing_id);

-- Індекс на тип події
CREATE INDEX idx_events_type ON events(type);

-- Індекс на дату (для аналітики)
CREATE INDEX idx_events_created_at ON events(created_at DESC);

-- Індекс на користувача
CREATE INDEX idx_events_user_id ON events(user_id);

-- Складений індекс для аналітики
CREATE INDEX idx_events_analytics ON events(listing_id, type, created_at DESC);
```

## 🔍 Приклади запитів

### Пошук оголошень
```sql
-- Активні оголошення загублених тварин у Києві
SELECT * FROM listings 
WHERE type = 'lost' 
  AND status = 'active' 
  AND city ILIKE '%Київ%'
ORDER BY created_at DESC;

-- Пошук по тексту (використовуючи ILIKE)
SELECT * FROM listings 
WHERE status = 'active' 
  AND (title ILIKE '%кіт%' OR description ILIKE '%кіт%')
ORDER BY created_at DESC;
```

### Робота з зображеннями
```sql
-- Отримати всі зображення оголошення
SELECT * FROM images 
WHERE imageable_type = 'listing' 
  AND imageable_id = 123 
ORDER BY sort_order ASC;

-- Отримати основне зображення користувача
SELECT * FROM images 
WHERE imageable_type = 'user' 
  AND imageable_id = 456 
  AND is_primary = TRUE;

-- Оголошення з їх основними зображеннями
SELECT l.*, i.url as primary_image_url
FROM listings l
LEFT JOIN images i ON (
    i.imageable_type = 'listing' 
    AND i.imageable_id = l.id 
    AND i.is_primary = TRUE
)
WHERE l.status = 'active'
ORDER BY l.created_at DESC;
```

### Аналітика
```sql
-- Кількість переглядів оголошення
SELECT COUNT(*) as views
FROM events 
WHERE listing_id = 123 
  AND type = 'view';

-- Популярні оголошення за останній тиждень
SELECT 
    l.id,
    l.title,
    COUNT(e.id) as views
FROM listings l
LEFT JOIN events e ON (e.listing_id = l.id AND e.type = 'view')
WHERE l.status = 'active' 
  AND e.created_at >= NOW() - INTERVAL '7 days'
GROUP BY l.id, l.title
ORDER BY views DESC
LIMIT 10;
```

### Статистика користувача
```sql
-- Всі оголошення користувача з кількістю переглядів
SELECT 
    l.*,
    COUNT(e.id) as total_views,
    COUNT(CASE WHEN e.type = 'contact_click' THEN 1 END) as contact_clicks
FROM listings l
LEFT JOIN events e ON e.listing_id = l.id
WHERE l.user_id = 789
GROUP BY l.id
ORDER BY l.created_at DESC;

-- Активність по оголошенням користувача за останній місяць
SELECT 
    DATE_TRUNC('day', e.created_at) as date,
    e.type,
    COUNT(*) as count
FROM events e
JOIN listings l ON l.id = e.listing_id
WHERE l.user_id = 789 
  AND e.created_at >= NOW() - INTERVAL '30 days'
GROUP BY DATE_TRUNC('day', e.created_at), e.type
ORDER BY date DESC, e.type;
```

## 🔄 Зв'язки між таблицями

```
users (1) ──→ (N) listings
  │
  └── (1) ──→ (N) images [imageable_type='user']

listings (1) ──→ (N) events
   │
   └── (1) ──→ (N) images [imageable_type='listing']

events (N) ──→ (1) listings
events (N) ──→ (1) users [optional]
```

## 📝 Примітки

1. **Поліморфна таблиця images** дозволяє зберігати зображення для різних типів сутностей (користувачі, оголошення) в одній таблиці
2. **avatar_url в users** дублює URL з images таблиці для швидкого доступу
3. **Тригери updated_at** автоматично оновлюють час модифікації записів
4. **Каскадне видалення** забезпечує цілісність даних при видаленні користувачів або оголошень
5. **JSONB payload** в events дозволяє зберігати додаткові дані про події в гнучкому форматі