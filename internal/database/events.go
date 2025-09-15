package database

import (
	"time"
)

// EventRepository handles event database operations
type EventRepository struct {
	db *DB
}

// NewEventRepository creates a new event repository
func NewEventRepository(db *DB) *EventRepository {
	return &EventRepository{db: db}
}

// Create creates a new event
func (r *EventRepository) Create(event *Event) error {
	query := `
		INSERT INTO events (user_id, listing_id, type, payload, ip_address, user_agent, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`

	err := r.db.QueryRow(query,
		event.UserID,
		event.ListingID,
		event.Type,
		event.Payload,
		event.IPAddress,
		event.UserAgent,
		time.Now()).
		Scan(&event.ID, &event.CreatedAt)

	return err
}

// GetByListingID retrieves events for a specific listing
func (r *EventRepository) GetByListingID(listingID, limit, offset int) ([]*Event, error) {
	events := []*Event{}
	query := `
		SELECT id, user_id, listing_id, type, payload, ip_address, user_agent, created_at 
		FROM events 
		WHERE listing_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`

	err := r.db.Select(&events, query, listingID, limit, offset)
	return events, err
}

// GetByUserID retrieves events for a specific user
func (r *EventRepository) GetByUserID(userID, limit, offset int) ([]*Event, error) {
	events := []*Event{}
	query := `
		SELECT id, user_id, listing_id, type, payload, ip_address, user_agent, created_at 
		FROM events 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`

	err := r.db.Select(&events, query, userID, limit, offset)
	return events, err
}

// GetAnalytics returns analytics data for a listing
func (r *EventRepository) GetAnalytics(listingID int) (map[string]int, error) {
	analytics := make(map[string]int)

	query := `
		SELECT type, COUNT(*) as count 
		FROM events 
		WHERE listing_id = $1 
		GROUP BY type`

	rows, err := r.db.Query(query, listingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var eventType string
		var count int

		if err := rows.Scan(&eventType, &count); err != nil {
			return nil, err
		}

		analytics[eventType] = count
	}

	return analytics, nil
}

// GetDailyAnalytics returns daily analytics for the last N days
func (r *EventRepository) GetDailyAnalytics(listingID, days int) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			DATE(created_at) as date,
			type,
			COUNT(*) as count
		FROM events 
		WHERE listing_id = $1 
		AND created_at >= NOW() - INTERVAL '%d days'
		GROUP BY DATE(created_at), type
		ORDER BY date DESC, type`

	rows, err := r.db.Query(query, listingID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analytics []map[string]interface{}

	for rows.Next() {
		var date time.Time
		var eventType string
		var count int

		if err := rows.Scan(&date, &eventType, &count); err != nil {
			return nil, err
		}

		analytics = append(analytics, map[string]interface{}{
			"date":  date.Format("2006-01-02"),
			"type":  eventType,
			"count": count,
		})
	}

	return analytics, nil
}

// CountByListing returns the total number of events for a listing
func (r *EventRepository) CountByListing(listingID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM events WHERE listing_id = $1`

	err := r.db.Get(&count, query, listingID)
	return count, err
}

// CountByUser returns the total number of events for a user
func (r *EventRepository) CountByUser(userID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM events WHERE user_id = $1`

	err := r.db.Get(&count, query, userID)
	return count, err
}
