-- Broadcast History Table
-- This table stores the history of all broadcast messages sent

CREATE TABLE IF NOT EXISTS `broadcast_history` (
  `id` VARCHAR(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `message` TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'The broadcast message content',
  `target_group` VARCHAR(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Target group: customers or team',
  `recipient_count` INT NOT NULL DEFAULT 0 COMMENT 'Number of recipients',
  `recipient_phones` JSON DEFAULT NULL COMMENT 'Array of phone numbers that received this broadcast',
  `status` ENUM('sent', 'failed', 'pending') NOT NULL DEFAULT 'pending' COMMENT 'Status of the broadcast',
  `template_key` VARCHAR(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Template key used (e.g., template_outage)',
  `sent_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'When the broadcast was sent',
  `created_by` VARCHAR(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'User ID who created the broadcast',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_target_group` (`target_group`),
  KEY `idx_status` (`status`),
  KEY `idx_sent_at` (`sent_at`),
  KEY `idx_created_by` (`created_by`),
  CONSTRAINT `broadcast_history_ibfk_1` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Broadcast message history';

