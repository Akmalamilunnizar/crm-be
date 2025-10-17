-- Remove the old 'classification' column from trouble_tickets table
-- This column conflicts with our new classification_id system

-- Check if the old classification column exists and remove it
SET @col_exists = (
    SELECT COUNT(*) 
    FROM INFORMATION_SCHEMA.COLUMNS 
    WHERE TABLE_SCHEMA = DATABASE() 
    AND TABLE_NAME = 'trouble_tickets' 
    AND COLUMN_NAME = 'classification'
);

-- Remove the old classification column if it exists
SET @sql = IF(@col_exists > 0, 
    'ALTER TABLE trouble_tickets DROP COLUMN classification', 
    'SELECT "Column classification does not exist, skipping removal" as message'
);

PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Verify the removal
SELECT 
    COLUMN_NAME,
    DATA_TYPE,
    IS_NULLABLE,
    COLUMN_DEFAULT
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_SCHEMA = DATABASE() 
AND TABLE_NAME = 'trouble_tickets' 
AND COLUMN_NAME LIKE '%classification%'
ORDER BY ORDINAL_POSITION;
