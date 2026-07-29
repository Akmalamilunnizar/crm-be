-- Script untuk sistem tracking aset keluar/masuk pada customer installations
-- Fitur Add Report Installation dengan aset keluar dan masuk

-- Step 1: Buat tabel untuk tracking aset keluar/masuk
CREATE TABLE IF NOT EXISTS `asset_transactions` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `customer_installation_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `asset_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `transaction_type` enum('out','in') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'out = aset keluar, in = aset masuk',
  `quantity` int NOT NULL DEFAULT 1 COMMENT 'Jumlah aset yang keluar/masuk',
  `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT 'Catatan transaksi aset',
  `transaction_date` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT 'Tanggal transaksi',
  `created_by` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'User yang melakukan transaksi',
  `createdAt` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updatedAt` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_asset_transactions_customer_installation_id` (`customer_installation_id`),
  KEY `idx_asset_transactions_asset_id` (`asset_id`),
  KEY `idx_asset_transactions_transaction_type` (`transaction_type`),
  KEY `idx_asset_transactions_transaction_date` (`transaction_date`),
  KEY `idx_asset_transactions_created_by` (`created_by`),
  CONSTRAINT `fk_asset_transactions_customer_installation` FOREIGN KEY (`customer_installation_id`) REFERENCES `customer_installations` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_asset_transactions_asset` FOREIGN KEY (`asset_id`) REFERENCES `assets` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_asset_transactions_user` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Step 2: Tambahkan kolom untuk tracking aset di customer_installations
ALTER TABLE customer_installations 
ADD COLUMN `installation_type` enum('new_installation','maintenance','upgrade','downgrade') DEFAULT 'new_installation' 
COMMENT 'Jenis instalasi yang dilakukan',
ADD COLUMN `total_assets_out` int DEFAULT 0 
COMMENT 'Total aset yang keluar untuk instalasi ini',
ADD COLUMN `total_assets_in` int DEFAULT 0 
COMMENT 'Total aset yang masuk dari instalasi ini',
ADD COLUMN `installation_completed_at` datetime(3) DEFAULT NULL 
COMMENT 'Tanggal instalasi selesai';

-- Step 3: Tambahkan index untuk kolom baru
ALTER TABLE customer_installations 
ADD INDEX `idx_customer_installations_installation_type` (`installation_type`),
ADD INDEX `idx_customer_installations_installation_completed_at` (`installation_completed_at`);

-- Step 4: Buat view untuk laporan aset keluar/masuk
CREATE OR REPLACE VIEW `asset_transaction_report` AS
SELECT 
    at.id as transaction_id,
    at.customer_installation_id,
    ci.customer_id,
    c.name as customer_name,
    at.asset_id,
    a.brand,
    a.type as asset_type,
    a.model,
    a.serial_number,
    at.transaction_type,
    at.quantity,
    at.notes as transaction_notes,
    at.transaction_date,
    at.created_by,
    u.name as created_by_name,
    ci.installation_type,
    ci.status as installation_status,
    ci.on_air_date
FROM asset_transactions at
JOIN customer_installations ci ON at.customer_installation_id = ci.id
JOIN customer c ON ci.customer_id = c.id
JOIN assets a ON at.asset_id = a.id
JOIN users u ON at.created_by = u.id
ORDER BY at.transaction_date DESC;

-- Step 5: Buat view untuk summary aset per instalasi
CREATE OR REPLACE VIEW `installation_asset_summary` AS
SELECT 
    ci.id as installation_id,
    ci.customer_id,
    c.name as customer_name,
    ci.installation_type,
    ci.status,
    ci.on_air_date,
    ci.installation_completed_at,
    COUNT(CASE WHEN at.transaction_type = 'out' THEN 1 END) as total_assets_out,
    COUNT(CASE WHEN at.transaction_type = 'in' THEN 1 END) as total_assets_in,
    SUM(CASE WHEN at.transaction_type = 'out' THEN at.quantity ELSE 0 END) as total_quantity_out,
    SUM(CASE WHEN at.transaction_type = 'in' THEN at.quantity ELSE 0 END) as total_quantity_in
FROM customer_installations ci
LEFT JOIN asset_transactions at ON ci.id = at.customer_installation_id
LEFT JOIN customer c ON ci.customer_id = c.id
GROUP BY ci.id, ci.customer_id, c.name, ci.installation_type, ci.status, ci.on_air_date, ci.installation_completed_at;

-- Step 6: Insert sample data untuk testing
INSERT INTO `asset_transactions` (
    `id`, 
    `customer_installation_id`, 
    `asset_id`, 
    `transaction_type`, 
    `quantity`, 
    `notes`, 
    `transaction_date`, 
    `created_by`, 
    `createdAt`, 
    `updatedAt`
) VALUES 
(
    UUID(), 
    '949189b8-97d1-11f0-9fa2-d843ae0f1e06', 
    '5ca1606b-66b3-4958-b7af-f48d4cda800a', 
    'out', 
    1, 
    'Router keluar untuk instalasi customer baru', 
    NOW(), 
    'c13a6c87-ec28-47ba-84c2-58b5ace2af57', 
    NOW(), 
    NOW()
),
(
    UUID(), 
    '949189b8-97d1-11f0-9fa2-d843ae0f1e06', 
    '842a48a4-6380-4280-96e1-a734f61a7d5b', 
    'out', 
    1, 
    'Modem keluar untuk instalasi customer baru', 
    NOW(), 
    'c13a6c87-ec28-47ba-84c2-58b5ace2af57', 
    NOW(), 
    NOW()
);

-- Step 7: Update customer_installations dengan data sample
UPDATE customer_installations 
SET 
    installation_type = 'new_installation',
    total_assets_out = 2,
    total_assets_in = 0,
    installation_completed_at = NOW()
WHERE id = '949189b8-97d1-11f0-9fa2-d843ae0f1e06';

-- Step 8: Verifikasi struktur tabel
SELECT 
    'asset_transactions' as table_name,
    COUNT(*) as column_count,
    'Asset tracking table created successfully' as status
FROM information_schema.COLUMNS 
WHERE TABLE_NAME = 'asset_transactions' 
AND TABLE_SCHEMA = 'iqgncnzy_skripsi';

-- Step 9: Tampilkan struktur tabel yang baru
DESCRIBE asset_transactions;

-- Step 10: Tampilkan sample data
SELECT 
    'Sample Asset Transactions' as info,
    COUNT(*) as total_transactions
FROM asset_transactions;
