-- Migration script to fix customer_installations table
-- This script adds missing ID column and other required fields to match the entity model

-- Step 1: Backup existing data (optional but recommended)
CREATE TABLE customer_installations_backup AS 
SELECT * FROM customer_installations;

-- Step 2: Drop the existing table (since it has no primary key and missing columns)
DROP TABLE IF EXISTS customer_installations;

-- Step 3: Create the new customer_installations table with proper structure
CREATE TABLE IF NOT EXISTS `customer_installations` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `customer_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Reference to customer table',
  `technician_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Reference to users table (technician)',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT 'Installation description',
  `date_trial` date NOT NULL COMMENT 'Trial installation date',
  `status` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'pending' COMMENT 'Installation status',
  `equipment_used` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT 'Equipment used in installation',
  `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT 'Installation notes',
  `completion_time` int DEFAULT 0 COMMENT 'Completion time in minutes',
  `customer_signature` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT 'Customer signature data',
  `on_air_date` date DEFAULT NULL COMMENT 'Date when service first went live/on-air',
  `createdAt` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updatedAt` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_customer_installations_customer_id` (`customer_id`),
  KEY `idx_customer_installations_technician_id` (`technician_id`),
  KEY `idx_customer_installations_on_air_date` (`on_air_date`),
  KEY `idx_customer_installations_status` (`status`),
  CONSTRAINT `fk_customer_installations_customer` FOREIGN KEY (`customer_id`) REFERENCES `customer` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_customer_installations_technician` FOREIGN KEY (`technician_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Step 4: Insert sample data (optional - you can remove this if not needed)
-- Note: You'll need to provide actual customer_id and technician_id values
INSERT INTO `customer_installations` (
  `id`, 
  `customer_id`, 
  `technician_id`, 
  `description`, 
  `date_trial`, 
  `status`, 
  `equipment_used`, 
  `notes`, 
  `completion_time`, 
  `on_air_date`, 
  `createdAt`, 
  `updatedAt`
) VALUES 
(
  UUID(), 
  '24225552-b0c3-43d7-9a67-20e04f36fa5f', 
  'c13a6c87-ec28-47ba-84c2-58b5ace2af57', 
  'Initial installation setup', 
  '2025-09-17', 
  'completed', 
  'Router, Cable, Modem', 
  'Installation completed successfully', 
  120, 
  '2025-09-17', 
  NOW(), 
  NOW()
),
(
  UUID(), 
  '55f03555-973d-425b-92bd-a0825c20b8ce', 
  'ba73859d-8957-418a-8311-28a6b31cad95', 
  'Standard installation', 
  '2025-09-17', 
  'pending', 
  'Router, Cable', 
  'Installation in progress', 
  0, 
  NULL, 
  NOW(), 
  NOW()
);

-- Step 5: Clean up backup table (uncomment if you want to remove backup)
-- DROP TABLE IF EXISTS customer_installations_backup;

-- Verification query
SELECT 
  'customer_installations' as table_name,
  COUNT(*) as record_count,
  'Table structure updated successfully' as status
FROM customer_installations;
