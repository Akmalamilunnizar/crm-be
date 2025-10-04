-- Add new fields to customer table for registration data
USE iqgncnzy_skripsi;

-- Add new columns to customer table
ALTER TABLE customer 
ADD COLUMN service_request_date DATE DEFAULT NULL COMMENT 'Date when customer requested service',
ADD COLUMN proposed_package VARCHAR(191) DEFAULT NULL COMMENT 'Package proposed by customer',
ADD COLUMN bandwidth_capacity VARCHAR(191) DEFAULT NULL COMMENT 'Bandwidth capacity requested';

-- Remove old columns that are no longer needed
ALTER TABLE customer 
DROP COLUMN IF EXISTS type_of_service,
DROP COLUMN IF EXISTS email,
DROP COLUMN IF EXISTS company_id,
DROP COLUMN IF EXISTS gender,
DROP COLUMN IF EXISTS card_identition,
DROP COLUMN IF EXISTS no_identition,
DROP COLUMN IF EXISTS job,
DROP COLUMN IF EXISTS password,
DROP COLUMN IF EXISTS status_user;

-- Show the updated table structure
DESCRIBE customer;

















