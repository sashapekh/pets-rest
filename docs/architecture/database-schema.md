# Схема бази даних MongoDB

## 📊 Огляд колекцій

```
pets_search
├── users           # Користувачі системи
├── listings        # Оголошення про тварин
└── events          # Аналітичні події
```

## 👤 Колекція `users`

### Структура документу
```javascript
{
  _id: ObjectId("507f1f77bcf86cd799439011"),
  email: "user@example.com",
  phone: "+380501234567",
  name: "Іван Петренко",
  created_at: ISODate("2024-01-01T12:00:00Z"),
  updated_at: ISODate("2024-01-01T12:00:00Z")
}
```

### Поля
| Поле | Тип | Обов'язкове | Опис |
|------|-----|-------------|------|
| `_id` | ObjectId | ✅ | Унікальний ідентифікатор |
| `email` | String | ✅ | Email адреса (унікальна) |
| `phone` | String | ❌ | Номер телефону |
| `name` | String | ❌ | Ім'я користувача |
| `created_at` | Date | ✅ | Дата створення |
| `updated_at` | Date | ❌ | Дата останнього оновлення |

### Індекси
```javascript
// Унікальний індекс на email
db.users.createIndex({ email: 1 }, { unique: true })

// Індекс на телефон
db.users.createIndex({ phone: 1 })
```

### Валідація
```javascript
{
  $jsonSchema: {
    bsonType: "object",
    required: ["email", "created_at"],
    properties: {
      email: {
        bsonType: "string",
        pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
      },
      phone: {
        bsonType: "string",
        pattern: "^\\+[1-9]\\d{1,14}$"
      },
      name: {
        bsonType: "string",
        minLength: 1,
        maxLength: 100
      }
    }
  }
}
```

## 📋 Колекція `listings`

### Структура документу
```javascript
{
  _id: ObjectId("507f1f77bcf86cd799439012"),
  user_id: ObjectId("507f1f77bcf86cd799439011"),
  type: "lost",
  title: "Загубився кіт Мурзик",
  description: "Сірий кіт з білими лапками, дуже ласкавий...",
  city: "Київ",
  location: "Район Печерський, поблизу метро Арсенальна",
  contact_phone: "+380501234567",
  contact_tg: "@username",
  status: "active",
  slug: "lost-cat-murzik-kyiv-abc123",
  images: [
    "https://minio.local/pets-photos/listing_id/image1.jpg",
    "https://minio.local/pets-photos/listing_id/image2.jpg"
  ],
  created_at: ISODate("2024-01-01T12:00:00Z"),
  updated_at: ISODate("2024-01-01T13:30:00Z")
}
```

### Поля
| Поле | Тип | Обов'язкове | Опис |
|------|-----|-------------|------|
| `_id` | ObjectId | ✅ | Унікальний ідентифікатор |
| `user_id` | ObjectId | ✅ | Посилання на користувача |
| `type` | String | ✅ | Тип: `lost`, `found`, `adopt` |
| `title` | String | ✅ | Заголовок оголошення |
| `description` | String | ❌ | Детальний опис |
| `city` | String | ❌ | Місто |
| `location` | String | ❌ | Конкретне місце |
| `contact_phone` | String | ❌ | Контактний телефон |
| `contact_tg` | String | ❌ | Telegram контакт |
| `status` | String | ✅ | Статус: `draft`, `active`, `archived` |
| `slug` | String | ❌ | URL slug (унікальний) |
| `images` | Array | ❌ | Масив URL зображень |
| `created_at` | Date | ✅ | Дата створення |
| `updated_at` | Date | ❌ | Дата останнього оновлення |

### Індекси
```javascript
// Індекс на власника
db.listings.createIndex({ user_id: 1 })

// Індекс на тип оголошення
db.listings.createIndex({ type: 1 })

// Індекс на статус
db.listings.createIndex({ status: 1 })

// Унікальний індекс на slug
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
