-- Migration: Add network_architecture field to technician_steps table
-- This field indicates which network architectures the step applies to
-- NULL = applies to all architectures
-- 'FTTH' = only for FTTH
-- 'HTB' = only for HTB
-- 'FTTH,HTB' = for both (comma-separated)

ALTER TABLE `technician_steps` 
ADD COLUMN `network_architecture` VARCHAR(50) NULL DEFAULT NULL 
COMMENT 'Network architecture(s) this step applies to. NULL = all, FTTH = FTTH only, HTB = HTB only, comma-separated for multiple';

-- Update existing steps:
-- Steps 2, 3, 4, 6 are HTB/Media Converter only
UPDATE `technician_steps` SET `network_architecture` = 'HTB' WHERE `step_order` IN (2, 3, 4, 6);

-- All other steps (including step 0 selfie) apply to all architectures (NULL)
-- This is already the default, so no update needed

