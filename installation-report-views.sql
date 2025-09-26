-- Script untuk membuat views untuk laporan instalasi yang lengkap
-- Menggabungkan semua data dari berbagai tabel yang berelasi dengan customer_installations

-- Step 1: Buat view untuk laporan instalasi lengkap
CREATE OR REPLACE VIEW `installation_report_complete` AS
SELECT 
    -- Basic Installation Info
    ci.id as installation_id,
    ci.customer_id,
    c.name as customer_name,
    c.address as customer_address,
    c.phone as customer_phone,
    ci.technician_id,
    u.name as technician_name,
    u.phone as technician_phone,
    
    -- Installation Details
    ci.status as installation_status,
    ci.installation_type,
    ci.notes as installation_notes,
    ci.on_air_date,
    ci.trial_end_date,
    ci.service_ready_date,
    ci.installation_completed_at,
    
    -- Document Info
    ci.document_type,
    ci.document_photo,
    
    -- Asset Info
    ci.total_assets_out,
    ci.total_assets_in,
    
    -- Network Device Info
    nd.id as network_device_id,
    nd.switch_id,
    nd.port_number,
    nd.remote_port,
    nd.eth_port,
    nd.mac_address,
    nd.ip_static,
    nd.status_perangkat,
    nd.kepemilikan_perangkat,
    nd.last_ping_status,
    nd.last_ping_timestamp,
    
    -- Asset Details
    a.brand as router_brand,
    a.type as router_type,
    a.model as router_model,
    a.serial_number as router_serial,
    
    -- Customer Service Info
    cs.id as customer_service_id,
    cs.user_login,
    cs.password,
    cs.user_status,
    cs.installation_notes as service_notes,
    cs.installation_team_phone,
    
    -- Cable Info
    cab.id as cable_id,
    cab.name as cable_name,
    cab.type as cable_type,
    cab.length as cable_length,
    cab.status as cable_status,
    
    -- End Port Type
    cs.end_port_type,
    
    -- Timestamps
    ci.createdAt as installation_created_at,
    ci.updatedAt as installation_updated_at
    
FROM customer_installations ci
LEFT JOIN customer c ON ci.customer_id = c.id
LEFT JOIN users u ON ci.technician_id = u.id
LEFT JOIN network_devices nd ON ci.id = nd.customer_installation_id
LEFT JOIN assets a ON nd.assets_id = a.id
LEFT JOIN customer_services cs ON ci.id = cs.customer_installation_id
LEFT JOIN cable cab ON ci.id = cab.customer_installation_id
ORDER BY ci.createdAt DESC;

-- Step 2: Buat view untuk summary instalasi per customer
CREATE OR REPLACE VIEW `installation_summary_per_customer` AS
SELECT 
    c.id as customer_id,
    c.name as customer_name,
    c.address as customer_address,
    c.phone as customer_phone,
    COUNT(ci.id) as total_installations,
    COUNT(CASE WHEN ci.status = 'completed' THEN 1 END) as completed_installations,
    COUNT(CASE WHEN ci.status = 'pending' THEN 1 END) as pending_installations,
    COUNT(CASE WHEN ci.status = 'in_progress' THEN 1 END) as in_progress_installations,
    MAX(ci.on_air_date) as latest_on_air_date,
    MAX(ci.installation_completed_at) as latest_completion_date
FROM customer c
LEFT JOIN customer_installations ci ON c.id = ci.customer_id
GROUP BY c.id, c.name, c.address, c.phone
ORDER BY total_installations DESC;

-- Step 3: Buat view untuk laporan aset per instalasi
CREATE OR REPLACE VIEW `installation_asset_report` AS
SELECT 
    ci.id as installation_id,
    c.name as customer_name,
    ci.installation_type,
    ci.status as installation_status,
    ci.on_air_date,
    ci.installation_completed_at,
    
    -- Asset Out Summary
    COUNT(CASE WHEN at.transaction_type = 'out' THEN 1 END) as total_assets_out,
    SUM(CASE WHEN at.transaction_type = 'out' THEN at.quantity ELSE 0 END) as total_quantity_out,
    
    -- Asset In Summary
    COUNT(CASE WHEN at.transaction_type = 'in' THEN 1 END) as total_assets_in,
    SUM(CASE WHEN at.transaction_type = 'in' THEN at.quantity ELSE 0 END) as total_quantity_in,
    
    -- Asset Details
    GROUP_CONCAT(DISTINCT 
        CASE WHEN at.transaction_type = 'out' THEN 
            CONCAT(a.brand, ' ', a.model, ' (', at.quantity, ' pcs)')
        END 
        SEPARATOR ', '
    ) as assets_out_details,
    
    GROUP_CONCAT(DISTINCT 
        CASE WHEN at.transaction_type = 'in' THEN 
            CONCAT(a.brand, ' ', a.model, ' (', at.quantity, ' pcs)')
        END 
        SEPARATOR ', '
    ) as assets_in_details
    
FROM customer_installations ci
LEFT JOIN customer c ON ci.customer_id = c.id
LEFT JOIN asset_transactions at ON ci.id = at.customer_installation_id
LEFT JOIN assets a ON at.asset_id = a.id
GROUP BY ci.id, c.name, ci.installation_type, ci.status, ci.on_air_date, ci.installation_completed_at
ORDER BY ci.createdAt DESC;

-- Step 4: Buat view untuk laporan teknisi per instalasi
CREATE OR REPLACE VIEW `installation_technician_report` AS
SELECT 
    u.id as technician_id,
    u.name as technician_name,
    u.phone as technician_phone,
    COUNT(ci.id) as total_installations,
    COUNT(CASE WHEN ci.status = 'completed' THEN 1 END) as completed_installations,
    COUNT(CASE WHEN ci.status = 'pending' THEN 1 END) as pending_installations,
    COUNT(CASE WHEN ci.status = 'in_progress' THEN 1 END) as in_progress_installations,
    AVG(DATEDIFF(ci.installation_completed_at, ci.createdAt)) as avg_completion_days,
    MAX(ci.installation_completed_at) as latest_completion_date
FROM users u
LEFT JOIN customer_installations ci ON u.id = ci.technician_id
WHERE u.role_id = (SELECT id FROM roles WHERE name = 'TECHNICIAN')
GROUP BY u.id, u.name, u.phone
ORDER BY total_installations DESC;

-- Step 5: Verifikasi views yang dibuat
SELECT 
    'installation_report_complete' as view_name,
    'Complete installation report with all related data' as description
UNION ALL
SELECT 
    'installation_summary_per_customer' as view_name,
    'Installation summary grouped by customer' as description
UNION ALL
SELECT 
    'installation_asset_report' as view_name,
    'Asset report per installation' as description
UNION ALL
SELECT 
    'installation_technician_report' as view_name,
    'Technician performance report' as description;

-- Step 6: Test view dengan sample data
SELECT 
    'Testing installation_report_complete' as test_info,
    COUNT(*) as total_records
FROM installation_report_complete;

SELECT 
    'Testing installation_summary_per_customer' as test_info,
    COUNT(*) as total_customers
FROM installation_summary_per_customer;

SELECT 
    'Testing installation_asset_report' as test_info,
    COUNT(*) as total_installations
FROM installation_asset_report;

SELECT 
    'Testing installation_technician_report' as test_info,
    COUNT(*) as total_technicians
FROM installation_technician_report;
