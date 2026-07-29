-- Add asset_item_id column to asset_transactions table
-- This allows tracking specific physical asset items in transactions
-- 
-- IMPORTANT: This ADDS a new column, does NOT replace the existing asset_id column
-- - asset_id: Links to product catalog (assets table) - "what type of device"
-- - asset_item_id: Links to physical inventory (asset_items table) - "which specific device"

-- Add the column
ALTER TABLE `asset_transactions` 
ADD COLUMN `asset_item_id` VARCHAR(191) COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'Links to specific asset item in asset_items table' AFTER `customer_installation_id`;

-- Add index for performance
ALTER TABLE `asset_transactions` 
ADD INDEX `idx_asset_transactions_asset_item_id` (`asset_item_id`);

-- Add foreign key constraint
ALTER TABLE `asset_transactions` 
ADD CONSTRAINT `asset_transactions_asset_item_id_fk` 
FOREIGN KEY (`asset_item_id`) REFERENCES `asset_items` (`id`) 
ON DELETE CASCADE ON UPDATE CASCADE;

-- Add comment to the table
ALTER TABLE `asset_transactions` 
COMMENT = 'Tracks asset transactions for installations and trouble tickets. Links to specific asset items for precise inventory tracking.';
