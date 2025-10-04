-- Remove old columns from customer table that are no longer needed
USE iqgncnzy_skripsi;

-- Remove old columns one by one
ALTER TABLE customer DROP COLUMN card_identition;
ALTER TABLE customer DROP COLUMN type_of_service;
ALTER TABLE customer DROP COLUMN email;
ALTER TABLE customer DROP COLUMN company_id;
ALTER TABLE customer DROP COLUMN gender;
ALTER TABLE customer DROP COLUMN no_identition;
ALTER TABLE customer DROP COLUMN job;
ALTER TABLE customer DROP COLUMN status_user;

-- Show the updated table structure
DESCRIBE customer;


















