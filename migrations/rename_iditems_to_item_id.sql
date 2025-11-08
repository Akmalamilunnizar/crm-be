-- Migration: Rename IdItems to item_id and update foreign keys to reference items catalog
-- This approach uses the existing IdItems column but changes its purpose to reference items catalog

SET @dbname = DATABASE();

-- Step 1: Drop existing foreign keys that reference assets table
-- For detail_itemskeluar
SET @tablename = "detail_itemskeluar";
SET @constraintname = "FK_detail_barangkeluar_assets";
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

-- For detail_itemsmasuk
SET @tablename = "detail_itemsmasuk";
SET @constraintname = "fk_detail_barangmasuk_assets";
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

-- Step 2: Rename IdItems to item_id (if not already renamed)
-- For detail_itemskeluar
SET @tablename = "detail_itemskeluar";
SET @oldcolumn = "IdItems";
SET @newcolumn = "item_id";

SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE
      (table_name = @tablename)
      AND (table_schema = @dbname)
      AND (column_name = @newcolumn)
  ) > 0,
  "SELECT 1",
  CONCAT("ALTER TABLE ", @tablename, " CHANGE COLUMN `", @oldcolumn, "` `", @newcolumn, "` VARCHAR(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL")
));
PREPARE renameIfNotExists FROM @preparedStatement;
EXECUTE renameIfNotExists;
DEALLOCATE PREPARE renameIfNotExists;

-- For detail_itemsmasuk
SET @tablename = "detail_itemsmasuk";
SET @oldcolumn = "IdItems";
SET @newcolumn = "item_id";

SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE
      (table_name = @tablename)
      AND (table_schema = @dbname)
      AND (column_name = @newcolumn)
  ) > 0,
  "SELECT 1",
  CONCAT("ALTER TABLE ", @tablename, " CHANGE COLUMN `", @oldcolumn, "` `", @newcolumn, "` VARCHAR(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL")
));
PREPARE renameIfNotExists FROM @preparedStatement;
EXECUTE renameIfNotExists;
DEALLOCATE PREPARE renameIfNotExists;

-- Step 3: Add foreign key to reference items catalog
-- For detail_itemskeluar
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
  "SELECT 1",
  CONCAT("ALTER TABLE ", @tablename, " ADD CONSTRAINT `", @constraintname, "` FOREIGN KEY (`item_id`) REFERENCES `items` (`id`) ON DELETE SET NULL ON UPDATE CASCADE")
));
PREPARE addIfNotExists FROM @preparedStatement;
EXECUTE addIfNotExists;
DEALLOCATE PREPARE addIfNotExists;

-- For detail_itemsmasuk
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
  "SELECT 1",
  CONCAT("ALTER TABLE ", @tablename, " ADD CONSTRAINT `", @constraintname, "` FOREIGN KEY (`item_id`) REFERENCES `items` (`id`) ON DELETE SET NULL ON UPDATE CASCADE")
));
PREPARE addIfNotExists FROM @preparedStatement;
EXECUTE addIfNotExists;
DEALLOCATE PREPARE addIfNotExists;

