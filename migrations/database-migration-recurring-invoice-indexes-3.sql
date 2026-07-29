-- Add history generated_at index
CREATE INDEX `idx_recurring_invoice_history_generated_at` ON `recurring_invoice_history` (`generated_at`);
