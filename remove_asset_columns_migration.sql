-- Migration to remove deprecated columns from assets table
-- This aligns with the new inventory management system design

USE `iqgncnzy_skripsi`;

-- Remove the deprecated columns from assets table
ALTER TABLE `assets` 
DROP COLUMN IF EXISTS `status`,
DROP COLUMN IF EXISTS `status_in_out`, 
DROP COLUMN IF EXISTS `quantity`;

-- Show completion message
SELECT 'Asset table columns removed successfully!' as message;
