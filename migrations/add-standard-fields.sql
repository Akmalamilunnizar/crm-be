-- Simple script to add standard fields for customer registration and installation
-- Run this script in your MySQL client (phpMyAdmin, MySQL Workbench, or command line)

-- 1. Add sales_representative_id to customer table
ALTER TABLE customer 
ADD COLUMN sales_representative_id VARCHAR(191) NULL 
COMMENT 'Sales representative who handled this customer';

-- 2. Add foreign key constraint for sales_representative_id
ALTER TABLE customer 
ADD CONSTRAINT fk_customer_sales_representative 
FOREIGN KEY (sales_representative_id) REFERENCES users(id) 
ON DELETE SET NULL ON UPDATE CASCADE;

-- 3. Add on_air_date to customer_installations table
ALTER TABLE customer_installations 
ADD COLUMN on_air_date DATE NULL 
COMMENT 'Date when service first went live/on-air';

-- 4. Add password and user_status to customer_services table
ALTER TABLE customer_services 
ADD COLUMN password VARCHAR(191) NULL 
COMMENT 'Customer service login password';

ALTER TABLE customer_services 
ADD COLUMN user_status ENUM('Active', 'Inactive', 'Suspended', 'Pending') 
DEFAULT 'Active' 
COMMENT 'Customer service user status';

-- 5. Add additional useful fields for better tracking
ALTER TABLE customer_services 
ADD COLUMN installation_notes TEXT NULL 
COMMENT 'Additional notes about the installation process';

ALTER TABLE customer_services 
ADD COLUMN installation_team_phone VARCHAR(20) NULL 
COMMENT 'Phone number of installation team';

-- 6. Add indexes for better performance
CREATE INDEX idx_customer_sales_representative_id ON customer(sales_representative_id);
CREATE INDEX idx_customer_installations_on_air_date ON customer_installations(on_air_date);
CREATE INDEX idx_customer_services_user_status ON customer_services(user_status);
CREATE INDEX idx_customer_services_installation_team_phone ON customer_services(installation_team_phone);

-- 7. Update existing records with default values
UPDATE customer_services 
SET user_status = 'Active' 
WHERE user_status IS NULL;

-- 8. Show success message
SELECT 'Standard fields added successfully!' as status;

