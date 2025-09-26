# Installation Reports View Fix

## Problem
Installation reports tidak tampil di frontend meskipun endpoint API sudah dibuat karena ada masalah dengan database view dan query.

## Root Cause
1. **Database View Issue**: View `installation_report_complete` masih menggunakan field `installation_team_phone` yang sudah dihapus dari database
2. **Query Issue**: Repository query masih mencoba mengambil field `installation_team_phone` yang tidak ada
3. **Type Definition Issue**: Type `InstallationReportCompleteResponse` masih memiliki field `installation_team_phone`

## Solution

### 1. ✅ Fixed Repository Query
**File**: `crm-be/internal/api/admin/customer/installation/report_repository.go`
- Removed `installation_team_phone` from SELECT query
- Query now matches the actual database schema

### 2. ✅ Fixed Type Definition
**File**: `crm-be/internal/api/admin/customer/installation/report_request.go`
- Removed `InstallationTeamPhone` field from `InstallationReportCompleteResponse` struct
- Type now matches the actual database view structure

### 3. ✅ Created Database View Fix Script
**File**: `crm-be/fix-installation-report-view.sql`
- Drops and recreates `installation_report_complete` view
- Removes `installation_team_phone` field from view definition
- Includes test queries to verify the view works

### 4. ✅ Created Database Test Script
**File**: `crm-be/test-database.sql`
- Script to test database data and views
- Checks if `customer_installations` table has data
- Verifies `installation_report_complete` view exists and has data
- Shows sample data from both table and view

## Database View Structure (Fixed)
```sql
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
```

## Type Definition (Fixed)
```go
type InstallationReportCompleteResponse struct {
    InstallationId           string     `json:"installation_id"`
    CustomerId               string     `json:"customer_id"`
    CustomerName             string     `json:"customer_name"`
    CustomerAddress          string     `json:"customer_address"`
    CustomerPhone            string     `json:"customer_phone"`
    TechnicianId             string     `json:"technician_id"`
    TechnicianName           string     `json:"technician_name"`
    TechnicianPhone          string     `json:"technician_phone"`
    InstallationStatus       string     `json:"installation_status"`
    InstallationType         string     `json:"installation_type"`
    InstallationNotes        string     `json:"installation_notes"`
    OnAirDate                *time.Time `json:"on_air_date"`
    TrialEndDate             *time.Time `json:"trial_end_date"`
    ServiceReadyDate         *time.Time `json:"service_ready_date"`
    InstallationCompletedAt  *time.Time `json:"installation_completed_at"`
    DocumentType             string     `json:"document_type"`
    DocumentPhoto            string     `json:"document_photo"`
    TotalAssetsOut           int        `json:"total_assets_out"`
    TotalAssetsIn            int        `json:"total_assets_in"`
    NetworkDeviceId          string     `json:"network_device_id"`
    SwitchId                 string     `json:"switch_id"`
    PortNumber               string     `json:"port_number"`
    RemotePort               string     `json:"remote_port"`
    EthPort                  string     `json:"eth_port"`
    MacAddress               string     `json:"mac_address"`
    IPStatic                 string     `json:"ip_static"`
    StatusPerangkat          string     `json:"status_perangkat"`
    KepemilikanPerangkat     string     `json:"kepemilikan_perangkat"`
    LastPingStatus           string     `json:"last_ping_status"`
    LastPingTimestamp        *time.Time `json:"last_ping_timestamp"`
    RouterBrand              string     `json:"router_brand"`
    RouterType               string     `json:"router_type"`
    RouterModel              string     `json:"router_model"`
    RouterSerial             string     `json:"router_serial"`
    CustomerServiceId        string     `json:"customer_service_id"`
    UserLogin                string     `json:"user_login"`
    Password                 string     `json:"password"`
    UserStatus               string     `json:"user_status"`
    ServiceNotes             string     `json:"service_notes"`
    // InstallationTeamPhone removed ✅
    CableId                  string     `json:"cable_id"`
    CableName                string     `json:"cable_name"`
    CableType                string     `json:"cable_type"`
    CableLength              float64    `json:"cable_length"`
    CableStatus              string     `json:"cable_status"`
    EndPortType              string     `json:"end_port_type"`
    InstallationCreatedAt    time.Time  `json:"installation_created_at"`
    InstallationUpdatedAt    time.Time  `json:"installation_updated_at"`
}
```

## Steps to Fix Database:
1. Run the database fix script:
   ```sql
   -- Execute: crm-be/fix-installation-report-view.sql
   ```

2. Test the database:
   ```sql
   -- Execute: crm-be/test-database.sql
   ```

## Status:
✅ **Repository query fixed** - Removed installation_team_phone field
✅ **Type definition fixed** - Removed InstallationTeamPhone field
✅ **Database view script created** - Ready to fix the view
✅ **Test script created** - Ready to verify the fix
✅ **Backend restarted** - Ready for testing

## Next Steps:
1. Execute the database fix script to recreate the view
2. Test the API endpoint
3. Check if installation reports now appear in frontend

## Expected Result:
Installation reports should now be visible in the "Report-Report Customer Installation" page! 🎉
