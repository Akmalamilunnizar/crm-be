-- Check database schema for customer_installations table

-- 1. Check table structure
DESCRIBE customer_installations;

-- 2. Check if document_photo column exists and its type
SELECT 
    COLUMN_NAME,
    DATA_TYPE,
    IS_NULLABLE,
    COLUMN_DEFAULT,
    COLUMN_COMMENT
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_SCHEMA = 'iqgncnzy_skripsi' 
AND TABLE_NAME = 'customer_installations'
AND COLUMN_NAME = 'document_photo';

-- 3. Check recent installations
SELECT 
    id,
    customer_id,
    technician_id,
    document_type,
    document_photo,
    status,
    createdAt,
    updatedAt
FROM customer_installations 
ORDER BY createdAt DESC 
LIMIT 5;

-- 4. Check if document_photo column has any data
SELECT 
    COUNT(*) as total_records,
    COUNT(document_photo) as records_with_document_photo,
    COUNT(CASE WHEN document_photo IS NOT NULL AND document_photo != '' THEN 1 END) as non_empty_document_photos
FROM customer_installations;
