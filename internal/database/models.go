package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// User represents a user in the system
type User struct {
	ID        int        `json:"id" db:"id"`
	Email     string     `json:"email" db:"email"`
	Phone     *string    `json:"phone,omitempty" db:"phone"`
	Name      *string    `json:"name,omitempty" db:"name"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// ListingType represents the type of listing
type ListingType string

const (
	ListingTypeLost  ListingType = "lost"
	ListingTypeFound ListingType = "found"
	ListingTypeAdopt ListingType = "adopt"
)

// ListingStatus represents the status of listing
type ListingStatus string

const (
	ListingStatusDraft    ListingStatus = "draft"
	ListingStatusActive   ListingStatus = "active"
	ListingStatusArchived ListingStatus = "archived"
)

// Listing represents a pet listing
type Listing struct {
	ID           int            `json:"id" db:"id"`
	UserID       int            `json:"user_id" db:"user_id"`
	Type         ListingType    `json:"type" db:"type"`
	Title        string         `json:"title" db:"title"`
	Description  *string        `json:"description,omitempty" db:"description"`
	City         *string        `json:"city,omitempty" db:"city"`
	Location     *string        `json:"location,omitempty" db:"location"`
	ContactPhone *string        `json:"contact_phone,omitempty" db:"contact_phone"`
	ContactTg    *string        `json:"contact_tg,omitempty" db:"contact_tg"`
	Status       ListingStatus  `json:"status" db:"status"`
	Slug         *string        `json:"slug,omitempty" db:"slug"`
	Images       pq.StringArray `json:"images" db:"images"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt    *time.Time     `json:"updated_at,omitempty" db:"updated_at"`
}

// EventType represents the type of event
type EventType string

const (
	EventTypeView         EventType = "view"
	EventTypeQRScan       EventType = "qr_scan"
	EventTypeContactClick EventType = "contact_click"
	EventTypePhoneClick   EventType = "phone_click"
)

// JSONPayload is a custom type for handling JSONB
type JSONPayload map[string]interface{}

// Value implements the driver.Valuer interface for JSONPayload
func (j JSONPayload) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for JSONPayload
func (j *JSONPayload) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into JSONPayload", value)
	}

	return json.Unmarshal(bytes, j)
}

// Event represents an analytics event
type Event struct {
	ID        int         `json:"id" db:"id"`
	UserID    *int        `json:"user_id,omitempty" db:"user_id"`
	ListingID int         `json:"listing_id" db:"listing_id"`
	Type      EventType   `json:"type" db:"type"`
	Payload   JSONPayload `json:"payload,omitempty" db:"payload"`
	IPAddress *string     `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent *string     `json:"user_agent,omitempty" db:"user_agent"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
}
