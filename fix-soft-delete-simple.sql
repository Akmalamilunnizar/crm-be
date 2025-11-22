-- Simple fix for installation_report_complete view to filter soft-deleted records
-- Run this on the iqgncnzy_skripsi database

-- Step 1: Drop the existing view
DROP VIEW IF EXISTS installation_report_complete;

-- Step 2: Create the new view with soft delete filter
CREATE VIEW installation_report_complete AS
SELECT
    ci.id AS installation_id,
    ci.customer_id,
    c.name AS customer_name,
    c.address AS customer_address,
    c.phone AS customer_phone,
    c.service_request_date AS tgl_permintaan_psb,
    ci.technician_id,
    u.name AS technician_name,
    u.phone AS technician_phone,
    ci.status AS installation_status,
    ci.installation_type,
    ci.notes AS installation_notes,
    ci.on_air_date,
    ci.trial_end_date,
    ci.service_ready_date,
    ci.installation_completed_at,
    CASE
        WHEN c.service_request_date IS NOT NULL AND ci.installation_completed_at IS NOT NULL
        THEN DATEDIFF(ci.installation_completed_at, c.service_request_date)
        ELSE NULL
    END AS durasi_psb,
    CASE
        WHEN c.service_request_date IS NOT NULL AND ci.installation_completed_at IS NOT NULL
        THEN CASE
            WHEN DATEDIFF(ci.installation_completed_at, c.service_request_date) <= 3
            THEN 'Tepat Waktu'
            ELSE 'Terlambat'
        END
        ELSE NULL
    END AS status_psb,
    ci.document_type,
    ci.document_photo,
    nd.id AS network_device_id,
    nd.switch_id,
    nd.port_number,
    nd.remote_port,
    nd.eth_port,
    nd.mac_address,
    nd.ip_static,
    nd.kepemilikan_perangkat,
    a.brand AS router_brand,
    a.type AS router_type,
    a.model AS router_model,
    a.serial_number AS router_serial,
    p.name AS product_name,
    p.description AS product_description,
    p.price AS product_price,
    p.download_speed_mbps,
    p.upload_speed_mbps,
    cs.id AS customer_service_id,
    cs.user_login,
    cs.password,
    cs.user_status,
    cs.installation_notes AS service_notes,
    cs.cable_type,
    cs.cable_length,
    cs.end_port_type,
    ci.createdAt AS installation_created_at,
    ci.updatedAt AS installation_updated_at,
    p.id AS product_id
FROM customer_installations ci
LEFT JOIN customer c ON ci.customer_id = c.id
LEFT JOIN users u ON ci.technician_id = u.id
LEFT JOIN network_devices nd ON ci.id = nd.customer_installation_id
LEFT JOIN assets a ON nd.assets_id = a.id
LEFT JOIN products p ON nd.product_id = p.id
LEFT JOIN customer_services cs ON ci.id = cs.customer_installation_id
WHERE ci.deleted_at IS NULL;
