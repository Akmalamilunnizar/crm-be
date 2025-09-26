-- Test database view and data

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
    createdAt
FROM customer_installations 
ORDER BY createdAt DESC 
LIMIT 3;

-- 3. Check if installation_report_complete view exists
SELECT 'installation_report_complete' as view_name, COUNT(*) as record_count FROM installation_report_complete;

-- 4. Show sample data from view
SELECT 
    installation_id,
    customer_name,
    technician_name,
    installation_status,
    installation_type,
    on_air_date
FROM installation_report_complete 
ORDER BY installation_created_at DESC 
LIMIT 3;

-- 5. Test the exact query used in repository
SELECT 
    installation_id,
    customer_id,
    customer_name,
    customer_address,
    customer_phone,
    technician_id,
    technician_name,
    technician_phone,
    installation_status,
    installation_type,
    installation_notes,
    on_air_date,
    trial_end_date,
    service_ready_date,
    installation_completed_at,
    document_type,
    document_photo,
    total_assets_out,
    total_assets_in,
    network_device_id,
    switch_id,
    port_number,
    remote_port,
    eth_port,
    mac_address,
    ip_static,
    status_perangkat,
    kepemilikan_perangkat,
    last_ping_status,
    last_ping_timestamp,
    router_brand,
    router_type,
    router_model,
    router_serial,
    customer_service_id,
    user_login,
    password,
    user_status,
    service_notes,
    cable_id,
    cable_name,
    cable_type,
    cable_length,
    cable_status,
    end_port_type,
    installation_created_at,
    installation_updated_at
FROM installation_report_complete
ORDER BY installation_created_at DESC
LIMIT 1;
