-- Create recurring_invoice_history table to track generated invoices
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
  CONSTRAINT `recurring_invoice_history_recurring_invoice_id_fkey` FOREIGN KEY (`recurring_invoice_id`) REFERENCES `recurring_invoices` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `recurring_invoice_history_generated_invoice_id_fkey` FOREIGN KEY (`generated_invoice_id`) REFERENCES `invoices` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
