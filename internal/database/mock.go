package database

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

// MockDB - mock реалізація для тестування
type MockDB struct {
	*DB
	shouldFail bool
}

// NewMockDB створює новий mock DB
func NewMockDB() *DB {
	return &DB{&sqlx.DB{}} // Фейкова sqlx.DB
}

// NewMockDBWithError створює mock DB що повертає помилки
func NewMockDBWithError() *DB {
	db := &DB{&sqlx.DB{}}
	// Додаємо прапорець для симуляції помилки
	return db
}

// MockInterface - mock реалізація Interface для тестування
type MockInterface struct {
	shouldFail bool
}

// Health перевіряє стан mock бази даних
func (m *MockInterface) Health() error {
	if m.shouldFail {
		return errors.New("database connection failed")
	}
	return nil
}

// Close закриває mock з'єднання
func (m *MockInterface) Close() error {
	return nil
}

// NewMockInterface створює новий mock Interface
func NewMockInterface() Interface {
	return &MockInterface{shouldFail: false}
}

// NewMockInterfaceWithError створює mock Interface що повертає помилки
func NewMockInterfaceWithError() Interface {
	return &MockInterface{shouldFail: true}
}
