-- Diagnostic script to verify database schema and column existence
-- Run this to check if latitude/longitude columns exist in customer_installations table

-- Check current database
SELECT DATABASE() as current_database;

-- Check if customer_installations table exists
SELECT 
    TABLE_NAME, 
    TABLE_COMMENT 
FROM 
    INFORMATION_SCHEMA.TABLES 
WHERE 
    TABLE_SCHEMA = DATABASE() 
    AND TABLE_NAME = 'customer_installations';

-- Check if latitude and longitude columns exist
SELECT 
    COLUMN_NAME, 
    DATA_TYPE, 
    COLUMN_DEFAULT, 
    IS_NULLABLE,
    COLUMN_COMMENT
FROM 
    INFORMATION_SCHEMA.COLUMNS 
WHERE 
    TABLE_SCHEMA = DATABASE() 
    AND TABLE_NAME = 'customer_installations' 
    AND COLUMN_NAME IN ('latitude', 'longitude')
ORDER BY 
    COLUMN_NAME;

-- If columns exist, check some sample data
SELECT 
    id,
    customer_id,
    latitude,
    longitude,
    createdAt
FROM 
    customer_installations 
LIMIT 3;
