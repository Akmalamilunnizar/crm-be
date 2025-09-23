-- Migration for Recurring Invoices
-- This migration adds support for recurring invoices similar to the reference system

-- Create recurring_invoices table
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
  `created_by` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  PRIMARY KEY (`id`),
  KEY `recurring_invoices_customer_id_fkey` (`customer_id`),
  KEY `recurring_invoices_created_by_fkey` (`created_by`),
  KEY `idx_recurring_invoices_status` (`status`),
  KEY `idx_recurring_invoices_next_invoice_date` (`next_invoice_date`),
  CONSTRAINT `recurring_invoices_customer_id_fkey` FOREIGN KEY (`customer_id`) REFERENCES `customer` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `recurring_invoices_created_by_fkey` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
