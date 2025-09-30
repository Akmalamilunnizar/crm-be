-- Database Migration: Adding Standard Fields for Customer Registration and Installation
-- Date: 2025-01-18
-- Description: Adding missing standard fields for customer registration and installation process

-- 1. Add sales_representative_id to customer table
ALTER TABLE customer 
ADD COLUMN sales_representative_id VARCHAR(191) NULL 
COMMENT 'Sales representative who handled this customer';

-- Add foreign key constraint for sales_representative_id
ALTER TABLE customer 
ADD CONSTRAINT fk_customer_sales_representative 
FOREIGN KEY (sales_representative_id) REFERENCES users(id) 
ON DELETE SET NULL ON UPDATE CASCADE;

-- Add index for better performance
CREATE INDEX idx_customer_sales_representative_id ON customer(sales_representative_id);

-- 2. Add on_air_date to customer_installations table
ALTER TABLE customer_installations 
ADD COLUMN on_air_date DATE NULL 
COMMENT 'Date when service first went live/on-air';

-- Add index for better performance
CREATE INDEX idx_customer_installations_on_air_date ON customer_installations(on_air_date);

-- 3. Add password and user_status to customer_services table
ALTER TABLE customer_services 
ADD COLUMN password VARCHAR(191) NULL 
COMMENT 'Customer service login password';

ALTER TABLE customer_services 
ADD COLUMN user_status ENUM('Active', 'Inactive', 'Suspended', 'Pending') 
DEFAULT 'Active' 
COMMENT 'Customer service user status';

-- Add index for better performance
CREATE INDEX idx_customer_services_user_status ON customer_services(user_status);

-- 4. Add additional useful fields for better tracking
ALTER TABLE customer_services 
ADD COLUMN installation_notes TEXT NULL 
COMMENT 'Additional notes about the installation process';

ALTER TABLE customer_services 
ADD COLUMN installation_team_phone VARCHAR(20) NULL 
COMMENT 'Phone number of installation team';

-- Add index for better performance
CREATE INDEX idx_customer_services_installation_team_phone ON customer_services(installation_team_phone);

-- 5. Update existing records with default values where needed
UPDATE customer_services 
SET user_status = 'Active' 
WHERE user_status IS NULL;

-- 6. Add comments to existing relevant columns for better documentation
ALTER TABLE customer 
MODIFY COLUMN password VARCHAR(191) NOT NULL 
COMMENT 'Customer account password for login';

ALTER TABLE customer 
MODIFY COLUMN status_user ENUM('Active','Nonactive') NOT NULL DEFAULT 'Active' 
COMMENT 'Customer account status (Active/Nonactive)';

-- 7. Create a view for easy access to complete customer installation data
CREATE OR REPLACE VIEW v_customer_installation_complete AS
SELECT 
    c.id as customer_id,
    c.name as customer_name,
    c.address,
    c.phone as customer_phone,
    c.alias,
    c.latitude,
    c.longitude,
    c.installation_date,
    c.status_user as customer_status,
    c.password as customer_password,
    c.sales_representative_id,
    u_sales.name as sales_representative_name,
    u_sales.phone as sales_representative_phone,
    ci.id as installation_id,
    ci.date_trial,
    ci.on_air_date,
    ci.status as installation_status,
    ci.equipment_used,
    ci.notes as installation_notes,
    ci.completion_time,
    ci.customer_signature,
    cs.id as service_id,
    cs.user_login,
    cs.password as service_password,
    cs.user_status as service_user_status,
    cs.service_activation_date,
    cs.cable_length,
    cs.end_port_type,
    cs.installation_team_phone,
    nd.id as device_id,
    nd.mac_address,
    nd.ip_static,
    nd.port_number,
    nd.remote_port,
    nd.eth_port,
    nd.status_perangkat as device_status,
    nd.kepemilikan_perangkat as device_ownership,
    nd.last_ping_status,
    nd.last_ping_timestamp,
    a.brand as router_brand,
    a.model as router_model,
    a.serial_number,
    s.name as switch_name,
    s.model as switch_model,
    s.ip_address as switch_ip,
    s.location as switch_location,
    cab.type as cable_type,
    cab.name as cable_name,
    p.name as product_name,
    p.price as product_price,
    p.description as product_description,
    u_tech.name as technician_name,
    u_tech.phone as technician_phone
FROM customer c
LEFT JOIN users u_sales ON c.sales_representative_id = u_sales.id
LEFT JOIN customer_installations ci ON c.id = ci.customer_id
LEFT JOIN customer_details cd ON ci.details_id = cd.id
LEFT JOIN users u_tech ON cd.technician_id = u_tech.id
LEFT JOIN customer_services cs ON c.id = cs.customer_id
LEFT JOIN network_devices nd ON c.id = nd.customer_id
LEFT JOIN assets a ON nd.assets_id = a.id
LEFT JOIN switch s ON nd.switch_id = s.id
LEFT JOIN cable cab ON cs.cable_id = cab.id
LEFT JOIN products p ON nd.product_id = p.id;

-- 8. Create indexes for better performance on the view
CREATE INDEX idx_customer_installation_complete_customer_id ON customer_installations(customer_id);
CREATE INDEX idx_customer_installation_complete_technician_id ON customer_details(technician_id);
CREATE INDEX idx_customer_services_customer_id ON customer_services(customer_id);
CREATE INDEX idx_network_devices_customer_id ON network_devices(customer_id);

-- 9. Add sample data for testing (optional)
-- INSERT INTO customer (id, name, address, area_id, card_identition, company_id, email, gender, status_user, job, latitude, longitude, alias, no_identition, password, phone, type_of_service, installation_date, next_payment_date, sales_representative_id)
-- VALUES ('test-customer-001', 'Test Customer', 'Test Address', '61362c06-7042-4ed4-a391-237b742912c3', 'ktp', '53140661-9c49-423d-a00e-c8fe21680213', 'test@example.com', 'male', 'Active', 'Employee', 0.0, 0.0, 'TestAlias', 123456789, 'password123', '081234567890', 'internet', '2025-01-18', '2025-02-18', 'c7d376ca-626c-4e33-bf47-342d2082c93d');

-- 10. Create a stored procedure to get complete customer installation data
DELIMITER //
CREATE PROCEDURE GetCustomerInstallationData(IN customer_id_param VARCHAR(191))
BEGIN
    SELECT * FROM v_customer_installation_complete 
    WHERE customer_id = customer_id_param;
END //
DELIMITER ;

-- 11. Create a stored procedure to update customer installation status
DELIMITER //
CREATE PROCEDURE UpdateInstallationStatus(
    IN installation_id_param VARCHAR(191),
    IN new_status VARCHAR(50),
    IN on_air_date_param DATE,
    IN completion_time_param INT
)
BEGIN
    UPDATE customer_installations 
    SET 
        status = new_status,
        on_air_date = on_air_date_param,
        completion_time = completion_time_param,
        updatedAt = NOW()
    WHERE id = installation_id_param;
END //
DELIMITER ;

-- 12. Create a trigger to automatically update on_air_date when status changes to 'completed'
DELIMITER //
CREATE TRIGGER tr_update_on_air_date
BEFORE UPDATE ON customer_installations
FOR EACH ROW
BEGIN
    IF NEW.status = 'completed' AND OLD.status != 'completed' AND NEW.on_air_date IS NULL THEN
        SET NEW.on_air_date = CURDATE();
    END IF;
END //
DELIMITER ;

-- 13. Add validation constraints
ALTER TABLE customer_services 
ADD CONSTRAINT chk_user_status 
CHECK (user_status IN ('Active', 'Inactive', 'Suspended', 'Pending'));

ALTER TABLE customer_installations 
ADD CONSTRAINT chk_on_air_date 
CHECK (on_air_date IS NULL OR on_air_date >= date_trial);

-- 14. Create a function to calculate installation duration
DELIMITER //
CREATE FUNCTION CalculateInstallationDuration(installation_id_param VARCHAR(191))
RETURNS INT
READS SQL DATA
DETERMINISTIC
BEGIN
    DECLARE duration INT DEFAULT 0;
    
    SELECT completion_time INTO duration
    FROM customer_installations
    WHERE id = installation_id_param;
    
    RETURN IFNULL(duration, 0);
END //
DELIMITER ;

-- 15. Create a function to get customer service status
DELIMITER //
CREATE FUNCTION GetCustomerServiceStatus(customer_id_param VARCHAR(191))
RETURNS VARCHAR(20)
READS SQL DATA
DETERMINISTIC
BEGIN
    DECLARE service_status VARCHAR(20) DEFAULT 'Unknown';
    
    SELECT user_status INTO service_status
    FROM customer_services
    WHERE customer_id = customer_id_param
    LIMIT 1;
    
    RETURN IFNULL(service_status, 'Unknown');
END //
DELIMITER ;

-- 16. Add comments to the database schema
ALTER TABLE customer 
COMMENT = 'Customer registration and basic information';

ALTER TABLE customer_installations 
COMMENT = 'Customer installation process and tracking';

ALTER TABLE customer_services 
COMMENT = 'Customer service configuration and access details';

-- 17. Create a summary table for reporting
CREATE TABLE IF NOT EXISTS customer_installation_summary (
    id VARCHAR(191) PRIMARY KEY,
    customer_id VARCHAR(191) NOT NULL,
    total_installations INT DEFAULT 0,
    completed_installations INT DEFAULT 0,
    pending_installations INT DEFAULT 0,
    average_completion_time DECIMAL(10,2) DEFAULT 0,
    last_installation_date DATE NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customer(id) ON DELETE CASCADE
);

-- 18. Create a trigger to update summary table
DELIMITER //
CREATE TRIGGER tr_update_installation_summary
AFTER INSERT ON customer_installations
FOR EACH ROW
BEGIN
    INSERT INTO customer_installation_summary (id, customer_id, total_installations, completed_installations, pending_installations, last_installation_date)
    VALUES (CONCAT('summary-', NEW.customer_id), NEW.customer_id, 1, 
            CASE WHEN NEW.status = 'completed' THEN 1 ELSE 0 END,
            CASE WHEN NEW.status = 'pending' THEN 1 ELSE 0 END,
            NEW.date_trial)
    ON DUPLICATE KEY UPDATE
        total_installations = total_installations + 1,
        completed_installations = completed_installations + CASE WHEN NEW.status = 'completed' THEN 1 ELSE 0 END,
        pending_installations = pending_installations + CASE WHEN NEW.status = 'pending' THEN 1 ELSE 0 END,
        last_installation_date = NEW.date_trial,
        updated_at = NOW();
END //
DELIMITER ;

-- 19. Create indexes for better performance
CREATE INDEX idx_customer_installation_summary_customer_id ON customer_installation_summary(customer_id);
CREATE INDEX idx_customer_installation_summary_status ON customer_installation_summary(completed_installations, pending_installations);

-- 20. Add sample data for testing the new fields
-- UPDATE customer_services 
-- SET password = 'service123', user_status = 'Active', installation_team_phone = '081234567890'
-- WHERE customer_id IN (SELECT id FROM customer LIMIT 5);

-- UPDATE customer_installations 
-- SET on_air_date = DATE_ADD(date_trial, INTERVAL 1 DAY)
-- WHERE status = 'completed' AND on_air_date IS NULL;

-- 21. Create a view for sales representative performance
CREATE OR REPLACE VIEW v_sales_representative_performance AS
SELECT 
    u.id as sales_rep_id,
    u.name as sales_rep_name,
    u.phone as sales_rep_phone,
    u.email as sales_rep_email,
    COUNT(c.id) as total_customers,
    COUNT(CASE WHEN c.status_user = 'Active' THEN 1 END) as active_customers,
    COUNT(CASE WHEN c.status_user = 'Nonactive' THEN 1 END) as inactive_customers,
    COUNT(ci.id) as total_installations,
    COUNT(CASE WHEN ci.status = 'completed' THEN 1 END) as completed_installations,
    AVG(ci.completion_time) as avg_completion_time,
    MAX(c.createdAt) as last_customer_created
FROM users u
LEFT JOIN customer c ON u.id = c.sales_representative_id
LEFT JOIN customer_installations ci ON c.id = ci.customer_id
WHERE u.role_id = (SELECT id FROM roles WHERE name = 'ADMIN' OR name = 'CUSTOMER_SERVICE')
GROUP BY u.id, u.name, u.phone, u.email;

-- 22. Final comments and documentation
SELECT 'Database migration completed successfully!' as status,
       'Added sales_representative_id to customer table' as change1,
       'Added on_air_date to customer_installations table' as change2,
       'Added password and user_status to customer_services table' as change3,
       'Created comprehensive views and procedures' as change4,
       'Added performance indexes and constraints' as change5;

