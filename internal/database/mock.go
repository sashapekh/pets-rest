package database

import "errors"

// MockDB - mock реалізація для тестування
type MockDB struct{}

// Health перевіряє стан mock бази даних
func (m *MockDB) Health() error {
	return nil
}

// Close закриває mock з'єднання
func (m *MockDB) Close() error {
	return nil
}

// MockDBWithError - mock з помилкою для тестування
type MockDBWithError struct{}

// Health повертає помилку для тестування
func (m *MockDBWithError) Health() error {
	return errors.New("database connection failed")
}

// Close закриває mock з'єднання з помилкою
func (m *MockDBWithError) Close() error {
	return nil
}
