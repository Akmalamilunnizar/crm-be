# Installation Reports Troubleshooting Guide

## Problem
Installation reports tidak tampil di frontend meskipun sudah dibuat dengan sukses.

## Current Status
✅ **Backend API endpoint created** - `/api/admin/customer-installation/report-complete`
✅ **Frontend API call implemented** - `getInstallationReportComplete()`
✅ **Type definitions added** - `InstallationReportCompleteResponse`
✅ **Repository query fixed** - Removed `installation_team_phone` field
✅ **Database view script created** - Ready to fix the view

## Root Cause Analysis
Masalah utama adalah **database view `installation_report_complete` masih menggunakan field `installation_team_phone` yang sudah dihapus dari database**.

## Step-by-Step Solution

### Step 1: Fix Database View
Jalankan script database untuk memperbaiki view:

**Option A: Manual SQL Execution**
```sql
-- Execute: crm-be/fix-installation-report-view.sql
```

**Option B: PowerShell Script**
```powershell
# Run: crm-be/run-database-fix.ps1
```

**Option C: Batch Script**
```batch
# Run: crm-be/run-database-fix.bat
```

### Step 2: Test Database
Jalankan script debug untuk memverifikasi:

```sql
-- Execute: crm-be/debug-installation-reports.sql
```

### Step 3: Start Backend
```bash
cd crm-be
go run cmd/myapp/main.go
```

### Step 4: Test API Endpoint
```bash
curl -X GET "http://localhost:8080/api/admin/customer-installation/report-complete" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

### Step 5: Test Frontend
1. Navigate to: `/dashboard/report/customer-installation/reports`
2. Check browser console for any errors
3. Verify API calls in Network tab

## Expected Results

### Database View Should Return:
```json
{
  "installation_id": "949189b8-97d1-11f0-9fa2-d843ae0f1e06",
  "customer_name": "Cha Eunwoo",
  "technician_name": "teknisi",
  "installation_status": "completed",
  "installation_type": "new_installation",
  "on_air_date": "2025-09-17T00:00:00Z",
  "router_brand": "ads",
  "router_model": "asd",
  "mac_address": "14:E2:O2:19:22:90"
}
```

### API Response Should Be:
```json
{
  "success": true,
  "message": "Installation reports retrieved successfully",
  "data": [
    {
      "installation_id": "949189b8-97d1-11f0-9fa2-d843ae0f1e06",
      "customer_name": "Cha Eunwoo",
      "technician_name": "teknisi",
      "installation_status": "completed",
      "installation_type": "new_installation",
      "on_air_date": "2025-09-17T00:00:00Z",
      "router_brand": "ads",
      "router_model": "asd",
      "mac_address": "14:E2:O2:19:22:90"
    }
  ]
}
```

### Frontend Should Display:
- Table with installation reports
- Customer names, technician names, status, etc.
- No "No installation reports found" message

## Troubleshooting Checklist

### ✅ Database Issues
- [ ] `customer_installations` table has data
- [ ] `installation_report_complete` view exists
- [ ] View definition doesn't include `installation_team_phone`
- [ ] View returns data when queried directly

### ✅ Backend Issues
- [ ] Backend server is running on port 8080
- [ ] API endpoint `/api/admin/customer-installation/report-complete` is accessible
- [ ] No errors in backend logs
- [ ] Database connection is working

### ✅ Frontend Issues
- [ ] API call `getInstallationReportComplete()` is working
- [ ] No JavaScript errors in browser console
- [ ] Network requests are successful
- [ ] Data is being received and processed

### ✅ Authentication Issues
- [ ] User is logged in
- [ ] Token is valid and not expired
- [ ] Authorization header is being sent correctly

## Common Issues and Solutions

### Issue 1: "No installation reports found"
**Cause**: Database view is not returning data
**Solution**: Run database fix script

### Issue 2: API returns 500 error
**Cause**: Database view has invalid field references
**Solution**: Recreate the view with correct field names

### Issue 3: Frontend shows loading forever
**Cause**: API endpoint is not accessible
**Solution**: Check if backend is running and endpoint exists

### Issue 4: Authentication errors
**Cause**: Token is invalid or expired
**Solution**: Re-login to get new token

## Files Modified
- `crm-be/internal/api/admin/customer/installation/report_repository.go` - Fixed query
- `crm-be/internal/api/admin/customer/installation/report_request.go` - Fixed type definition
- `crm-be/fix-installation-report-view.sql` - Database view fix script
- `crm-be/debug-installation-reports.sql` - Debug script
- `crm-be/run-database-fix.ps1` - PowerShell automation script

## Next Steps
1. **Execute database fix script** to recreate the view
2. **Start backend server** and verify it's running
3. **Test API endpoint** to ensure it returns data
4. **Refresh frontend** and check if reports appear
5. **Verify data display** in the reports table

## Expected Outcome
After completing all steps, installation reports should be visible in the frontend at `/dashboard/report/customer-installation/reports`! 🎉

