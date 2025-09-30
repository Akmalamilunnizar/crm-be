-- Fix database relation between customer_installations and images tables

-- First, let's check if there are any orphaned images (images without valid installation IDs)
SELECT 
    i.id as image_id,
    i.archive_installation_id,
    i.file,
    i.full_path,
    ci.id as installation_id
FROM images i
LEFT JOIN customer_installations ci ON i.archive_installation_id = ci.id
WHERE ci.id IS NULL;

-- Add foreign key constraint to images table
-- This will ensure that archive_installation_id references a valid customer_installations.id
ALTER TABLE images 
ADD CONSTRAINT fk_images_customer_installations 
FOREIGN KEY (archive_installation_id) 
REFERENCES customer_installations(id) 
ON DELETE CASCADE 
ON UPDATE CASCADE;

-- Create index for better performance
CREATE INDEX idx_images_archive_installation_id ON images(archive_installation_id);

-- Verify the constraint was added
SELECT 
    CONSTRAINT_NAME,
    TABLE_NAME,
    COLUMN_NAME,
    REFERENCED_TABLE_NAME,
    REFERENCED_COLUMN_NAME
FROM information_schema.KEY_COLUMN_USAGE 
WHERE TABLE_SCHEMA = 'iqgncnzy_skripsi' 
AND TABLE_NAME = 'images' 
AND REFERENCED_TABLE_NAME IS NOT NULL;


