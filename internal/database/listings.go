package database

import (
	"fmt"
	"time"

	"github.com/lib/pq"
)

// ListingRepository handles listing database operations
type ListingRepository struct {
	db *DB
}

// NewListingRepository creates a new listing repository
func NewListingRepository(db *DB) *ListingRepository {
	return &ListingRepository{db: db}
}

// Create creates a new listing
func (r *ListingRepository) Create(listing *Listing) error {
	query := `
		INSERT INTO listings (user_id, type, title, description, city, location, contact_phone, contact_tg, status, slug, images, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at`

	err := r.db.QueryRow(query,
		listing.UserID,
		listing.Type,
		listing.Title,
		listing.Description,
		listing.City,
		listing.Location,
		listing.ContactPhone,
		listing.ContactTg,
		listing.Status,
		listing.Slug,
		pq.Array(listing.Images),
		time.Now()).
		Scan(&listing.ID, &listing.CreatedAt)

	return err
}

// GetByID retrieves a listing by ID
func (r *ListingRepository) GetByID(id int) (*Listing, error) {
	listing := &Listing{}
	query := `
		SELECT id, user_id, type, title, description, city, location, contact_phone, contact_tg, status, slug, images, created_at, updated_at 
		FROM listings 
		WHERE id = $1`

	err := r.db.Get(listing, query, id)
	if err != nil {
		return nil, err
	}

	return listing, nil
}

// GetBySlug retrieves a listing by slug
func (r *ListingRepository) GetBySlug(slug string) (*Listing, error) {
	listing := &Listing{}
	query := `
		SELECT id, user_id, type, title, description, city, location, contact_phone, contact_tg, status, slug, images, created_at, updated_at 
		FROM listings 
		WHERE slug = $1 AND status = 'active'`

	err := r.db.Get(listing, query, slug)
	if err != nil {
		return nil, err
	}

	return listing, nil
}

// Update updates a listing
func (r *ListingRepository) Update(listing *Listing) error {
	query := `
		UPDATE listings 
		SET type = $2, title = $3, description = $4, city = $5, location = $6, contact_phone = $7, contact_tg = $8, status = $9, slug = $10, images = $11, updated_at = $12
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRow(query,
		listing.ID,
		listing.Type,
		listing.Title,
		listing.Description,
		listing.City,
		listing.Location,
		listing.ContactPhone,
		listing.ContactTg,
		listing.Status,
		listing.Slug,
		pq.Array(listing.Images),
		time.Now()).
		Scan(&listing.UpdatedAt)

	return err
}

// Delete deletes a listing
func (r *ListingRepository) Delete(id int) error {
	query := `DELETE FROM listings WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("listing with id %d not found", id)
	}

	return nil
}

// ListByUser retrieves all listings for a user
func (r *ListingRepository) ListByUser(userID int, limit, offset int) ([]*Listing, error) {
	listings := []*Listing{}
	query := `
		SELECT id, user_id, type, title, description, city, location, contact_phone, contact_tg, status, slug, images, created_at, updated_at 
		FROM listings 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`

	err := r.db.Select(&listings, query, userID, limit, offset)
	return listings, err
}

// ListActive retrieves all active listings with optional filters
func (r *ListingRepository) ListActive(listingType *ListingType, city *string, limit, offset int) ([]*Listing, error) {
	listings := []*Listing{}

	baseQuery := `
		SELECT id, user_id, type, title, description, city, location, contact_phone, contact_tg, status, slug, images, created_at, updated_at 
		FROM listings 
		WHERE status = 'active'`

	args := []interface{}{}
	argCount := 0

	if listingType != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND type = $%d", argCount)
		args = append(args, *listingType)
	}

	if city != nil {
		argCount++
		baseQuery += fmt.Sprintf(" AND city ILIKE $%d", argCount)
		args = append(args, "%"+*city+"%")
	}

	baseQuery += " ORDER BY created_at DESC"

	if limit > 0 {
		argCount++
		baseQuery += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, limit)

		if offset > 0 {
			argCount++
			baseQuery += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	err := r.db.Select(&listings, baseQuery, args...)
	return listings, err
}

// Search searches listings by title and description
func (r *ListingRepository) Search(query string, limit, offset int) ([]*Listing, error) {
	listings := []*Listing{}
	searchQuery := `
		SELECT id, user_id, type, title, description, city, location, contact_phone, contact_tg, status, slug, images, created_at, updated_at 
		FROM listings 
		WHERE status = 'active' 
		AND (title ILIKE $1 OR description ILIKE $1)
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`

	searchTerm := "%" + query + "%"
	err := r.db.Select(&listings, searchQuery, searchTerm, limit, offset)
	return listings, err
}

// Count returns the total number of listings
func (r *ListingRepository) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM listings`

	err := r.db.Get(&count, query)
	return count, err
}

// CountByUser returns the total number of listings for a user
func (r *ListingRepository) CountByUser(userID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM listings WHERE user_id = $1`

	err := r.db.Get(&count, query, userID)
	return count, err
}
