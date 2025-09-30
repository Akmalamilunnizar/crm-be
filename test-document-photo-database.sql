-- Test document photo in database

-- 1. Check customer_installations table structure
DESCRIBE customer_installations;

-- 2. Check recent installations with document photos
SELECT 
    id,
    customer_id,
    technician_id,
    document_type,
    document_photo,
    installation_type,
    status,
    createdAt,
    updatedAt
FROM customer_installations 
ORDER BY createdAt DESC 
LIMIT 5;

-- 3. Check if document_photo column has data
SELECT 
    COUNT(*) as total_installations,
    COUNT(document_photo) as installations_with_document_photo,
    COUNT(CASE WHEN document_photo IS NOT NULL AND document_photo != '' THEN 1 END) as non_empty_document_photos
FROM customer_installations;

-- 4. Show installations with document photos
SELECT 
    id,
    document_type,
    document_photo,
    LENGTH(document_photo) as photo_path_length,
    createdAt
FROM customer_installations 
WHERE document_photo IS NOT NULL 
AND document_photo != ''
ORDER BY createdAt DESC;

-- 5. Check uploads directory structure (if possible)
-- This would need to be run from file system
-- SELECT 'uploads/installations/documents/' as upload_directory;
