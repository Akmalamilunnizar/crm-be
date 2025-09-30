-- Migration to support multiple products per customer
-- This allows customers to have multiple network devices with different internet speeds

-- Step 1: Add product_id to network_devices table
ALTER TABLE `network_devices` 
ADD COLUMN `product_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL AFTER `assets_id`;

-- Step 2: Add foreign key constraint for product_id in network_devices
ALTER TABLE `network_devices` 
ADD CONSTRAINT `network_devices_ibfk_3` 
FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- Step 3: Add index for better performance
ALTER TABLE `network_devices` 
ADD KEY `idx_network_devices_product_id` (`product_id`);

-- Step 4: Migrate existing data (copy product_id from customer to network_devices)
-- Note: Since customer table doesn't have product_id column, we'll skip this step
-- You may need to manually set product_id for existing network_devices based on your business logic
-- UPDATE `network_devices` nd
-- JOIN `customer` c ON nd.customer_id = c.id
-- SET nd.product_id = c.product_id
-- WHERE nd.product_id IS NULL;

-- Step 5: Remove product_id from customer table (commented out for safety)
-- Uncomment the following lines after verifying the migration worked correctly:
-- ALTER TABLE `customer` DROP FOREIGN KEY `customer_ibfk_3`;
-- ALTER TABLE `customer` DROP KEY `customer_ibfk_3`;
-- ALTER TABLE `customer` DROP KEY `idx_customer_product_id`;
-- ALTER TABLE `customer` DROP COLUMN `product_id`;
