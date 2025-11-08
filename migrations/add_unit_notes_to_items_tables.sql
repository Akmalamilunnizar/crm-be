-- Migration: Add unit and notes columns to detail_itemskeluar and detail_itemsmasuk tables
-- Run this migration to add the new columns if they don't exist

-- Note: This script uses prepared statements to check if columns exist before adding them
-- because MySQL doesn't support IF NOT EXISTS for ALTER TABLE ADD COLUMN

SET @dbname = DATABASE();

-- Add unit and notes to detail_itemskeluar
SET @tablename = "detail_itemskeluar";
SET @columnname = "unit";
SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE
      (table_name = @tablename)
      AND (table_schema = @dbname)
      AND (column_name = @columnname)
  ) > 0,
  "SELECT 1",
  CONCAT("ALTER TABLE ", @tablename, " ADD COLUMN `", @columnname, "` VARCHAR(20) DEFAULT 'PCS' AFTER `QtyKeluar`")
));
PREPARE alterIfNotExists FROM @preparedStatement;
EXECUTE alterIfNotExists;
DEALLOCATE PREPARE alterIfNotExists;

SET @columnname = "notes";
SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE
      (table_name = @tablename)
      AND (table_schema = @dbname)
      AND (column_name = @columnname)
  ) > 0,
  "SELECT 1",
  CONCAT("ALTER TABLE ", @tablename, " ADD COLUMN `", @columnname, "` TEXT NULL AFTER `unit`")
));
PREPARE alterIfNotExists FROM @preparedStatement;
EXECUTE alterIfNotExists;
DEALLOCATE PREPARE alterIfNotExists;

-- Add unit and notes to detail_itemsmasuk
SET @tablename = "detail_itemsmasuk";
SET @columnname = "unit";
SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE
      (table_name = @tablename)
      AND (table_schema = @dbname)
      AND (column_name = @columnname)
  ) > 0,
  "SELECT 1",
  CONCAT("ALTER TABLE ", @tablename, " ADD COLUMN `", @columnname, "` VARCHAR(20) DEFAULT 'PCS' AFTER `SubTotal`")
));
PREPARE alterIfNotExists FROM @preparedStatement;
EXECUTE alterIfNotExists;
DEALLOCATE PREPARE alterIfNotExists;

SET @columnname = "notes";
SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE
      (table_name = @tablename)
      AND (table_schema = @dbname)
      AND (column_name = @columnname)
  ) > 0,
  "SELECT 1",
  CONCAT("ALTER TABLE ", @tablename, " ADD COLUMN `", @columnname, "` TEXT NULL AFTER `unit`")
));
PREPARE alterIfNotExists FROM @preparedStatement;
EXECUTE alterIfNotExists;
DEALLOCATE PREPARE alterIfNotExists;

-- Create parent tables if they don't exist (itemskeluar and itemsmasuk)
CREATE TABLE IF NOT EXISTS `itemskeluar` (
  `id` VARCHAR(6) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `date` DATE NOT NULL,
  `notes` TEXT NULL,
  `created_by` VARCHAR(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_created_by` (`created_by`),
  KEY `idx_date` (`date`),
  CONSTRAINT `fk_itemskeluar_user` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `itemsmasuk` (
  `id` VARCHAR(6) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `date` DATE NOT NULL,
  `notes` TEXT NULL,
  `created_by` VARCHAR(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_created_by` (`created_by`),
  KEY `idx_date` (`date`),
  CONSTRAINT `fk_itemsmasuk_user` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Update foreign keys in detail tables to reference itemskeluar and itemsmasuk
-- Drop old foreign keys if they exist and reference different tables
SET @tablename = "detail_itemskeluar";
SET @constraintname = "FK_detail_assetkeluar_assetkeluar";
SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS
    WHERE
      (table_name = @tablename)
      AND (table_schema = @dbname)
      AND (constraint_name = @constraintname)
      AND (constraint_type = 'FOREIGN KEY')
  ) > 0,
  CONCAT("ALTER TABLE ", @tablename, " DROP FOREIGN KEY `", @constraintname, "`"),
  "SELECT 1"
));
PREPARE dropIfExists FROM @preparedStatement;
EXECUTE dropIfExists;
DEALLOCATE PREPARE dropIfExists;

SET @tablename = "detail_itemsmasuk";
SET @constraintname = "FK_detail_assetmasuk_assetmasuk";
SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS
    WHERE
      (table_name = @tablename)
      AND (table_schema = @dbname)
      AND (constraint_name = @constraintname)
      AND (constraint_type = 'FOREIGN KEY')
  ) > 0,
  CONCAT("ALTER TABLE ", @tablename, " DROP FOREIGN KEY `", @constraintname, "`"),
  "SELECT 1"
));
PREPARE dropIfExists FROM @preparedStatement;
EXECUTE dropIfExists;
DEALLOCATE PREPARE dropIfExists;

-- Add new foreign keys to itemskeluar and itemsmasuk (only if they don't exist)
SET @tablename = "detail_itemskeluar";
SET @constraintname = "FK_detail_itemskeluar_itemskeluar";
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
  CONCAT("ALTER TABLE ", @tablename, " ADD CONSTRAINT `", @constraintname, "` FOREIGN KEY (`IdKeluar`) REFERENCES `itemskeluar` (`id`) ON DELETE CASCADE ON UPDATE CASCADE")
));
PREPARE addIfNotExists FROM @preparedStatement;
EXECUTE addIfNotExists;
DEALLOCATE PREPARE addIfNotExists;

SET @tablename = "detail_itemsmasuk";
SET @constraintname = "FK_detail_itemsmasuk_itemsmasuk";
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
  CONCAT("ALTER TABLE ", @tablename, " ADD CONSTRAINT `", @constraintname, "` FOREIGN KEY (`IdMasuk`) REFERENCES `itemsmasuk` (`id`) ON DELETE CASCADE ON UPDATE CASCADE")
));
PREPARE addIfNotExists FROM @preparedStatement;
EXECUTE addIfNotExists;
DEALLOCATE PREPARE addIfNotExists;
