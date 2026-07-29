-- Add new fields to customer_installations table (removed installation_type as it's only for new installations)
ALTER TABLE `customer_installations` 
ADD COLUMN `status` VARCHAR(50) DEFAULT 'pending' AFTER `date`,
ADD COLUMN `equipment_used` TEXT AFTER `status`,
ADD COLUMN `notes` TEXT AFTER `equipment_used`,
ADD COLUMN `completion_time` INT DEFAULT 0 AFTER `notes`,
ADD COLUMN `customer_signature` TEXT AFTER `completion_time`;
