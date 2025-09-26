-- Debug script untuk installation reports

-- 1. Check if customer_installations table has data
SELECT 'customer_installations' as table_name, COUNT(*) as record_count FROM customer_installations;

-- 2. Show sample data from customer_installations
SELECT 
    id,
    customer_id,
    technician_id,
    status,
    installation_type,
    on_air_date,
    createdAt,
    updatedAt
FROM customer_installations 
ORDER BY createdAt DESC 
LIMIT 3;

-- 3. Check if installation_report_complete view exists
SHOW TABLES LIKE 'installation_report_complete';

-- 4. Check view definition
SHOW CREATE VIEW installation_report_complete;

-- 5. Test the view directly
SELECT COUNT(*) as total_records FROM installation_report_complete;

-- 6. Show sample data from view
SELECT 
    installation_id,
    customer_name,
    technician_name,
    installation_status,
    installation_type,
    on_air_date,
    installation_created_at
FROM installation_report_complete 
ORDER BY installation_created_at DESC 
LIMIT 3;

-- 7. Check if there are any errors in the view
SELECT * FROM installation_report_complete LIMIT 1;

