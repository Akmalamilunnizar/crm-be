-- Script untuk menambahkan kolom dokumen ke tabel customer_installations
-- 1. Kolom document_type untuk dropdown pilihan dokumen (KTP, SIM, Paspor)
-- 2. Kolom document_photo untuk menyimpan foto dokumen

-- Step 1: Tambahkan kolom document_type dengan ENUM
ALTER TABLE customer_installations 
ADD COLUMN document_type ENUM('KTP', 'SIM', 'Paspor') DEFAULT NULL 
COMMENT 'Jenis dokumen identitas yang digunakan';

-- Step 2: Tambahkan kolom document_photo untuk menyimpan path foto dokumen
ALTER TABLE customer_installations 
ADD COLUMN document_photo VARCHAR(255) DEFAULT NULL 
COMMENT 'Path file foto dokumen identitas';

-- Step 3: Tambahkan index untuk document_type untuk performa query
ALTER TABLE customer_installations 
ADD INDEX idx_customer_installations_document_type (document_type);

-- Step 4: Verifikasi struktur tabel setelah perubahan
SELECT 
    'customer_installations' as table_name,
    COUNT(*) as column_count,
    'Document fields added successfully' as status
FROM information_schema.COLUMNS 
WHERE TABLE_NAME = 'customer_installations' 
AND TABLE_SCHEMA = 'iqgncnzy_skripsi';

-- Step 5: Tampilkan struktur tabel yang baru
DESCRIBE customer_installations;

-- Step 6: Tampilkan sample data dengan field baru
SELECT 
    id,
    customer_id,
    technician_id,
    status,
    document_type,
    document_photo,
    notes,
    on_air_date,
    createdAt
FROM customer_installations 
LIMIT 5;
