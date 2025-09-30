-- Simple test for installation reports

-- 1. Check customer_installations table
SELECT 'customer_installations' as table_name, COUNT(*) as record_count FROM customer_installations;

-- 2. Check if view exists and has data
SELECT 'installation_report_complete' as view_name, COUNT(*) as record_count FROM installation_report_complete;

-- 3. Show sample data from view
SELECT 
    installation_id,
    customer_name,
    technician_name,
    installation_status,
    installation_type,
    on_air_date
FROM installation_report_complete 
LIMIT 3;

-- 4. Test the exact query used in repository
SELECT 
    installation_id,
    customer_name,
    technician_name,
    installation_status,
    installation_type,
    on_air_date,
    router_brand,
    router_model,
    mac_address
FROM installation_report_complete
ORDER BY installation_created_at DESC
LIMIT 2;
