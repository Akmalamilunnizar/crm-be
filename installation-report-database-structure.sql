-- Script untuk menyesuaikan struktur database untuk fitur Add Report Installation
-- Membuat relasi yang tepat antara customer_installations dengan tabel lain

-- Step 1: Tambahkan relasi customer_installations ke network_devices
-- Ini akan memungkinkan satu instalasi memiliki multiple network devices
ALTER TABLE network_devices 
ADD COLUMN `customer_installation_id` varchar(191) DEFAULT NULL 
COMMENT 'Reference to customer_installations table';

-- Step 2: Tambahkan relasi customer_installations ke customer_services
-- Ini akan memungkinkan satu instalasi memiliki multiple customer services
ALTER TABLE customer_services 
ADD COLUMN `customer_installation_id` varchar(191) DEFAULT NULL 
COMMENT 'Reference to customer_installations table';

-- Step 3: Tambahkan relasi customer_installations ke cable
-- Ini akan memungkinkan tracking kabel yang digunakan dalam instalasi
ALTER TABLE cable 
ADD COLUMN `customer_installation_id` varchar(191) DEFAULT NULL 
COMMENT 'Reference to customer_installations table';

-- Step 4: Tambahkan index untuk relasi baru
ALTER TABLE network_devices 
ADD INDEX `idx_network_devices_customer_installation_id` (`customer_installation_id`);

ALTER TABLE customer_services 
ADD INDEX `idx_customer_services_customer_installation_id` (`customer_installation_id`);

ALTER TABLE cable 
ADD INDEX `idx_cable_customer_installation_id` (`customer_installation_id`);

-- Step 5: Tambahkan foreign key constraints
ALTER TABLE network_devices 
ADD CONSTRAINT `fk_network_devices_customer_installation` 
FOREIGN KEY (`customer_installation_id`) REFERENCES `customer_installations` (`id`) 
ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE customer_services 
ADD CONSTRAINT `fk_customer_services_customer_installation` 
FOREIGN KEY (`customer_installation_id`) REFERENCES `customer_installations` (`id`) 
ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE cable 
ADD CONSTRAINT `fk_cable_customer_installation` 
FOREIGN KEY (`customer_installation_id`) REFERENCES `customer_installations` (`id`) 
ON DELETE SET NULL ON UPDATE CASCADE;

-- Step 6: Update data existing untuk testing
-- Update network_devices yang sudah ada untuk relasi dengan customer_installations
UPDATE network_devices 
SET customer_installation_id = '949189b8-97d1-11f0-9fa2-d843ae0f1e06'
WHERE customer_id = '24225552-b0c3-43d7-9a67-20e04f36fa5f';

UPDATE network_devices 
SET customer_installation_id = '94923543-97d1-11f0-9fa2-d843ae0f1e06'
WHERE customer_id = '55f03555-973d-425b-92bd-a0825c20b8ce';

-- Step 7: Insert sample data untuk customer_services
INSERT INTO customer_services (
    id, 
    customer_id, 
    customer_installation_id,
    device_id, 
    cable_id, 
    cable_length, 
    end_port_type, 
    user_login, 
    password, 
    user_status, 
    installation_notes, 
    installation_team_phone
) VALUES 
(
    UUID(), 
    '24225552-b0c3-43d7-9a67-20e04f36fa5f', 
    '949189b8-97d1-11f0-9fa2-d843ae0f1e06',
    'sadvwabafbsdgd', 
    'cable-1', 
    100.00, 
    'RJ45', 
    'cha_eunwoo', 
    'password123', 
    'Active', 
    'Installation completed successfully', 
    '081234567890'
),
(
    UUID(), 
    '55f03555-973d-425b-92bd-a0825c20b8ce', 
    '94923543-97d1-11f0-9fa2-d843ae0f1e06',
    'device-55f03555', 
    'cable-2', 
    50.00, 
    'Fiber', 
    'rita_darmayanti', 
    'password456', 
    'Active', 
    'Installation in progress', 
    '081234567891'
);

-- Step 8: Insert sample data untuk cable yang digunakan dalam instalasi
INSERT INTO cable (
    id, 
    name, 
    type, 
    length, 
    status, 
    customer_installation_id
) VALUES 
(
    'cable-installation-1', 
    'UTP Cable for Cha Eunwoo Installation', 
    'UTP Cat6', 
    100.00, 
    'in_use', 
    '949189b8-97d1-11f0-9fa2-d843ae0f1e06'
),
(
    'cable-installation-2', 
    'Fiber Cable for Rita Installation', 
    'Single Mode Fiber', 
    50.00, 
    'in_use', 
    '94923543-97d1-11f0-9fa2-d843ae0f1e06'
);

-- Step 9: Verifikasi struktur tabel setelah perubahan
SELECT 
    'network_devices' as table_name,
    COUNT(*) as column_count,
    'Customer installation relation added' as status
FROM information_schema.COLUMNS 
WHERE TABLE_NAME = 'network_devices' 
AND TABLE_SCHEMA = 'iqgncnzy_skripsi';

SELECT 
    'customer_services' as table_name,
    COUNT(*) as column_count,
    'Customer installation relation added' as status
FROM information_schema.COLUMNS 
WHERE TABLE_NAME = 'customer_services' 
AND TABLE_SCHEMA = 'iqgncnzy_skripsi';

SELECT 
    'cable' as table_name,
    COUNT(*) as column_count,
    'Customer installation relation added' as status
FROM information_schema.COLUMNS 
WHERE TABLE_NAME = 'cable' 
AND TABLE_SCHEMA = 'iqgncnzy_skripsi';

-- Step 10: Tampilkan sample data dengan relasi baru
SELECT 
    'Sample Installation Report Data' as info,
    ci.id as installation_id,
    c.name as customer_name,
    u.name as technician_name,
    ci.status as installation_status,
    ci.on_air_date,
    ci.trial_end_date,
    ci.service_ready_date
FROM customer_installations ci
JOIN customer c ON ci.customer_id = c.id
JOIN users u ON ci.technician_id = u.id
LIMIT 2;
