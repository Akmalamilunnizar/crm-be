-- Fix foreign key constraints for items catalog
-- This script ensures item_id columns match the items.id column type exactly

SET @dbname = DATABASE();

-- Drop existing foreign key if it exists (in case of type mismatch)
SET @tablename = "detail_itemskeluar";
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
  CONCAT("ALTER TABLE ", @tablename, " DROP FOREIGN KEY `", @constraintname, "`"),
  "SELECT 1"
));
PREPARE dropIfExists FROM @preparedStatement;
EXECUTE dropIfExists;
DEALLOCATE PREPARE dropIfExists;

-- Ensure item_id column exists and has correct type (VARCHAR(191) with utf8mb4_unicode_ci)
-- Check column type first
SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE
      (table_name = @tablename)
      AND (table_schema = @dbname)
      AND (column_name = 'item_id')
      AND (column_type LIKE 'varchar(191)%')
      AND (collation_name = 'utf8mb4_unicode_ci')
  ) > 0,
  "SELECT 1",
  CONCAT("ALTER TABLE ", @tablename, " MODIFY COLUMN `item_id` VARCHAR(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL")
));
PREPARE modifyColumn FROM @preparedStatement;
EXECUTE modifyColumn;
DEALLOCATE PREPARE modifyColumn;

-- Add foreign key constraint (now that types match)
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

-- Same for detail_itemsmasuk
SET @tablename = "detail_itemsmasuk";
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
  CONCAT("ALTER TABLE ", @tablename, " DROP FOREIGN KEY `", @constraintname, "`"),
  "SELECT 1"
));
PREPARE dropIfExists FROM @preparedStatement;
EXECUTE dropIfExists;
DEALLOCATE PREPARE dropIfExists;

-- Ensure item_id column has correct type
SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE
      (table_name = @tablename)
      AND (table_schema = @dbname)
      AND (column_name = 'item_id')
      AND (column_type LIKE 'varchar(191)%')
      AND (collation_name = 'utf8mb4_unicode_ci')
  ) > 0,
  "SELECT 1",
  CONCAT("ALTER TABLE ", @tablename, " MODIFY COLUMN `item_id` VARCHAR(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL")
));
PREPARE modifyColumn FROM @preparedStatement;
EXECUTE modifyColumn;
DEALLOCATE PREPARE modifyColumn;

-- Add foreign key constraint
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

