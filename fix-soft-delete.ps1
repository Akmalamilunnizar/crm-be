# PowerShell script to update installation_report_complete view
# Run from crm-be directory: .\fix-soft-delete.ps1

Write-Host "Updating installation_report_complete view to filter soft-deleted records..." -ForegroundColor Yellow
Write-Host ""

# MySQL connection details
$mysql_host = "103.63.24.139"
$mysql_user = "iqgncnzy_skripsi"
$mysql_pass = "XhYJOWlwNgsk"
$mysql_db = "iqgncnzy_skripsi"

# Build connection string
$connection_string = "-h $mysql_host -u $mysql_user -p$mysql_pass $mysql_db"

# SQL commands (note: using ` instead of " to avoid escaping issues)
$sql_commands = @"
DROP VIEW IF EXISTS installation_report_complete;

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
"@

Write-Host "Executing SQL commands..." -ForegroundColor Cyan

try {
    # Execute the SQL commands
    $result = Invoke-Expression "mysql $connection_string -e '$sql_commands'" 2>&1

    if ($LASTEXITCODE -eq 0) {
        Write-Host ""
        Write-Host "SUCCESS: installation_report_complete view updated successfully!" -ForegroundColor Green
        Write-Host "Installation reports will now only show records where deleted_at IS NULL." -ForegroundColor Green
        Write-Host ""
        Write-Host "Both frontend pages will now automatically hide soft-deleted installation records:" -ForegroundColor Cyan
        Write-Host "- Customer Dashboard (crm-fe/pages/customer/index.vue)" -ForegroundColor White
        Write-Host "- Installation Reports (crm-fe/pages/dashboard/report/customer-installation/reports.vue)" -ForegroundColor White
    } else {
        Write-Host ""
        Write-Host "ERROR: Failed to update the view." -ForegroundColor Red
        Write-Host "Exit code: $LASTEXITCODE" -ForegroundColor Red
        if ($result) {
            Write-Host "MySQL output:" -ForegroundColor Red
            Write-Host $result -ForegroundColor Red
        }
    }
} catch {
    Write-Host ""
    Write-Host "EXCEPTION: $_" -ForegroundColor Red
}

Write-Host ""
Read-Host "Press Enter to continue"
