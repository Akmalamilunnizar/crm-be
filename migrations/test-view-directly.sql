-- Test installation_report_complete view directly

-- 1. Check if view exists
SHOW TABLES LIKE 'installation_report_complete';

-- 2. Check view definition
SHOW CREATE VIEW installation_report_complete;

-- 3. Test simple query
SELECT COUNT(*) as total_records FROM installation_report_complete;

-- 4. Test with limit
SELECT * FROM installation_report_complete LIMIT 1;

-- 5. Check if there are any errors
SELECT 
    installation_id,
    customer_name,
    technician_name,
    installation_status,
    installation_type
FROM installation_report_complete 
LIMIT 3;
