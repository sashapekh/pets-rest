-- Drop indexes
DROP INDEX IF EXISTS idx_events_user_id;
DROP INDEX IF EXISTS idx_events_created_at;
DROP INDEX IF EXISTS idx_events_type;
DROP INDEX IF EXISTS idx_events_listing_id;

-- Drop table
DROP TABLE IF EXISTS events;
