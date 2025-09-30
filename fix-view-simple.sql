-- Fix installation_report_complete view

-- Drop existing view
DROP VIEW IF EXISTS `installation_report_complete`;

-- Create the view without installation_team_phone field
CREATE VIEW `installation_report_complete` AS
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

-- Test the view
SELECT 'View created successfully' as status;
SELECT COUNT(*) as total_records FROM installation_report_complete;
