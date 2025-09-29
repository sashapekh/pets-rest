-- Create images table with polymorphic relationship
CREATE TABLE IF NOT EXISTS images (
    id SERIAL PRIMARY KEY,
    imageable_type VARCHAR(50) NOT NULL,
    imageable_id INTEGER NOT NULL,
    url TEXT NOT NULL,
    filename VARCHAR(255),
    size_bytes INTEGER,
    mime_type VARCHAR(100),
    alt_text TEXT,
    sort_order INTEGER DEFAULT 0,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for polymorphic queries
CREATE INDEX IF NOT EXISTS idx_images_polymorphic ON images(imageable_type, imageable_id);
CREATE INDEX IF NOT EXISTS idx_images_sort ON images(imageable_type, imageable_id, sort_order);
CREATE INDEX IF NOT EXISTS idx_images_primary ON images(imageable_type, imageable_id, is_primary) WHERE is_primary = TRUE;
CREATE INDEX IF NOT EXISTS idx_images_created_at ON images(created_at DESC);

-- Add constraint to ensure valid imageable_type
ALTER TABLE images ADD CONSTRAINT check_imageable_type 
CHECK (imageable_type IN ('user', 'listing'));

-- Add trigger for updated_at timestamps on images
CREATE TRIGGER update_images_updated_at BEFORE UPDATE ON images
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add constraint to ensure only one primary image per entity
CREATE UNIQUE INDEX idx_images_unique_primary 
ON images(imageable_type, imageable_id) 
WHERE is_primary = TRUE;
