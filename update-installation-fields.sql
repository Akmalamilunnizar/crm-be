-- Update customer_installations table with all required installation fields
USE iqgncnzy_skripsi;

-- Add missing fields to customer_installations table
ALTER TABLE customer_installations 
ADD COLUMN IF NOT EXISTS installation_date DATE DEFAULT NULL COMMENT 'Installation date',
ADD COLUMN IF NOT EXISTS phone VARCHAR(20) DEFAULT NULL COMMENT 'Installation team phone number',
ADD COLUMN IF NOT EXISTS on_air_date DATE DEFAULT NULL COMMENT 'Date when service went live',
ADD COLUMN IF NOT EXISTS date_trial DATE DEFAULT NULL COMMENT 'Trial period end date',
ADD COLUMN IF NOT EXISTS service_activation_date DATE DEFAULT NULL COMMENT 'Service activation date',
ADD COLUMN IF NOT EXISTS switch_id VARCHAR(191) DEFAULT NULL COMMENT 'Switch ID',
ADD COLUMN IF NOT EXISTS port_number VARCHAR(50) DEFAULT NULL COMMENT 'Port number',
ADD COLUMN IF NOT EXISTS cable_type VARCHAR(100) DEFAULT NULL COMMENT 'Cable type',
ADD COLUMN IF NOT EXISTS cable_length DECIMAL(10,2) DEFAULT NULL COMMENT 'Cable length in meters',
ADD COLUMN IF NOT EXISTS end_port_type VARCHAR(50) DEFAULT NULL COMMENT 'End port type',
ADD COLUMN IF NOT EXISTS router_brand VARCHAR(100) DEFAULT NULL COMMENT 'Router brand',
ADD COLUMN IF NOT EXISTS router_model VARCHAR(100) DEFAULT NULL COMMENT 'Router model/series',
ADD COLUMN IF NOT EXISTS mac_address VARCHAR(191) DEFAULT NULL COMMENT 'MAC address',
ADD COLUMN IF NOT EXISTS eth_port VARCHAR(50) DEFAULT NULL COMMENT 'Ethernet port',
ADD COLUMN IF NOT EXISTS ip_static VARCHAR(45) DEFAULT NULL COMMENT 'Static IP address',
ADD COLUMN IF NOT EXISTS remote_port VARCHAR(50) DEFAULT NULL COMMENT 'Remote port',
ADD COLUMN IF NOT EXISTS user_login VARCHAR(191) DEFAULT NULL COMMENT 'User login',
ADD COLUMN IF NOT EXISTS password VARCHAR(191) DEFAULT NULL COMMENT 'Password',
ADD COLUMN IF NOT EXISTS status_user ENUM('Active','Inactive','Suspended','Pending') DEFAULT 'Active' COMMENT 'User status',
ADD COLUMN IF NOT EXISTS status_perangkat ENUM('active','inactive','maintenance','faulty') DEFAULT 'active' COMMENT 'Device status',
ADD COLUMN IF NOT EXISTS kepemilikan_perangkat ENUM('owned','leased','customer') DEFAULT 'owned' COMMENT 'Device ownership',
ADD COLUMN IF NOT EXISTS last_ping_status ENUM('up','down','unknown') DEFAULT 'unknown' COMMENT 'Last ping status';

-- Show the updated table structure
DESCRIBE customer_installations;

















