# Поліморфна таблиця для зображень

## Огляд

Ми реалізували поліморфну таблицю `images` для зберігання зображень, які можуть належати різним сутностям (користувачі, оголошення, тощо).

## Структура таблиці

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

## Типи сутностей

- `user` - зображення користувача (аватарки)
- `listing` - зображення оголошення

## Go моделі

### Image
```go
type Image struct {
    ID            int           `json:"id" db:"id"`
    ImageableType ImageableType `json:"imageable_type" db:"imageable_type"`
    ImageableID   int           `json:"imageable_id" db:"imageable_id"`
    URL           string        `json:"url" db:"url"`
    Filename      *string       `json:"filename,omitempty" db:"filename"`
    SizeBytes     *int          `json:"size_bytes,omitempty" db:"size_bytes"`
    MimeType      *string       `json:"mime_type,omitempty" db:"mime_type"`
    AltText       *string       `json:"alt_text,omitempty" db:"alt_text"`
    SortOrder     int           `json:"sort_order" db:"sort_order"`
    IsPrimary     bool          `json:"is_primary" db:"is_primary"`
    CreatedAt     time.Time     `json:"created_at" db:"created_at"`
    UpdatedAt     *time.Time    `json:"updated_at,omitempty" db:"updated_at"`
}
```

### ImageableType
```go
type ImageableType string

const (
    ImageableTypeUser    ImageableType = "user"
    ImageableTypeListing ImageableType = "listing"
)
```

## ImageService

### Основні методи

#### Створення зображення
```go
// Створити зображення з URL
image, err := imageService.CreateImageFromURL(
    database.ImageableTypeUser, 
    userID, 
    "https://example.com/avatar.jpg", 
    true, // isPrimary
)

// Створити зображення вручну
req := CreateImageRequest{
    ImageableType: database.ImageableTypeListing,
    ImageableID:   listingID,
    URL:           "https://example.com/pet.jpg",
    SortOrder:     0,
    IsPrimary:     true,
}
image, err := imageService.CreateImage(req)
```

#### Отримання зображень
```go
// Всі зображення для сутності
images, err := imageService.GetImagesByEntity(database.ImageableTypeUser, userID)

// Основне зображення
primaryImage, err := imageService.GetPrimaryImage(database.ImageableTypeUser, userID)
```

#### Управління зображеннями
```go
// Встановити як основне
err := imageService.SetPrimaryImage(imageID)

// Видалити зображення
err := imageService.DeleteImage(imageID)

// Видалити всі зображення сутності
err := imageService.DeleteImagesByEntity(database.ImageableTypeUser, userID)

// Оновити порядок зображень
orders := []struct {
    ID    int `json:"id"`
    Order int `json:"order"`
}{
    {ID: 1, Order: 0},
    {ID: 2, Order: 1},
}
err := imageService.UpdateImageOrder(database.ImageableTypeListing, listingID, orders)
```

## UserService інтеграція

UserService автоматично зберігає аватарки користувачів при OAuth авторизації:

```go
// При вході через Google OAuth
func (s *UserService) FirstOrNewUserForRegister(user *oauth.User) (*database.User, error) {
    // ... логіка створення/пошуку користувача
    
    // Автоматично зберігає аватарку в images таблицю
    if user.AvatarURL != "" {
        s.saveUserAvatar(newUser.ID, user.AvatarURL)
    }
    
    return newUser, nil
}
```

### Додаткові методи UserService

```go
// Отримати аватарку користувача
avatar, err := userService.GetUserAvatar(userID)

// Отримати всі зображення користувача
images, err := userService.GetUserImages(userID)
```

## Переваги поліморфної структури

1. **Гнучкість** - легко додавати нові типи сутностей
2. **Уніфікованість** - одна таблиця для всіх зображень
3. **Метадані** - зберігання додаткової інформації про зображення
4. **Сортування** - контроль порядку зображень
5. **Основне зображення** - позначення головного зображення
6. **Масштабованість** - легко розширювати функціональність

## Міграція даних

Існуючі зображення з `listings.images` були автоматично перенесені до нової таблиці за допомогою міграції `006_migrate_listing_images_to_images_table.up.sql`.

## Приклади використання в API

### Отримання зображень оголошення
```go
func (h *ListingHandler) GetListingImages(c fiber.Ctx) error {
    listingID := c.Params("id")
    id, _ := strconv.Atoi(listingID)
    
    images, err := h.imageService.GetImagesByEntity(database.ImageableTypeListing, id)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(images)
}
```

### Додавання зображення до оголошення
```go
func (h *ListingHandler) AddListingImage(c fiber.Ctx) error {
    listingID := c.Params("id")
    id, _ := strconv.Atoi(listingID)
    
    var req services.CreateImageRequest
    if err := c.Bind().Body(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": err.Error()})
    }
    
    req.ImageableType = database.ImageableTypeListing
    req.ImageableID = id
    
    image, err := h.imageService.CreateImage(req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.Status(201).JSON(image)
}
```
