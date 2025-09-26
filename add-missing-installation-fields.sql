-- Script untuk menambahkan kolom yang belum ada di database
-- Tgl. Batas Percobaan dan Tgl. Siap Layanan

-- Step 1: Tambahkan kolom Tgl. Batas Percobaan
ALTER TABLE customer_installations 
ADD COLUMN `trial_end_date` DATE DEFAULT NULL 
COMMENT 'Tanggal batas percobaan layanan';

-- Step 2: Tambahkan kolom Tgl. Siap Layanan
ALTER TABLE customer_installations 
ADD COLUMN `service_ready_date` DATE DEFAULT NULL 
COMMENT 'Tanggal siap layanan';

-- Step 3: Tambahkan index untuk kolom baru
ALTER TABLE customer_installations 
ADD INDEX `idx_customer_installations_trial_end_date` (`trial_end_date`),
ADD INDEX `idx_customer_installations_service_ready_date` (`service_ready_date`);

-- Step 4: Update sample data untuk testing
UPDATE customer_installations 
SET 
    trial_end_date = DATE_ADD(on_air_date, INTERVAL 30 DAY),
    service_ready_date = on_air_date
WHERE id = '949189b8-97d1-11f0-9fa2-d843ae0f1e06';

-- Step 5: Verifikasi struktur tabel setelah perubahan
SELECT 
    'customer_installations' as table_name,
    COUNT(*) as column_count,
    'Missing fields added successfully' as status
FROM information_schema.COLUMNS 
WHERE TABLE_NAME = 'customer_installations' 
AND TABLE_SCHEMA = 'iqgncnzy_skripsi';

-- Step 6: Tampilkan struktur tabel yang baru
DESCRIBE customer_installations;

-- Step 7: Tampilkan sample data dengan kolom baru
SELECT 
    id,
    customer_id,
    technician_id,
    status,
    on_air_date,
    trial_end_date,
    service_ready_date,
    installation_type,
    total_assets_out,
    total_assets_in,
    installation_completed_at
FROM customer_installations 
WHERE id = '949189b8-97d1-11f0-9fa2-d843ae0f1e06';
