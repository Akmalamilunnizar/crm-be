# Inventory System Design Explanation

## Overview
This inventory system tracks goods going in and out using:
- **Items Catalog** (`items` table) - Master data for all inventory items
- **Items Transactions** (`itemskeluar` / `itemsmasuk`) - Parent transaction records
- **Transaction Details** (`detail_itemskeluar` / `detail_itemsmasuk`) - Line items for each transaction

## Database Structure

### 1. Items Catalog (`items` table)
**Purpose**: Master catalog of all inventory items
- Stores item names (Drop 1C, Patchcore, Cvt single A, Kabel lan, etc.)
- Default unit for each item (M, PCS, KG, etc.)
- Optional link to `assets` table for items that are also tracked as assets (like routers)

**Example Data**:
```
id: "item-001"
name: "Drop 1C"
default_unit: "M"
category: "Cable"
```

### 2. Items Transactions (`itemskeluar` / `itemsmasuk`)
**Purpose**: Parent transaction records that group multiple items
- Each transaction has a date, notes, and created_by
- One transaction can contain multiple items (multiple detail records)

**Example**:
```
id: "IK0001"
date: "2024-01-15"
notes: "Installation for customer ABC"
created_by: "user-123"
```

### 3. Transaction Details (`detail_itemskeluar` / `detail_itemsmasuk`)
**Purpose**: Line items in each transaction

**Columns Explained**:
- **`item_id`**: References `items` catalog (NEW - links to item catalog)
- **`IdItems`**: References `assets` table (LEGACY - for backward compatibility with existing system)
- **`QtyKeluar` / `QtyMasuk`**: Quantity going out/in
- **`unit`**: Unit used for THIS transaction (e.g., "M", "PCS", "KG")
  - **Why separate from items.default_unit?** 
    - Items might be used in different units in different transactions
    - Example: "Kabel lan" default is "M", but one transaction might use "50 M" and another "100 M"
    - Records the actual unit used in that specific transaction
- **`notes`**: Additional notes for THIS specific line item
  - **Purpose**: 
    - Record special conditions (e.g., "Used for installation at customer XYZ")
    - Note item condition (e.g., "Damaged, needs replacement")
    - Track item source/destination details
    - Document any special handling instructions

## Example Transaction Flow

### Scenario: Recording "ALAT KELUAR" (Items Going Out)

**Transaction 1: Installation for Customer ABC**
```
itemskeluar:
  id: "IK0001"
  date: "2024-01-15"
  notes: "Installation for customer ABC"

detail_itemskeluar (multiple items in one transaction):
  Item 1:
    item_id: "item-001" (Drop 1C)
    QtyKeluar: 50
    unit: "M" (matches default unit)
    notes: "Used for fiber connection"
    
  Item 2:
    item_id: "item-002" (Patchcore)
    QtyKeluar: 1
    unit: "PCS" (matches default unit)
    notes: NULL
    
  Item 3:
    item_id: "item-003" (Kabel lan)
    QtyKeluar: 20
    unit: "M" (matches default unit)
    notes: "Cat6 cable for ethernet connection"
```

**Why `unit` column is needed**:
- Item "Kabel lan" has default_unit = "M"
- But in transaction, you specify "20 M" - the `unit` column records that "M" was used
- If someone later uses "Kabel lan" in a different unit (rare but possible), it's recorded

**Why `notes` column is useful**:
- Track WHERE item was used: "Used for installation at customer ABC"
- Note CONDITION: "Damaged during installation, needs replacement"
- Document SPECIAL REQUIREMENTS: "Requires special connector"
- Record DETAILS: "Installed in office building, 3rd floor"

## Benefits of This Design

1. **Flexibility**: Items can be used with different units in different transactions
2. **Traceability**: Notes allow tracking where/how items were used
3. **Catalog Management**: Items catalog centralizes all item information
4. **Backward Compatibility**: `IdItems` (asset_id) still works for existing system
5. **Reporting**: Can generate reports by item, by unit, by notes, etc.

## Usage Example

**User wants to record:**
```
ALAT KELUAR:
1. Drop 1C 50 M
2. Patchcore 1 PCS
3. Cvt single A 1 pcs
5. Kabel lan 1M
6. router TOTOLINK
```

**In the system:**
1. Create `itemskeluar` transaction (IK0001)
2. Add 5 detail records:
   - Drop 1C: 50 M (from items catalog)
   - Patchcore: 1 PCS (from items catalog)
   - Cvt single A: 1 PCS (from items catalog)
   - Kabel lan: 1 M (from items catalog)
   - Router TOTOLINK: 1 PCS (from assets table, linked via IdItems)

**Notes can be added:**
- "Drop 1C" notes: "Used for fiber connection to building"
- "Router TOTOLINK" notes: "Serial: RTR-12345, assigned to customer ABC"

