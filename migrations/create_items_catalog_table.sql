-- Migration: Create items catalog table and update relationships
-- This creates a proper item catalog system for inventory management

-- Create items catalog table (master data for all items)
CREATE TABLE IF NOT EXISTS `items` (
  `id` VARCHAR(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` VARCHAR(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Item name (e.g., Drop 1C, Patchcore, Cvt single A, Kabel lan)',
  `default_unit` VARCHAR(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'PCS' COMMENT 'Default unit (PCS, M, KG, BOX, ROLL)',
  `category` VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT 'Item category (e.g., Cable, Connector, Hardware)',
  `description` TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT 'Item description',
  `asset_id` VARCHAR(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT 'Optional: Link to asset if this item is tracked as an asset',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`),
  KEY `idx_category` (`category`),
  KEY `idx_asset_id` (`asset_id`),
  CONSTRAINT `fk_items_asset` FOREIGN KEY (`asset_id`) REFERENCES `assets` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Item catalog - master data for all inventory items';

-- Insert sample items based on user requirements
INSERT INTO `items` (`id`, `name`, `default_unit`, `category`, `description`) VALUES
  (UUID(), 'Drop 1C', 'M', 'Cable', 'Drop cable 1 core'),
  (UUID(), 'Patchcore', 'PCS', 'Connector', 'Patchcore connector'),
  (UUID(), 'Cvt single A', 'PCS', 'Connector', 'Cvt single type A connector'),
  (UUID(), 'Kabel lan', 'M', 'Cable', 'LAN cable (UTP/STP)')
ON DUPLICATE KEY UPDATE `name` = `name`;

-- Note: Router items should be linked to assets table
-- We'll handle routers separately since they're tracked in asset management

-- Update detail_itemskeluar to reference items instead of assets
-- First, check if we need to add item_id column
SET @dbname = DATABASE();
SET @tablename = "detail_itemskeluar";
SET @columnname = "item_id";

SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE
      (table_name = @tablename)
      AND (table_schema = @dbname)
      AND (column_name = @columnname)
  ) > 0,
  "SELECT 1",
  CONCAT("ALTER TABLE ", @tablename, " ADD COLUMN `", @columnname, "` VARCHAR(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL AFTER `IdItems`, ADD KEY `idx_item_id` (`", @columnname, "`)")
));
PREPARE alterIfNotExists FROM @preparedStatement;
EXECUTE alterIfNotExists;
DEALLOCATE PREPARE alterIfNotExists;

-- Add foreign key for item_id if it doesn't exist
SET @constraintname = "FK_detail_itemskeluar_items";
SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS
    WHERE
      (table_name = @tablename)
      AND (table_schema = @dbname)
      AND (constraint_name = @constraintname)
      AND (constraint_type = 'FOREIGN KEY')
  ) > 0,
  "SELECT 1",
  CONCAT("ALTER TABLE ", @tablename, " ADD CONSTRAINT `", @constraintname, "` FOREIGN KEY (`item_id`) REFERENCES `items` (`id`) ON DELETE SET NULL ON UPDATE CASCADE")
));
PREPARE addIfNotExists FROM @preparedStatement;
EXECUTE addIfNotExists;
DEALLOCATE PREPARE addIfNotExists;

-- Update detail_itemsmasuk to reference items instead of assets
SET @tablename = "detail_itemsmasuk";
SET @columnname = "item_id";

SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE
      (table_name = @tablename)
      AND (table_schema = @dbname)
      AND (column_name = @columnname)
  ) > 0,
  "SELECT 1",
  CONCAT("ALTER TABLE ", @tablename, " ADD COLUMN `", @columnname, "` VARCHAR(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL AFTER `IdItems`, ADD KEY `idx_item_id` (`", @columnname, "`)")
));
PREPARE alterIfNotExists FROM @preparedStatement;
EXECUTE alterIfNotExists;
DEALLOCATE PREPARE alterIfNotExists;

-- Add foreign key for item_id if it doesn't exist
SET @constraintname = "FK_detail_itemsmasuk_items";
SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS
    WHERE
      (table_name = @tablename)
      AND (table_schema = @dbname)
      AND (constraint_name = @constraintname)
      AND (constraint_type = 'FOREIGN KEY')
  ) > 0,
  "SELECT 1",
  CONCAT("ALTER TABLE ", @tablename, " ADD CONSTRAINT `", @constraintname, "` FOREIGN KEY (`item_id`) REFERENCES `items` (`id`) ON DELETE SET NULL ON UPDATE CASCADE")
));
PREPARE addIfNotExists FROM @preparedStatement;
EXECUTE addIfNotExists;
DEALLOCATE PREPARE addIfNotExists;

