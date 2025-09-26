# Customer Multiple Products Schema Migration

## Overview
This migration allows customers to have multiple network devices with different internet speeds/products, instead of being limited to a single product per customer.

## Problem with Current Schema
- `customer` table has `product_id` (one-to-one relationship)
- Each customer can only have one internet service/product
- Cannot support multiple devices with different speeds for the same customer

## New Schema Design

### Before (Current):
```
customer (1) -----> (1) product
customer (1) -----> (n) network_devices
```

### After (New):
```
customer (1) -----> (n) network_devices (n) -----> (1) product
```

## Database Changes

### 1. Add `product_id` to `network_devices` table
```sql
ALTER TABLE `network_devices` 
ADD COLUMN `product_id` varchar(191) DEFAULT NULL AFTER `assets_id`;
```

### 2. Add foreign key constraint
```sql
ALTER TABLE `network_devices` 
ADD CONSTRAINT `network_devices_ibfk_3` 
FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;
```

### 3. Add index for performance
```sql
ALTER TABLE `network_devices` 
ADD KEY `idx_network_devices_product_id` (`product_id`);
```

### 4. Migrate existing data
```sql
UPDATE `network_devices` nd
JOIN `customer` c ON nd.customer_id = c.id
SET nd.product_id = c.product_id
WHERE nd.product_id IS NULL;
```

### 5. Remove `product_id` from `customer` table (after verification)
```sql
ALTER TABLE `customer` DROP FOREIGN KEY `customer_ibfk_3`;
ALTER TABLE `customer` DROP KEY `customer_ibfk_3`;
ALTER TABLE `customer` DROP KEY `idx_customer_product_id`;
ALTER TABLE `customer` DROP COLUMN `product_id`;
```

## Benefits of New Schema

### 1. **Multiple Products per Customer**
- Customer can have multiple network devices
- Each device can have different internet speeds
- Example: Customer has both 100 Mbps and 50 Mbps connections

### 2. **Better Data Modeling**
- More accurate representation of real-world scenarios
- Network devices are the actual service endpoints
- Products are tied to specific devices, not customers

### 3. **Flexible Billing**
- Can bill different amounts for different devices
- Support for multiple service tiers per customer
- Better tracking of which device uses which service

## Example Use Cases

### Scenario 1: Business Customer
```
Customer: "PT Maju Bersama"
├── Device 1: Office Router (100 Mbps) - Product: "Business 100 Mbps"
├── Device 2: Warehouse Router (50 Mbps) - Product: "Standard 50 Mbps"
└── Device 3: Branch Office Router (25 Mbps) - Product: "Basic 25 Mbps"
```

### Scenario 2: Residential Customer with Multiple Services
```
Customer: "Cha Eunwoo"
├── Device 1: Main Router (100 Mbps) - Product: "Premium 100 Mbps"
└── Device 2: Guest Network Router (20 Mbps) - Product: "Basic 20 Mbps"
```

## Application Code Changes Required

### 1. Backend API Changes
- Update customer creation to not require `product_id`
- Update network device creation to require `product_id`
- Modify queries to join through `network_devices` table

### 2. Frontend Changes
- Update customer forms to not include product selection
- Update network device forms to include product selection
- Modify customer detail views to show multiple products

### 3. Billing System Updates
- Generate invoices per network device, not per customer
- Support multiple invoices per customer
- Update recurring invoice logic

## Migration Steps

### Step 1: Run Migration
```bash
# Set environment variables
export DB_HOST="localhost"
export DB_PORT="3306"
export DB_USER="root"
export DB_PASSWORD=""
export DB_NAME="iqgncnzy_skripsi"

# Run migration
go run migrate-customer-multiple-products.go
```

### Step 2: Verify Migration
```sql
-- Check that network_devices now have product_id
SELECT nd.id, nd.customer_id, nd.product_id, p.name as product_name
FROM network_devices nd
LEFT JOIN products p ON nd.product_id = p.id;

-- Verify data migration
SELECT c.id, c.name, c.product_id as old_product_id, nd.product_id as new_product_id
FROM customer c
JOIN network_devices nd ON c.id = nd.customer_id;
```

### Step 3: Update Application Code
- Modify backend APIs
- Update frontend components
- Test thoroughly

### Step 4: Remove Old Column (After Testing)
Uncomment the final DROP statements in the migration file and run again.

## Rollback Plan
If issues occur, you can rollback by:
1. Adding `product_id` back to `customer` table
2. Copying data from `network_devices.product_id` to `customer.product_id`
3. Removing `product_id` from `network_devices` table

## Testing Checklist
- [ ] Migration runs without errors
- [ ] Existing data is preserved
- [ ] New customers can be created without product_id
- [ ] Network devices can be created with product_id
- [ ] Customer detail views show multiple products
- [ ] Billing system works with multiple products
- [ ] All existing functionality still works
