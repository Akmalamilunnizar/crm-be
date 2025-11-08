-- Simulation: Example of adding items to inventory
-- This shows how the system works with items catalog

-- ============================================
-- STEP 1: Items are added to catalog (master data)
-- ============================================
-- This is typically done once when setting up the system or when new item types are introduced

INSERT INTO `items` (`id`, `name`, `default_unit`, `category`, `description`) VALUES
  ('item-drop-1c', 'Drop 1C', 'M', 'Cable', 'Drop cable 1 core'),
  ('item-patchcore', 'Patchcore', 'PCS', 'Connector', 'Patchcore connector'),
  ('item-cvt-single-a', 'Cvt single A', 'PCS', 'Connector', 'Cvt single type A connector'),
  ('item-kabel-lan', 'Kabel lan', 'M', 'Cable', 'LAN cable (UTP/STP)')
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- ============================================
-- STEP 2: Items come IN (purchase/stock in)
-- ============================================
-- Example: Purchasing items from supplier

-- Create parent transaction for items coming in
INSERT INTO `itemsmasuk` (`id`, `date`, `notes`, `created_by`) VALUES
  ('IM0001', '2024-01-10', 'Purchase from supplier XYZ', 'user-123');

-- Add detail items that came in
INSERT INTO `detail_itemsmasuk` (`IdMasuk`, `item_id`, `QtyMasuk`, `unit`, `HargaSatuan`, `SubTotal`, `notes`) VALUES
  ('IM0001', 'item-drop-1c', 500, 'M', 5000, 2500000, 'Received 5 rolls of 100M each'),
  ('IM0001', 'item-patchcore', 50, 'PCS', 15000, 750000, 'Standard patchcore connectors'),
  ('IM0001', 'item-cvt-single-a', 30, 'PCS', 20000, 600000, 'Type A connectors'),
  ('IM0001', 'item-kabel-lan', 1000, 'M', 3000, 3000000, 'Cat6 cables');

-- Result: Inventory now has:
-- - Drop 1C: 500 M
-- - Patchcore: 50 PCS
-- - Cvt single A: 30 PCS
-- - Kabel lan: 1000 M

-- ============================================
-- STEP 3: Items go OUT (usage/installation)
-- ============================================
-- Example: Recording "ALAT KELUAR" for installation

-- Create parent transaction for items going out
INSERT INTO `itemskeluar` (`id`, `date`, `notes`, `created_by`) VALUES
  ('IK0001', '2024-01-15', 'Installation for customer ABC - Building installation', 'user-123');

-- Add detail items that went out
INSERT INTO `detail_itemskeluar` (`IdKeluar`, `item_id`, `QtyKeluar`, `unit`, `notes`) VALUES
  ('IK0001', 'item-drop-1c', 50, 'M', 'Used for fiber connection to building entrance'),
  ('IK0001', 'item-patchcore', 1, 'PCS', 'Installed at termination point'),
  ('IK0001', 'item-cvt-single-a', 1, 'PCS', 'Connector for main fiber line'),
  ('IK0001', 'item-kabel-lan', 20, 'M', 'Cat6 cable for office ethernet connection');

-- Result: Inventory reduced by:
-- - Drop 1C: -50 M (now 450 M remaining)
-- - Patchcore: -1 PCS (now 49 PCS remaining)
-- - Cvt single A: -1 PCS (now 29 PCS remaining)
-- - Kabel lan: -20 M (now 980 M remaining)

-- ============================================
-- STEP 4: Router from assets (special case)
-- ============================================
-- Routers are tracked in assets table, not items catalog
-- They can be included in transactions using IdItems (asset_id)

-- Example: Router going out
INSERT INTO `itemskeluar` (`id`, `date`, `notes`, `created_by`) VALUES
  ('IK0002', '2024-01-15', 'Router deployment for customer ABC', 'user-123');

-- Router from assets table (not from items catalog)
-- Assuming router asset_id = 'asset-router-totolink-001'
INSERT INTO `detail_itemskeluar` (`IdKeluar`, `IdItems`, `QtyKeluar`, `unit`, `notes`) VALUES
  ('IK0002', 'asset-router-totolink-001', 1, 'PCS', 'TOTOLINK router, Serial: RTR-12345, assigned to customer ABC');

-- ============================================
-- EXPLANATION OF COLUMNS:
-- ============================================

-- **`unit` column in detail tables:**
-- - Records the unit used in THIS specific transaction
-- - Example: "Drop 1C" has default_unit = "M" in items catalog
--   But in transaction IK0001, we used "50 M" - the unit column records "M"
-- - Why separate? Items might be used in different units:
--   - "Kabel lan" might be used as "100 M" in one transaction, "50 M" in another
--   - The unit column records the actual unit used in each transaction

-- **`notes` column in detail tables:**
-- - Purpose 1: Track WHERE item was used
--   Example: "Used for installation at customer ABC building, 3rd floor"
--   
-- - Purpose 2: Note CONDITION of item
--   Example: "Damaged during installation, needs replacement"
--   
-- - Purpose 3: Document SPECIAL REQUIREMENTS
--   Example: "Requires special connector adapter"
--   
-- - Purpose 4: Record DETAILED INFORMATION
--   Example: "Router Serial: RTR-12345, MAC: AA:BB:CC:DD:EE:FF"
--   
-- - Purpose 5: Track SOURCE/DESTINATION
--   Example: "Returned from customer XYZ, needs refurbishment"
--   
-- - Purpose 6: Installation SPECIFICS
--   Example: "Installed at port 10, switch SW-001"

-- ============================================
-- QUERY EXAMPLES:
-- ============================================

-- View all items in catalog
SELECT * FROM `items` ORDER BY `category`, `name`;

-- View inventory transactions with items
SELECT 
  ik.id as transaction_id,
  ik.date,
  ik.notes as transaction_notes,
  i.name as item_name,
  dik.QtyKeluar as quantity,
  dik.unit,
  dik.notes as item_notes
FROM `itemskeluar` ik
JOIN `detail_itemskeluar` dik ON dik.IdKeluar = ik.id
LEFT JOIN `items` i ON i.id = dik.item_id
ORDER BY ik.date DESC;

-- Calculate remaining inventory (simplified - would need proper tracking)
SELECT 
  i.name,
  i.default_unit,
  COALESCE(SUM(CASE WHEN dik.QtyKeluar IS NOT NULL THEN dik.QtyKeluar ELSE 0 END), 0) as total_out,
  COALESCE(SUM(CASE WHEN dim.QtyMasuk IS NOT NULL THEN dim.QtyMasuk ELSE 0 END), 0) as total_in,
  (COALESCE(SUM(CASE WHEN dim.QtyMasuk IS NOT NULL THEN dim.QtyMasuk ELSE 0 END), 0) - 
   COALESCE(SUM(CASE WHEN dik.QtyKeluar IS NOT NULL THEN dik.QtyKeluar ELSE 0 END), 0)) as remaining
FROM `items` i
LEFT JOIN `detail_itemskeluar` dik ON dik.item_id = i.id
LEFT JOIN `detail_itemsmasuk` dim ON dim.item_id = i.id
GROUP BY i.id, i.name, i.default_unit;

