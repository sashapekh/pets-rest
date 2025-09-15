# Схема бази даних PostgreSQL

## 📊 Огляд таблиць

```
pets_search
├── users           # Користувачі системи
├── listings        # Оголошення про тварин
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
    images TEXT[], -- Array of image URLs
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
| `images` | TEXT[] | ❌ | Масив URL зображень |
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
db.listings.createIndex({ slug: 1 }, { unique: true, sparse: true })

// Індекс на місто
db.listings.createIndex({ city: 1 })

// Індекс на дату створення (сортування)
db.listings.createIndex({ created_at: -1 })

// Складений індекс для пошуку
db.listings.createIndex({ 
  type: 1, 
  status: 1, 
  city: 1, 
  created_at: -1 
})
```

### Валідація
```javascript
{
  $jsonSchema: {
    bsonType: "object",
    required: ["user_id", "type", "title", "status", "created_at"],
    properties: {
      user_id: {
        bsonType: "objectId"
      },
      type: {
        bsonType: "string",
        enum: ["lost", "found", "adopt"]
      },
      title: {
        bsonType: "string",
        minLength: 3,
        maxLength: 200
      },
      description: {
        bsonType: "string",
        maxLength: 2000
      },
      status: {
        bsonType: "string",
        enum: ["draft", "active", "archived"]
      },
      slug: {
        bsonType: "string",
        pattern: "^[a-z0-9-]+$"
      },
      images: {
        bsonType: "array",
        maxItems: 5,
        items: {
          bsonType: "string"
        }
      }
    }
  }
}
```

## 📈 Колекція `events`

### Структура документу
```javascript
{
  _id: ObjectId("507f1f77bcf86cd799439013"),
  user_id: ObjectId("507f1f77bcf86cd799439011"), // опціонально
  listing_id: ObjectId("507f1f77bcf86cd799439012"),
  type: "view",
  payload: {
    referrer: "https://google.com",
    device: "mobile",
    browser: "Chrome",
    coordinates: {
      lat: 50.4501,
      lng: 30.5234
    }
  },
  ip_address: "192.168.1.1",
  user_agent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)...",
  created_at: ISODate("2024-01-01T12:00:00Z")
}
```

### Поля
| Поле | Тип | Обов'язкове | Опис |
|------|-----|-------------|------|
| `_id` | ObjectId | ✅ | Унікальний ідентифікатор |
| `user_id` | ObjectId | ❌ | Користувач (якщо авторизований) |
| `listing_id` | ObjectId | ✅ | Посилання на оголошення |
| `type` | String | ✅ | Тип події |
| `payload` | Object | ❌ | Додаткові дані |
| `ip_address` | String | ❌ | IP адреса клієнта |
| `user_agent` | String | ❌ | User Agent браузера |
| `created_at` | Date | ✅ | Час події |

### Типи подій
| Тип | Опис |
|-----|------|
| `view` | Перегляд публічної сторінки |
| `qr_scan` | Скан QR коду |
| `contact_click` | Клік по контактним даним |
| `phone_click` | Клік по номеру телефону |
| `telegram_click` | Клік по Telegram |
| `image_view` | Перегляд фото |

### Індекси
```javascript
// Індекс на оголошення
db.events.createIndex({ listing_id: 1 })

// Індекс на тип події
db.events.createIndex({ type: 1 })

// Індекс на дату (для аналітики)
db.events.createIndex({ created_at: -1 })

// Індекс на користувача
db.events.createIndex({ user_id: 1 })

// Складений індекс для аналітики
db.events.createIndex({ 
  listing_id: 1, 
  type: 1, 
  created_at: -1 
})
```

## 🔍 Приклади запитів

### Пошук оголошень
```javascript
// Активні оголошення загублених тварин у Києві
db.listings.find({
  type: "lost",
  status: "active",
  city: "Київ"
}).sort({ created_at: -1 })

// Пошук по тексту (потребує text index)
db.listings.find({
  $text: { $search: "кіт сірий" },
  status: "active"
})
```

### Аналітика
```javascript
// Кількість переглядів оголошення
db.events.countDocuments({
  listing_id: ObjectId("507f1f77bcf86cd799439012"),
  type: "view"
})

// Популярні оголошення за останній тиждень
db.events.aggregate([
  {
    $match: {
      type: "view",
      created_at: { 
        $gte: ISODate("2024-01-01T00:00:00Z") 
      }
    }
  },
  {
    $group: {
      _id: "$listing_id",
      views: { $sum: 1 }
    }
  },
  { $sort: { views: -1 } },
  { $limit: 10 }
])
```

### Статистика користувача
```javascript
// Всі оголошення користувача
db.listings.find({
  user_id: ObjectId("507f1f77bcf86cd799439011")
}).sort({ created_at: -1 })

// Активність по оголошенням користувача
db.events.aggregate([
  {
    $lookup: {
      from: "listings",
      localField: "listing_id",
      foreignField: "_id",
      as: "listing"
    }
  },
  {
    $match: {
      "listing.user_id": ObjectId("507f1f77bcf86cd799439011")
    }
  },
  {
    $group: {
      _id: "$listing_id",
      total_events: { $sum: 1 },
      views: {
        $sum: { $cond: [{ $eq: ["$type", "view"] }, 1, 0] }
      },
      contacts: {
        $sum: { $cond: [{ $eq: ["$type", "contact_click"] }, 1, 0] }
      }
    }
  }
])
```
