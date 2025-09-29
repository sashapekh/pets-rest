-- Migrate existing listing images from images array to the new images table
-- First, create a function to migrate the data
CREATE OR REPLACE FUNCTION migrate_listing_images() RETURNS void AS $$
DECLARE
    listing_record RECORD;
    image_url TEXT;
    sort_order_counter INTEGER;
BEGIN
    -- Loop through all listings that have images
    FOR listing_record IN 
        SELECT id, images FROM listings WHERE images IS NOT NULL AND array_length(images, 1) > 0
    LOOP
        sort_order_counter := 0;
        
        -- Loop through each image URL in the array
        FOREACH image_url IN ARRAY listing_record.images
        LOOP
            -- Insert into images table
            INSERT INTO images (
                imageable_type,
                imageable_id, 
                url,
                sort_order,
                is_primary,
                created_at
            ) VALUES (
                'listing',
                listing_record.id,
                image_url,
                sort_order_counter,
                (sort_order_counter = 0), -- First image is primary
                NOW()
            );
            
            sort_order_counter := sort_order_counter + 1;
        END LOOP;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- Execute the migration function
SELECT migrate_listing_images();

-- Drop the migration function as it's no longer needed
DROP FUNCTION migrate_listing_images();

-- Remove the images column from listings table
ALTER TABLE listings DROP COLUMN IF EXISTS images;
