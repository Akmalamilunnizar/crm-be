-- Fix technician_team_members table foreign key constraint
-- This script creates the missing table with correct foreign key references

-- Drop the table if it exists (to avoid conflicts)
DROP TABLE IF EXISTS `technician_team_members`;

-- Create the technician_team_members table with correct foreign key constraint
CREATE TABLE `technician_team_members` (
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

