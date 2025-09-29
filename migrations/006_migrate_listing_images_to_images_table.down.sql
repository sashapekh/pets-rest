-- Rollback migration: Add images column back to listings and restore data from images table
-- Add the images column back
ALTER TABLE listings ADD COLUMN images TEXT[];

-- Migrate data back from images table to listings
UPDATE listings 
SET images = (
    SELECT array_agg(url ORDER BY sort_order)
    FROM images 
    WHERE imageable_type = 'listing' AND imageable_id = listings.id
)
WHERE id IN (
    SELECT DISTINCT imageable_id 
    FROM images 
    WHERE imageable_type = 'listing'
);

-- Delete listing images from images table
DELETE FROM images WHERE imageable_type = 'listing';
