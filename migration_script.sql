-- Migration Script: Update iqgncnzy_skripsi database to match lilly.sql structure
-- Created: $(date)
-- Purpose: Add missing tables and structures from lilly.sql to local database

USE iqgncnzy_skripsi;

-- ==============================================
-- CREATE NEW TABLES
-- ==============================================

-- 1. Create netwatch_configs table
CREATE TABLE IF NOT EXISTS `netwatch_configs` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `host` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `port` int DEFAULT '8728',
  `username` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `use_ssl` tinyint(1) DEFAULT '0',
  `active` tinyint(1) DEFAULT '1',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 2. Create netwatch_devices table
CREATE TABLE IF NOT EXISTS `netwatch_devices` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ip_address` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL,
  `customer_id` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` enum('up','down') COLLATE utf8mb4_unicode_ci DEFAULT 'up',
  `last_seen` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ip_address` (`ip_address`),
  KEY `customer_id` (`customer_id`),
  CONSTRAINT `netwatch_devices_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customer` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3. Create netwatch_events table
CREATE TABLE IF NOT EXISTS `netwatch_events` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `device_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `event_type` enum('up','down') COLLATE utf8mb4_unicode_ci NOT NULL,
  `event_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `raw_data` text COLLATE utf8mb4_unicode_ci,
  `processed` tinyint(1) DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `device_id` (`device_id`),
  CONSTRAINT `netwatch_events_ibfk_1` FOREIGN KEY (`device_id`) REFERENCES `netwatch_devices` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 4. Create recurring_invoices table
CREATE TABLE IF NOT EXISTS `recurring_invoices` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `customer_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `amount` int NOT NULL,
  `invoice_date` date NOT NULL,
  `due_date` date NOT NULL,
  `next_invoice_date` date NOT NULL,
  `frequency` enum('monthly','quarterly','yearly') CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL DEFAULT 'monthly',
  `status` enum('active','stopped','completed') CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL DEFAULT 'active',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `invoice_items` json NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL,
  `created_by` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `recurring_invoices_customer_id_fkey` (`customer_id`),
  KEY `recurring_invoices_created_by_fkey` (`created_by`),
  KEY `idx_recurring_invoices_status` (`status`),
  KEY `idx_recurring_invoices_next_invoice_date` (`next_invoice_date`),
  KEY `idx_recurring_invoices_customer_status` (`customer_id`,`status`),
  KEY `idx_recurring_invoices_frequency` (`frequency`),
  CONSTRAINT `recurring_invoices_created_by_fkey` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE RESTRICT,
  CONSTRAINT `recurring_invoices_customer_id_fkey` FOREIGN KEY (`customer_id`) REFERENCES `customer` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 5. Create recurring_invoice_history table
CREATE TABLE IF NOT EXISTS `recurring_invoice_history` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `recurring_invoice_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `generated_invoice_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `generated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `invoice_date` date NOT NULL,
  `due_date` date NOT NULL,
  PRIMARY KEY (`id`),
  KEY `recurring_invoice_history_recurring_invoice_id_fkey` (`recurring_invoice_id`),
  KEY `recurring_invoice_history_generated_invoice_id_fkey` (`generated_invoice_id`),
  KEY `idx_recurring_invoice_history_generated_at` (`generated_at`),
  CONSTRAINT `recurring_invoice_history_generated_invoice_id_fkey` FOREIGN KEY (`generated_invoice_id`) REFERENCES `invoices` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `recurring_invoice_history_recurring_invoice_id_fkey` FOREIGN KEY (`recurring_invoice_id`) REFERENCES `recurring_invoices` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 6. Create spare_parts table
CREATE TABLE IF NOT EXISTS `spare_parts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text COLLATE utf8mb4_unicode_ci,
  `category` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `is_active` tinyint(1) DEFAULT '1',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_spare_parts_category` (`category`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 7. Create technician_steps table
CREATE TABLE IF NOT EXISTS `technician_steps` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `step_order` int NOT NULL,
  `title` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text COLLATE utf8mb4_unicode_ci,
  `tools` text COLLATE utf8mb4_unicode_ci,
  `spare_parts` text COLLATE utf8mb4_unicode_ci,
  `procedure` text COLLATE utf8mb4_unicode_ci,
  `solution` text COLLATE utf8mb4_unicode_ci,
  `is_active` tinyint(1) DEFAULT '1',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_technician_steps_order` (`step_order`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 8. Create trouble_tickets_original table (must be created before technician_team_members)
CREATE TABLE IF NOT EXISTS `trouble_tickets_original` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `customer_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `title` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `status` enum('finished','ongoing','unfinished') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `assigned_to` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `current_assignee_role` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '1631f7d5-7d01-40af-8d24-1692cefa205a',
  `customer_note` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `noc_note` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `technician_note` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `img_cs` varchar(60) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'Image uploaded by Customer Service',
  `img_noc` varchar(60) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'Image uploaded by NOC',
  `img_tech_bf` varchar(60) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'Image uploaded by Technician (Before)',
  `img_tech_af` varchar(60) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'Image uploaded by Technician (After)',
  `verified_by_cs` tinyint(1) DEFAULT '0',
  `verified_at` datetime DEFAULT NULL,
  `technician_completed` tinyint(1) DEFAULT '0',
  `network_architecture` varchar(50) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `accumulation` int DEFAULT '1' COMMENT 'Number of customers affected by the same problem',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `Index 3` (`assigned_to`,`current_assignee_role`),
  KEY `Index 4` (`customer_id`),
  KEY `Index 5` (`status`),
  KEY `Index 6` (`created_at`),
  KEY `Index 7` (`type`),
  CONSTRAINT `trouble_tickets_original_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customer` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `trouble_tickets_original_ibfk_2` FOREIGN KEY (`assigned_to`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 9. Create technician_team_members table
CREATE TABLE IF NOT EXISTS `technician_team_members` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `ticket_id` bigint unsigned NOT NULL,
  `user_id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `role` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_team_ticket` (`ticket_id`),
  KEY `fk_team_user` (`user_id`),
  CONSTRAINT `fk_team_ticket` FOREIGN KEY (`ticket_id`) REFERENCES `trouble_tickets_original` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_team_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 10. Create ticket_logs table
CREATE TABLE IF NOT EXISTS `ticket_logs` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ticket_id` bigint unsigned NOT NULL,
  `log_type` enum('manual','netwatch','system') COLLATE utf8mb4_unicode_ci DEFAULT 'manual',
  `message` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `event_id` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_by` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `event_id` (`event_id`),
  KEY `created_by` (`created_by`),
  KEY `idx_ticket_logs_ticket` (`ticket_id`),
  KEY `idx_ticket_logs_type` (`log_type`),
  CONSTRAINT `ticket_logs_ibfk_1` FOREIGN KEY (`ticket_id`) REFERENCES `trouble_tickets_original` (`id`) ON DELETE CASCADE,
  CONSTRAINT `ticket_logs_ibfk_2` FOREIGN KEY (`event_id`) REFERENCES `netwatch_events` (`id`) ON DELETE SET NULL,
  CONSTRAINT `ticket_logs_ibfk_3` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 11. Create ticket_technician_steps table
CREATE TABLE IF NOT EXISTS `ticket_technician_steps` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `ticket_id` bigint unsigned NOT NULL,
  `step_id` bigint unsigned NOT NULL,
  `technician_id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT 'pending',
  `notes` text COLLATE utf8mb4_unicode_ci,
  `spare_parts_used` text COLLATE utf8mb4_unicode_ci,
  `image_paths` text COLLATE utf8mb4_unicode_ci,
  `completed_at` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_ticket_step_technician` (`ticket_id`,`step_id`,`technician_id`),
  KEY `idx_ticket_tech_steps_ticket` (`ticket_id`),
  KEY `idx_ticket_tech_steps_step` (`step_id`),
  KEY `idx_ticket_tech_steps_technician` (`technician_id`),
  CONSTRAINT `fk_ticket_tech_steps_step` FOREIGN KEY (`step_id`) REFERENCES `technician_steps` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_ticket_tech_steps_technician` FOREIGN KEY (`technician_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_ticket_tech_steps_ticket` FOREIGN KEY (`ticket_id`) REFERENCES `trouble_tickets_original` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==============================================
-- ALTER EXISTING TABLES
-- ==============================================

-- Add UNIQUE constraint to network_devices.mac_address if it doesn't exist
-- Note: This will only work if there are no duplicate mac_address values
-- ALTER TABLE `network_devices` ADD UNIQUE KEY `Index 5` (`mac_address`);

-- ==============================================
-- INSERT SAMPLE DATA (Optional)
-- ==============================================

-- Insert sample spare parts data
INSERT IGNORE INTO `spare_parts` (`id`, `name`, `description`, `category`, `is_active`, `created_at`, `updated_at`) VALUES
(1, 'Stop Kontak', 'Power outlet untuk peralatan', 'Electrical', 1, '2025-09-22 08:44:12', '2025-09-22 08:44:12'),
(2, 'Kabel Listrik', 'Kabel power untuk peralatan', 'Electrical', 1, '2025-09-22 08:44:12', '2025-09-22 08:44:12'),
(3, 'CVT', 'Media Converter', 'Network Device', 1, '2025-09-22 08:44:12', '2025-09-22 08:44:12'),
(4, 'Fastcon', 'Konektor fiber optik', 'Fiber Optic', 1, '2025-09-22 08:44:12', '2025-09-22 08:44:12'),
(5, 'PatchCord', 'Kabel patch untuk koneksi', 'Fiber Optic', 1, '2025-09-22 08:44:12', '2025-09-22 08:44:12'),
(6, 'Splicer', 'Alat untuk menyambung kabel optik', 'Fiber Optic', 1, '2025-09-22 08:44:12', '2025-09-22 08:44:12'),
(7, 'UTP', 'Kabel UTP untuk jaringan', 'Network Cable', 1, '2025-09-22 08:44:13', '2025-09-22 08:44:13'),
(8, 'RJ45', 'Konektor kabel UTP', 'Network Connector', 1, '2025-09-22 08:44:13', '2025-09-22 08:44:13'),
(9, 'Adaptor', 'Power adapter untuk peralatan', 'Electrical', 1, '2025-09-22 08:44:13', '2025-09-22 08:44:13'),
(10, 'Router', 'Perangkat Router Wireless', 'Network Device', 1, '2025-09-22 08:44:13', '2025-09-22 08:44:13');

-- Insert sample technician steps data
INSERT IGNORE INTO `technician_steps` (`id`, `step_order`, `title`, `description`, `tools`, `spare_parts`, `procedure`, `solution`, `is_active`, `created_at`, `updated_at`) VALUES
(1, 1, 'Pengecekan Power Listrik Utama', 'Memeriksa kondisi power listrik utama untuk memastikan suplai listrik stabil', 'Multimeter, Test Pen, Kabel Tester', 'Stop Kontak, Kabel Listrik, Fuse', '1. Periksa tegangan listrik dengan multimeter\n2. Cek kondisi stop kontak\n3. Test koneksi kabel listrik\n4. Verifikasi fuse tidak putus', 'Ganti komponen yang rusak, perbaiki koneksi yang longgar', 1, '2025-09-22 08:44:12', '2025-09-22 08:44:12'),
(2, 2, 'Pengecekan Adaptor Perangkat Media Converter', 'Memeriksa kondisi adaptor power untuk media converter', 'Multimeter, Test Pen', 'Adaptor, Kabel Power', '1. Periksa output voltage adaptor\n2. Cek kondisi kabel adaptor\n3. Test koneksi ke media converter', 'Ganti adaptor jika output tidak sesuai, perbaiki koneksi kabel', 1, '2025-09-22 08:44:12', '2025-09-22 08:44:12'),
(3, 3, 'Pengecekan Perangkat Media Converter (CVT)', 'Memeriksa kondisi dan fungsi media converter', 'Multimeter, Test Pen, Laptop', 'CVT, Fastcon, PatchCord', '1. Periksa LED indikator CVT\n2. Test koneksi input/output\n3. Cek konfigurasi CVT\n4. Verifikasi link status', 'Reset CVT, ganti jika rusak, perbaiki konfigurasi', 1, '2025-09-22 08:44:12', '2025-09-22 08:44:12'),
(4, 4, 'Pengecekan Status Link Upstream (Kedua Sisi Ujung)', 'Memeriksa status koneksi upstream dari kedua sisi', 'Laptop, Ping Test, Traceroute', 'PatchCord, Splicer', '1. Ping test ke gateway\n2. Traceroute ke server\n3. Cek latency dan packet loss\n4. Verifikasi routing', 'Perbaiki routing, ganti kabel jika rusak, konfigurasi ulang', 1, '2025-09-22 08:44:12', '2025-09-22 08:44:12'),
(5, 5, 'Pencarian & Perbaikan Titik Kesalahan/Putus/Kerusakan Kabel (Optik)', 'Mencari dan memperbaiki kerusakan pada kabel optik', 'OTDR, Power Meter, Splicer', 'Kabel Optik, Splicer, Fastcon', '1. Gunakan OTDR untuk deteksi kerusakan\n2. Cek power loss pada kabel\n3. Identifikasi titik putus\n4. Lakukan splicing jika diperlukan', 'Splice kabel yang putus, ganti kabel jika rusak parah', 1, '2025-09-22 08:44:12', '2025-09-22 08:44:12'),
(6, 6, 'Pengecekan Fungsi Kabel UTP', 'Memeriksa kondisi dan fungsi kabel UTP', 'Cable Tester, Multimeter', 'UTP, RJ45', '1. Test continuity kabel UTP\n2. Cek koneksi RJ45\n3. Verifikasi pinout kabel\n4. Test kecepatan kabel', 'Ganti kabel UTP jika rusak, crimping ulang RJ45', 1, '2025-09-22 08:44:12', '2025-09-22 08:44:12'),
(7, 7, 'Pengecekan Adaptor Perangkat Router Wireless', 'Memeriksa kondisi adaptor power untuk router wireless', 'Multimeter, Test Pen', 'Adaptor, Kabel Power', '1. Periksa output voltage adaptor\n2. Cek kondisi kabel adaptor\n3. Test koneksi ke router', 'Ganti adaptor jika output tidak sesuai', 1, '2025-09-22 08:44:12', '2025-09-22 08:44:12'),
(8, 8, 'Pengecekan Perangkat Router', 'Memeriksa kondisi dan fungsi router', 'Laptop, Multimeter, Test Pen', 'Router, Adaptor', '1. Periksa LED indikator router\n2. Test koneksi LAN/WAN\n3. Cek konfigurasi router\n4. Verifikasi DHCP server', 'Reset router, ganti jika rusak, update firmware', 1, '2025-09-22 08:44:12', '2025-09-22 08:44:12');

-- ==============================================
-- MIGRATION COMPLETED
-- ==============================================

SELECT 'Migration completed successfully!' as status;
SELECT COUNT(*) as total_tables FROM information_schema.tables WHERE table_schema = 'iqgncnzy_skripsi';

