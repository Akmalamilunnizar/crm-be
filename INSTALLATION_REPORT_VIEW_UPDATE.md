# Installation Report View Update Guide

## Problem
The `installation_report_complete` database view was returning empty values for network device and product fields because it wasn't properly joining the necessary tables.

## Solution
The view has been updated to include proper JOINs to the following tables:
- `network_devices` - for network device information
- `assets` - for router/device information  
- `products` - for product/package information
- `customer_services` - for customer service details

## Files Updated

### 1. Frontend: `crm-fe/pages/dashboard/report/customer-installation/detail/[id].vue`
- Updated the `Report` interface to include all fields from the updated view
- Added missing fields: `customer_id`, `technician_id`, `installation_updated_at`, `network_device_id`, `product_id`, `customer_service_id`

### 2. Backend: `crm-be/create_view.go`
- Updated the view definition to include proper JOINs
- Changed from empty strings/NULL values to actual data from joined tables
- Added `WHERE ci.deleted_at IS NULL` to filter out soft-deleted records

## How to Apply the Update

### Option 1: Run the SQL File Directly
Execute the SQL file `update-installation-report-view.sql` in your MySQL database:

```bash
mysql -h [host] -u [username] -p [database_name] < update-installation-report-view.sql
```

### Option 2: Run the Go Program
```bash
cd crm-be
go run create_view.go
```

### Option 3: Manual SQL Execution
Copy the SQL from `update-installation-report-view.sql` and execute it directly in your MySQL client (phpMyAdmin, MySQL Workbench, etc.)

## Verification

After updating the view, test by:
1. Navigate to the installation report detail page in the frontend
2. Check that all fields are populated (network device, product, router info, etc.)
3. Verify that PSB (Pasang Baru) calculations are working correctly

## Key Changes in the View

**Before:**
- Network device fields were empty strings: `'' AS network_device_id`
- Product fields were empty/NULL
- Router fields were empty

**After:**
- Network device fields from actual data: `nd.id AS network_device_id`
- Product fields from products table: `p.name AS product_name`, `p.price AS product_price`
- Router fields from assets table: `a.brand AS router_brand`, `a.model AS router_model`

## Database View Structure

The updated view now includes:
- Customer information (from `customer` table)
- Technician information (from `users` table)
- Installation details (from `customer_installations` table)
- Network device details (from `network_devices` table)
- Router/Asset details (from `assets` table)
- Product/Package details (from `products` table)
- Customer service details (from `customer_services` table)
- Computed PSB fields (durasi_psb, status_psb)

## Notes
- The view uses `LEFT JOIN` to ensure installations without network devices/products still appear
- The `WHERE ci.deleted_at IS NULL` clause filters out soft-deleted installations
- PSB duration is calculated using `TO_DAYS()` function for accurate day differences
