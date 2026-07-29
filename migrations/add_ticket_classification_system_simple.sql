-- Simplified classification system migration (MySQL 5.7+ compatible)
-- This version works without foreign keys and uses older MySQL syntax

-- Step 1: Create ticket_classification lookup table
CREATE TABLE IF NOT EXISTS ticket_classification (
    id VARCHAR(20) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    show_noc_action TINYINT(1) DEFAULT 1 COMMENT 'Whether to show NOC action button',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Step 2: Insert classification types
INSERT INTO ticket_classification (id, name, description, show_noc_action) VALUES
('gangguan', 'Gangguan', 'Trouble/issue tickets - customer experiencing problems', 1),
('psb', 'PSB (Pasang Baru)', 'New installation requests', 0),
('lainnya', 'Lainnya', 'Other miscellaneous tickets', 0),
('dismantle', 'Dismantle', 'Service termination and equipment removal', 0)
ON DUPLICATE KEY UPDATE 
    name = VALUES(name),
    description = VALUES(description),
    show_noc_action = VALUES(show_noc_action);

-- Step 3: Ensure classification_id column has correct default
ALTER TABLE trouble_tickets 
MODIFY COLUMN classification_id VARCHAR(20) NULL DEFAULT 'gangguan';

-- Step 4: Update all NULL or empty values to 'gangguan'
UPDATE trouble_tickets 
SET classification_id = 'gangguan' 
WHERE classification_id IS NULL OR classification_id = '';

-- Step 5: Make it NOT NULL
ALTER TABLE trouble_tickets 
MODIFY COLUMN classification_id VARCHAR(20) NOT NULL DEFAULT 'gangguan';

-- Step 6: Remove old verified_by_cs column (with error handling)
SET @column_exists = (
    SELECT COUNT(*) 
    FROM information_schema.COLUMNS 
    WHERE TABLE_SCHEMA = DATABASE() 
    AND TABLE_NAME = 'trouble_tickets' 
    AND COLUMN_NAME = 'verified_by_cs'
);

SET @drop_sql = IF(@column_exists > 0, 
    'ALTER TABLE trouble_tickets DROP COLUMN verified_by_cs', 
    'SELECT "Column verified_by_cs does not exist" as Info');
PREPARE stmt FROM @drop_sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Step 7: Remove old verified_at column (with error handling)
SET @column_exists_at = (
    SELECT COUNT(*) 
    FROM information_schema.COLUMNS 
    WHERE TABLE_SCHEMA = DATABASE() 
    AND TABLE_NAME = 'trouble_tickets' 
    AND COLUMN_NAME = 'verified_at'
);

SET @drop_sql_at = IF(@column_exists_at > 0, 
    'ALTER TABLE trouble_tickets DROP COLUMN verified_at', 
    'SELECT "Column verified_at does not exist" as Info');
PREPARE stmt FROM @drop_sql_at;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Step 8: Clean up trouble_type table
-- Remove PSB and Dismantle (they are classifications, not trouble types)
DELETE FROM trouble_type WHERE id IN ('psb', 'dismantle', 'gangguan');

-- Step 9: Ensure "Lainnya" trouble type exists
INSERT INTO trouble_type (id, name) 
VALUES ('lainnya', 'Lainnya')
ON DUPLICATE KEY UPDATE name = 'Lainnya';

-- Note: Foreign key not added due to character set compatibility issues
-- Application logic enforces referential integrity
-- Valid classification_id values: 'gangguan', 'psb', 'lainnya', 'dismantle'

-- Verification queries:
SELECT 'Classification types:' as Info;
SELECT * FROM ticket_classification ORDER BY id;

SELECT 'Ticket distribution by classification:' as Info;
SELECT classification_id, COUNT(*) as count FROM trouble_tickets GROUP BY classification_id;

SELECT 'Available trouble types:' as Info;
SELECT * FROM trouble_type ORDER BY id;
