-- Test script to check database data and views

-- 1. Check if customer_installations table has data
SELECT 'customer_installations' as table_name, COUNT(*) as record_count FROM customer_installations;

-- 2. Check if installation_report_complete view exists and has data
SELECT 'installation_report_complete' as view_name, COUNT(*) as record_count FROM installation_report_complete;

-- 3. Show sample data from customer_installations
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
LIMIT 5;

-- 4. Show sample data from installation_report_complete view
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
LIMIT 5;

-- 5. Check if view exists
SHOW TABLES LIKE 'installation_report_complete';

-- 6. Check view definition
SHOW CREATE VIEW installation_report_complete;
