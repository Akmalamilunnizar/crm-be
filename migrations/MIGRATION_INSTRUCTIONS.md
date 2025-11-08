# Migration Instructions for Items Inventory System

## Overview
This migration renames `IdItems` column to `item_id` in `detail_itemskeluar` and `detail_itemsmasuk` tables, and updates foreign keys to reference the `items` catalog table instead of `assets` table.

## Migration Order

Run these migrations in order:

### 1. Create Items Catalog Table
```sql
-- Run: create_items_catalog_table.sql
```
This creates the `items` table (catalog) and adds sample items.

### 2. Rename IdItems to item_id
```sql
-- Run: rename_iditems_to_item_id.sql
```
This:
- Drops old foreign keys to `assets` table
- Renames `IdItems` → `item_id` in both detail tables
- Adds new foreign keys to `items` table

### 3. Add Unit and Notes Columns
```sql
-- Run: add_unit_notes_to_items_tables.sql
```
This adds `unit` and `notes` columns to detail tables.

### 4. Fix Foreign Keys (if needed)
```sql
-- Run: fix_items_foreign_keys.sql
```
This ensures foreign key types match exactly.

## What Changed

### Before:
```
detail_itemskeluar / detail_itemsmasuk
  - IdItems (references assets.id)
```

### After:
```
detail_itemskeluar / detail_itemsmasuk
  - item_id (references items.id) ← renamed from IdItems
```

## Important Notes

1. **Router/Asset Items**: For routers and other assets, you need to:
   - Create an entry in `items` catalog
   - Link it to the asset via `items.asset_id`
   - Use that item_id in transactions

2. **Existing Data**: If you have existing data in `IdItems`:
   - You'll need to migrate it to the `items` catalog first
   - Or create corresponding items in the catalog

3. **Backward Compatibility**: The old `IdItems` column is completely replaced with `item_id`. Make sure all existing data is migrated.

