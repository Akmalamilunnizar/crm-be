-- Migration to add company_id to assets table
-- This allows assets to be associated with specific companies

-- Step 1: Add company_id column to assets table
ALTER TABLE `assets` 
ADD COLUMN `company_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL AFTER `serial_number`;

-- Step 2: Remove the old site column (replace with company_id)
ALTER TABLE `assets` 
DROP COLUMN `site`;

-- Step 3: Add foreign key constraint
ALTER TABLE `assets` 
ADD CONSTRAINT `assets_company_fk` 
FOREIGN KEY (`company_id`) REFERENCES `company` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- Step 4: Add index for better performance
ALTER TABLE `assets` 
ADD KEY `assets_company_idx` (`company_id`);

