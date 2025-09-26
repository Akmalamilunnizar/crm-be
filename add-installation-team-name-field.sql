-- Add installation team name field to customer_services table
-- This field is needed for "Nama Anggota Tim Install" in the form

ALTER TABLE `customer_services` 
ADD COLUMN `installation_team_name` VARCHAR(191) DEFAULT NULL COMMENT 'Nama anggota tim install' AFTER `installation_team_phone`;

-- Add index for better performance
CREATE INDEX `idx_customer_services_installation_team_name` ON `customer_services` (`installation_team_name`);
