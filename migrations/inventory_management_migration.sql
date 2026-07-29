-- Inventory Management System Database Migration
-- This script creates the necessary tables for the robust inventory and asset management system

USE `iqgncnzy_skripsi`;

-- Create the Product Catalog table (assets) - Updated version
CREATE TABLE IF NOT EXISTS `assets` (
  `id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `brand` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `model` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `serial_number` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'This is a non-unique Product Code/SKU',
  `price` double NOT NULL,
  `quantity` double DEFAULT NULL COMMENT 'Reference for total purchased, not for live inventory.',
  `description` text COLLATE utf8mb4_unicode_ci,
  `status` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT 'active',
  `status_in_out` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT 'in',
  `site` varchar(191) COLLATE utf8mb4_unicode_ci,
  `date` varchar(191) COLLATE utf8mb4_unicode_ci,
  `company_id` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `createdAt` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updatedAt` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `assets_company_fk` (`company_id`),
  CONSTRAINT `assets_company_fk` FOREIGN KEY (`company_id`) REFERENCES `company` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create the Physical Inventory table (asset_items) - One row per physical device
CREATE TABLE IF NOT EXISTS `asset_items` (
  `id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `asset_id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Links to the asset type in the assets table',
  `mac_address` varchar(17) COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` enum('in_stock','in_use','maintenance','damaged','retired') NOT NULL DEFAULT 'in_stock',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `mac_address_unique` (`mac_address`),
  KEY `asset_id_fk` (`asset_id`),
  CONSTRAINT `asset_items_asset_id_fk` FOREIGN KEY (`asset_id`) REFERENCES `assets` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create the main purchase record table (barangmasuk)
CREATE TABLE IF NOT EXISTS `barangmasuk` (
  `IdMasuk` varchar(6) COLLATE utf8mb4_general_ci NOT NULL,
  `date` date NOT NULL,
  `notes` text COLLATE utf8mb4_unicode_ci,
  `created_by` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`IdMasuk`),
  KEY `created_by_fk` (`created_by`),
  CONSTRAINT `barangmasuk_created_by_fk` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Create the purchase detail table (detail_barangmasuk)
CREATE TABLE IF NOT EXISTS `detail_barangmasuk` (
  `IdMasuk` varchar(6) COLLATE utf8mb4_general_ci NOT NULL,
  `asset_id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Links to the assets catalog table',
  `serial_number` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Serial number of the box/batch',
  `QtyMasuk` int DEFAULT NULL,
  `HargaSatuan` int NOT NULL,
  `SubTotal` int NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  KEY `IdMasuk` (`IdMasuk`),
  KEY `asset_id` (`asset_id`),
  CONSTRAINT `detail_barangmasuk_asset_id_fk` FOREIGN KEY (`asset_id`) REFERENCES `assets` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `detail_barangmasuk_id_masuk_fk` FOREIGN KEY (`IdMasuk`) REFERENCES `barangmasuk` (`IdMasuk`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Create the installation asset transactions table (asset_transactions)
CREATE TABLE IF NOT EXISTS `asset_transactions` (
  `id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `customer_installation_id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `asset_item_id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `transaction_type` enum('out','in') COLLATE utf8mb4_unicode_ci NOT NULL,
  `notes` text COLLATE utf8mb4_unicode_ci,
  `created_by` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `customer_installation_id_fk` (`customer_installation_id`),
  KEY `asset_item_id_fk` (`asset_item_id`),
  KEY `created_by_fk` (`created_by`),
  CONSTRAINT `asset_transactions_customer_installation_fk` FOREIGN KEY (`customer_installation_id`) REFERENCES `customer_installations` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `asset_transactions_asset_item_fk` FOREIGN KEY (`asset_item_id`) REFERENCES `asset_items` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `asset_transactions_created_by_fk` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create the trouble ticket asset transactions table (ticket_asset_transactions)
CREATE TABLE IF NOT EXISTS `ticket_asset_transactions` (
  `id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `trouble_ticket_id` bigint unsigned NOT NULL,
  `asset_item_id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `transaction_type` enum('out','in') COLLATE utf8mb4_unicode_ci NOT NULL,
  `notes` text COLLATE utf8mb4_unicode_ci,
  `created_by` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `trouble_ticket_id_fk` (`trouble_ticket_id`),
  KEY `asset_item_id_fk` (`asset_item_id`),
  KEY `created_by_fk` (`created_by`),
  CONSTRAINT `ticket_asset_transactions_trouble_ticket_fk` FOREIGN KEY (`trouble_ticket_id`) REFERENCES `trouble_tickets` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `ticket_asset_transactions_asset_item_fk` FOREIGN KEY (`asset_item_id`) REFERENCES `asset_items` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `ticket_asset_transactions_created_by_fk` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create indexes for better performance
CREATE INDEX `idx_asset_items_status` ON `asset_items` (`status`);
CREATE INDEX `idx_asset_items_created_at` ON `asset_items` (`created_at`);
CREATE INDEX `idx_asset_transactions_type` ON `asset_transactions` (`transaction_type`);
CREATE INDEX `idx_ticket_asset_transactions_type` ON `ticket_asset_transactions` (`transaction_type`);

-- Insert sample data for testing (optional)
-- INSERT INTO `assets` (`id`, `brand`, `type`, `model`, `serial_number`, `price`, `quantity`, `description`, `status`, `status_in_out`, `site`, `date`, `createdAt`, `updatedAt`) VALUES
-- ('asset-001', 'TP-Link', 'Router', 'Archer C7', 'TPL-AC7-001', 250000, 100, 'Wireless AC1750 Dual Band Router', 'active', 'in', 'Warehouse', '2024-01-15', NOW(), NOW()),
-- ('asset-002', 'MikroTik', 'Router', 'RB750Gr3', 'MT-RB750-001', 450000, 50, 'hEX S Gigabit Ethernet Router', 'active', 'in', 'Warehouse', '2024-01-15', NOW(), NOW());

-- Create a view for inventory summary
CREATE OR REPLACE VIEW `inventory_summary` AS
SELECT 
    a.id as asset_id,
    a.brand,
    a.type,
    a.model,
    a.serial_number as product_code,
    COUNT(ai.id) as total_items,
    SUM(CASE WHEN ai.status = 'in_stock' THEN 1 ELSE 0 END) as in_stock_count,
    SUM(CASE WHEN ai.status = 'in_use' THEN 1 ELSE 0 END) as in_use_count,
    SUM(CASE WHEN ai.status = 'maintenance' THEN 1 ELSE 0 END) as maintenance_count,
    SUM(CASE WHEN ai.status = 'damaged' THEN 1 ELSE 0 END) as damaged_count,
    SUM(CASE WHEN ai.status = 'retired' THEN 1 ELSE 0 END) as retired_count
FROM `assets` a
LEFT JOIN `asset_items` ai ON a.id = ai.asset_id
GROUP BY a.id, a.brand, a.type, a.model, a.serial_number;

-- Create a view for transaction history
CREATE OR REPLACE VIEW `asset_transaction_history` AS
SELECT 
    'installation' as transaction_source,
    at.id,
    at.customer_installation_id as reference_id,
    at.asset_item_id,
    ai.mac_address,
    a.brand,
    a.model,
    at.transaction_type,
    at.notes,
    at.created_by,
    u.name as created_by_name,
    at.created_at
FROM `asset_transactions` at
JOIN `asset_items` ai ON at.asset_item_id = ai.id
JOIN `assets` a ON ai.asset_id = a.id
LEFT JOIN `users` u ON at.created_by = u.id

UNION ALL

SELECT 
    'trouble_ticket' as transaction_source,
    tat.id,
    CAST(tat.trouble_ticket_id AS CHAR) as reference_id,
    tat.asset_item_id,
    ai.mac_address,
    a.brand,
    a.model,
    tat.transaction_type,
    tat.notes,
    tat.created_by,
    u.name as created_by_name,
    tat.created_at
FROM `ticket_asset_transactions` tat
JOIN `asset_items` ai ON tat.asset_item_id = ai.id
JOIN `assets` a ON ai.asset_id = a.id
LEFT JOIN `users` u ON tat.created_by = u.id

ORDER BY created_at DESC;

-- Show completion message
SELECT 'Inventory Management System migration completed successfully!' as message;

