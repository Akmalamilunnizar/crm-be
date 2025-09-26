-- Netwatch Integration Database Migration
-- Run this script to add Netwatch monitoring tables to your database





-- Create ticket_logs table
CREATE TABLE IF NOT EXISTS ticket_logs (
    id VARCHAR(36) PRIMARY KEY,
    ticket_id BIGINT UNSIGNED NOT NULL,
    log_type ENUM('manual', 'netwatch', 'system') DEFAULT 'manual',
    message TEXT NOT NULL,
    event_id VARCHAR(36),
    created_by VARCHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (ticket_id) REFERENCES trouble_tickets(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);




-- Add indexes for better performance (idempotent)
SET @__idx1 := (SELECT COUNT(1) FROM INFORMATION_SCHEMA.STATISTICS 
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='ticket_logs' AND INDEX_NAME='idx_ticket_logs_ticket');
SET @__sql1 := IF(@__idx1=0, 'CREATE INDEX idx_ticket_logs_ticket ON ticket_logs(ticket_id)', 'SELECT 1');
PREPARE __stmt1 FROM @__sql1; EXECUTE __stmt1; DEALLOCATE PREPARE __stmt1;

SET @__idx2 := (SELECT COUNT(1) FROM INFORMATION_SCHEMA.STATISTICS 
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='ticket_logs' AND INDEX_NAME='idx_ticket_logs_type');
SET @__sql2 := IF(@__idx2=0, 'CREATE INDEX idx_ticket_logs_type ON ticket_logs(log_type)', 'SELECT 1');
PREPARE __stmt2 FROM @__sql2; EXECUTE __stmt2; DEALLOCATE PREPARE __stmt2;

-- Insert sample data for testing (optional)


-- Insert sample trouble type for network issues
INSERT INTO trouble_type (id, name) VALUES 
('network', 'Network Issue')
ON DUPLICATE KEY UPDATE name = name;

-- Technician team members table
CREATE TABLE IF NOT EXISTS `technician_team_members` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `ticket_id` BIGINT UNSIGNED NOT NULL,
  `user_id` VARCHAR(191) NOT NULL,
  `role` VARCHAR(50) NOT NULL,
  `created_at` DATETIME NULL,
  `updated_at` DATETIME NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_team_ticket` (`ticket_id`),
  CONSTRAINT `fk_team_ticket` FOREIGN KEY (`ticket_id`) REFERENCES `trouble_tickets` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_team_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Ticket troubleshooting steps table (max 7 steps enforced at application level)
CREATE TABLE IF NOT EXISTS `ticket_steps` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `ticket_id` BIGINT UNSIGNED NOT NULL,
  `step_order` INT NOT NULL,
  `description` TEXT NULL,
  `image_path` VARCHAR(191) NULL,
  `created_at` DATETIME NULL,
  `updated_at` DATETIME NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_steps_ticket` (`ticket_id`),
  UNIQUE KEY `uq_ticket_step_order` (`ticket_id`, `step_order`),
  CONSTRAINT `fk_steps_ticket` FOREIGN KEY (`ticket_id`) REFERENCES `trouble_tickets` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Ticket step images (multiple images per step)
CREATE TABLE IF NOT EXISTS `ticket_step_images` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `step_id` BIGINT UNSIGNED NOT NULL,
  `path` VARCHAR(191) NOT NULL,
  `created_at` DATETIME NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_step_images_step` (`step_id`),
  CONSTRAINT `fk_step_images_step` FOREIGN KEY (`step_id`) REFERENCES `ticket_steps` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Add CS verification fields to trouble_tickets (idempotent, works on MySQL variants without IF NOT EXISTS)
SET @__col1 := (SELECT COUNT(1) FROM INFORMATION_SCHEMA.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'trouble_tickets' AND COLUMN_NAME = 'verified_by_cs');
SET @__sqlc1 := IF(@__col1=0, 'ALTER TABLE trouble_tickets ADD COLUMN verified_by_cs TINYINT(1) DEFAULT 0', 'SELECT 1');
PREPARE __stmta FROM @__sqlc1; EXECUTE __stmta; DEALLOCATE PREPARE __stmta;

SET @__col2 := (SELECT COUNT(1) FROM INFORMATION_SCHEMA.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'trouble_tickets' AND COLUMN_NAME = 'verified_at');
SET @__sqlc2 := IF(@__col2=0, 'ALTER TABLE trouble_tickets ADD COLUMN verified_at DATETIME NULL', 'SELECT 1');
PREPARE __stmtb FROM @__sqlc2; EXECUTE __stmtb; DEALLOCATE PREPARE __stmtb;

-- Add technician completion tracking column
SET @__col3 := (SELECT COUNT(1) FROM INFORMATION_SCHEMA.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'trouble_tickets' AND COLUMN_NAME = 'technician_completed');
SET @__sqlc3 := IF(@__col3=0, 'ALTER TABLE trouble_tickets ADD COLUMN technician_completed TINYINT(1) DEFAULT 0', 'SELECT 1');
PREPARE __stmtc FROM @__sqlc3; EXECUTE __stmtc; DEALLOCATE PREPARE __stmtc;

-- Add network architecture type column
SET @__col4 := (SELECT COUNT(1) FROM INFORMATION_SCHEMA.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'trouble_tickets' AND COLUMN_NAME = 'network_architecture');
SET @__sqlc4 := IF(@__col4=0, 'ALTER TABLE trouble_tickets ADD COLUMN network_architecture VARCHAR(50) NULL', 'SELECT 1');
PREPARE __stmtd FROM @__sqlc4; EXECUTE __stmtd; DEALLOCATE PREPARE __stmtd;

-- Create technician steps table (predefined checklist)
CREATE TABLE IF NOT EXISTS `technician_steps` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `step_order` INT NOT NULL,
  `title` VARCHAR(255) NOT NULL,
  `description` TEXT NULL,
  `tools` TEXT NULL,
  `spare_parts` TEXT NULL,
  `procedure` TEXT NULL,
  `solution` TEXT NULL,
  `is_active` TINYINT(1) DEFAULT 1,
  `created_at` DATETIME NULL,
  `updated_at` DATETIME NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_technician_steps_order` (`step_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create ticket technician steps table (progress tracking)
CREATE TABLE IF NOT EXISTS `ticket_technician_steps` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `ticket_id` BIGINT UNSIGNED NOT NULL,
  `step_id` BIGINT UNSIGNED NOT NULL,
  `technician_id` VARCHAR(191) NOT NULL,
  `status` VARCHAR(50) DEFAULT 'pending',
  `notes` TEXT NULL,
  `spare_parts_used` TEXT NULL,
  `completed_at` DATETIME NULL,
  `created_at` DATETIME NULL,
  `updated_at` DATETIME NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_ticket_tech_steps_ticket` (`ticket_id`),
  INDEX `idx_ticket_tech_steps_step` (`step_id`),
  INDEX `idx_ticket_tech_steps_technician` (`technician_id`),
  UNIQUE KEY `uq_ticket_step_technician` (`ticket_id`, `step_id`, `technician_id`),
  CONSTRAINT `fk_ticket_tech_steps_ticket` FOREIGN KEY (`ticket_id`) REFERENCES `trouble_tickets` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_ticket_tech_steps_step` FOREIGN KEY (`step_id`) REFERENCES `technician_steps` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_ticket_tech_steps_technician` FOREIGN KEY (`technician_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create spare parts table
CREATE TABLE IF NOT EXISTS `spare_parts` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL,
  `description` TEXT NULL,
  `category` VARCHAR(100) NULL,
  `is_active` TINYINT(1) DEFAULT 1,
  `created_at` DATETIME NULL,
  `updated_at` DATETIME NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_spare_parts_category` (`category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;