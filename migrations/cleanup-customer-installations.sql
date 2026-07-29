-- Script untuk membersihkan tabel customer_installations
-- 1. Hapus tabel backup yang tidak digunakan
-- 2. Hapus kolom yang tidak diperlukan

-- Step 1: Hapus tabel backup
DROP TABLE IF EXISTS customer_installations_backup;

-- Step 2: Hapus kolom yang tidak diperlukan dari customer_installations
-- Kolom yang akan dihapus: description, date_trial, equipment_used, completion_time, customer_signature

-- Hapus kolom satu per satu
ALTER TABLE customer_installations DROP COLUMN description;
ALTER TABLE customer_installations DROP COLUMN date_trial;
ALTER TABLE customer_installations DROP COLUMN equipment_used;
ALTER TABLE customer_installations DROP COLUMN completion_time;
ALTER TABLE customer_installations DROP COLUMN customer_signature;

-- Step 3: Verifikasi struktur tabel setelah perubahan
SELECT 
    'customer_installations' as table_name,
    COUNT(*) as column_count,
    'Columns removed successfully' as status
FROM information_schema.COLUMNS 
WHERE TABLE_NAME = 'customer_installations' 
AND TABLE_SCHEMA = 'iqgncnzy_skripsi';

-- Step 4: Tampilkan struktur tabel yang baru
DESCRIBE customer_installations;
