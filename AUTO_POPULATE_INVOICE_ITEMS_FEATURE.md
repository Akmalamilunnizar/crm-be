# Auto-Populate Invoice Items from Network Devices Feature

## Overview
This feature automatically populates recurring invoice items based on a customer's network devices and their associated products. When you select a customer, the system will automatically fetch their network devices and create invoice items for each device that has a product assigned.

## How It Works

### 1. **Customer Selection**
- When you select a customer in the recurring invoice form, the system automatically fetches their network devices
- Only devices that have a `product_id` assigned will be included in the invoice items

### 2. **Automatic Population**
- Each network device with a product becomes an invoice item
- Item name: `{Product Name} - {MAC Address or 'Device'}`
- Quantity: 1 (can be modified)
- Price: Product price (can be modified)
- Total: Automatically calculated

### 3. **Manual Control**
- You can still manually add, edit, or remove items
- Use the "Auto-fill from Devices" button to refresh items from network devices
- Use the "Add Item" button to add custom items

## Prerequisites

### 1. **Database Migration**
Before using this feature, you must run the database migration to add `product_id` to the `network_devices` table:

```bash
# Run the migration
go run migrate-customer-multiple-products.go
```

### 2. **Network Device Setup**
Each network device must have a `product_id` assigned to be included in automatic invoice generation.

## Example Scenarios

### Scenario 1: Business Customer with Multiple Devices
```
Customer: "PT Maju Bersama"
├── Device 1: Office Router (MAC: 14:E2:O2:19:22:90)
│   ├── Product: "Business 100 Mbps" (800,000/month)
│   └── Invoice Item: "Business 100 Mbps - 14:E2:O2:19:22:90"
├── Device 2: Warehouse Router (MAC: 40:EE:15:C8:67:5D)
│   ├── Product: "Standard 50 Mbps" (500,000/month)
│   └── Invoice Item: "Standard 50 Mbps - 40:EE:15:C8:67:5D"
└── Total: 1,300,000/month
```

### Scenario 2: Residential Customer
```
Customer: "Cha Eunwoo"
├── Device 1: Main Router (MAC: 12:34:56:78:90:AB)
│   ├── Product: "Premium 100 Mbps" (100,000/month)
│   └── Invoice Item: "Premium 100 Mbps - 12:34:56:78:90:AB"
└── Total: 100,000/month
```

### Scenario 3: Customer with No Product-Assigned Devices
```
Customer: "New Customer"
├── No devices with products assigned
└── Default Item: "Internet Service" (Price: 0 - to be filled manually)
```

## User Interface

### 1. **Customer Selection**
- Select a customer from the dropdown
- System automatically fetches and populates invoice items

### 2. **Invoice Items Section**
- **Auto-fill from Devices** button: Manually refresh items from network devices
- **Add Item** button: Add custom invoice items
- **Loading indicator**: Shows when fetching network devices

### 3. **Item Display**
- Each item shows: Name, Quantity, Price, Total
- Items can be edited or removed
- Total amount updates automatically

## Technical Implementation

### Backend Changes
1. **Network Device Model**: Added `ProductID` field and `Product` relationship
2. **Repository**: Updated to preload product information
3. **API**: Existing endpoint `/api/admin/network-device/customer/{customerId}` returns product data

### Frontend Changes
1. **Form Component**: Added auto-population logic
2. **API Client**: Uses existing `getNetworkDevicesByCustomer` method
3. **Watchers**: Automatically triggers when customer is selected

## API Endpoints Used

### Get Network Devices by Customer
```
GET /api/admin/network-device/customer/{customerId}
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "device-id",
      "customer_id": "customer-id",
      "product_id": "product-id",
      "product": {
        "id": "product-id",
        "name": "Business 100 Mbps",
        "price": 800000
      },
      "mac_address": "14:E2:O2:19:22:90",
      "ip_static": "172.168.92.12"
    }
  ]
}
```

## Benefits

### 1. **Automation**
- Reduces manual data entry
- Ensures accuracy by using actual device data
- Saves time when creating recurring invoices

### 2. **Flexibility**
- Can still manually edit items
- Supports mixed scenarios (auto + manual items)
- Easy to refresh from device data

### 3. **Accuracy**
- Uses real product prices from the database
- Includes actual device information
- Prevents pricing errors

## Troubleshooting

### Issue: No items are auto-populated
**Solution:**
1. Check if the customer has network devices
2. Verify that network devices have `product_id` assigned
3. Run the database migration if not done yet

### Issue: Items show with price 0
**Solution:**
1. Check if the products have correct prices in the database
2. Verify the product relationship is working

### Issue: Auto-fill button not working
**Solution:**
1. Check browser console for errors
2. Verify API endpoint is accessible
3. Check authentication token

## Future Enhancements

### 1. **Bulk Operations**
- Auto-populate for multiple customers
- Batch invoice generation

### 2. **Advanced Filtering**
- Filter by device status
- Include/exclude specific device types

### 3. **Customization**
- Custom item naming templates
- Automatic quantity calculation based on usage

### 4. **Integration**
- Sync with billing systems
- Automatic invoice generation scheduling
