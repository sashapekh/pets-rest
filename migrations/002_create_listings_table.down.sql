-- Drop trigger
DROP TRIGGER IF EXISTS update_listings_updated_at ON listings;

-- Drop indexes
DROP INDEX IF EXISTS idx_listings_created_at;
DROP INDEX IF EXISTS idx_listings_city;
DROP INDEX IF EXISTS idx_listings_slug;
DROP INDEX IF EXISTS idx_listings_status;
DROP INDEX IF EXISTS idx_listings_type;
DROP INDEX IF EXISTS idx_listings_user_id;

-- Drop table
DROP TABLE IF EXISTS listings;
