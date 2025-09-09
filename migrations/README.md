# Database Migrations

Цей проект використовує [golang-migrate](https://github.com/golang-migrate/migrate) для управління схемою бази даних PostgreSQL.

## 📁 Структура міграцій

```
migrations/
├── 001_create_users_table.up.sql       # Створення таблиці користувачів
├── 001_create_users_table.down.sql     # Видалення таблиці користувачів
├── 002_create_listings_table.up.sql    # Створення таблиці оголошень
├── 002_create_listings_table.down.sql  # Видалення таблиці оголошень
├── 003_create_events_table.up.sql      # Створення таблиці подій
├── 003_create_events_table.down.sql    # Видалення таблиці подій
└── README.md                           # Цей файл
```

## 🚀 Команди

### Використання через Makefile
```bash
# Запустити всі міграції (up)
make migrate-up

# Відкотити всі міграції (down)
make migrate-down

# Перевірити поточну версію
make migrate-version

# Скинути базу даних (небезпечно!)
make db-reset
```

### Прямі команди
```bash
# Запустити міграції
go run ./cmd/migrate -up

# Відкотити міграції
go run ./cmd/migrate -down

# Перевірити версію
go run ./cmd/migrate -version
```

## 🔄 Автоматичне виконання

Міграції **автоматично виконуються** при запуску API сервера. Це означає:

1. При старті контейнера `api` міграції запускаються автоматично
2. Якщо міграції вже виконані, вони пропускаються
3. При помилці міграції сервер не запуститься

## 📝 Створення нових міграцій

### 1. Найменування файлів
Формат: `{version}_{description}.{up|down}.sql`

Приклад:
```
004_add_user_avatar.up.sql
004_add_user_avatar.down.sql
```

### 2. Приклад UP міграції
```sql
-- 004_add_user_avatar.up.sql
ALTER TABLE users ADD COLUMN avatar_url VARCHAR(255);
CREATE INDEX idx_users_avatar ON users(avatar_url);
```

### 3. Приклад DOWN міграції
```sql
-- 004_add_user_avatar.down.sql
DROP INDEX IF EXISTS idx_users_avatar;
ALTER TABLE users DROP COLUMN IF EXISTS avatar_url;
```

## ⚠️ Важливі правила

### DO ✅
- Завжди створюйте як UP, так і DOWN міграції
- Використовуйте `IF EXISTS` / `IF NOT EXISTS` де можливо
- Тестуйте міграції на копії production даних
- Робіть backup перед важливими змінами
- Дотримуйтесь послідовної нумерації

### DON'T ❌
- Не редагуйте існуючі міграції після їх застосування
- Не видаляйте міграції, які вже були виконані
- Не використовуйте DDL команди, які не можна відкотити
- Не додавайте NOT NULL колонки без DEFAULT значень

## 🛠️ Troubleshooting

### Помилка "dirty database"
```bash
# Перевірити стан
make migrate-version

# Якщо база "dirty", можна спробувати:
# 1. Виправити проблему вручну в базі
# 2. Або скинути базу (втрата даних!)
make db-reset
```

### Міграція застрягла
```bash
# Подивитися логи бази даних
make db-logs

# Підключитися до бази та перевірити стан
docker compose exec postgres psql -U pets_user -d pets_search
```

### Відкат конкретної міграції
Для відкату до конкретної версії потрібно використовувати інструмент migrate напряму:

```bash
# Встановити migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Відкотити до версії N
migrate -database "postgres://pets_user:pets_password@localhost:5432/pets_search?sslmode=disable" -path migrations goto N
```

## 📊 Поточна схема

### Версія 1: Користувачі
- Таблиця `users` з базовими полями
- Валідація email та телефону
- Тригери для `updated_at`

### Версія 2: Оголошення
- Таблиця `listings` з зовнішнім ключем до `users`
- Підтримка масивів зображень (PostgreSQL arrays)
- Check constraints для enum значень

### Версія 3: Аналітика
- Таблиця `events` для збору метрик
- JSONB поле для гнучких даних
- Індекси для швидких запитів
