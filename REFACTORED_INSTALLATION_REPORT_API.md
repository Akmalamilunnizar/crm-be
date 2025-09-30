# Refactored Installation Report API Documentation

## Overview
This document describes the refactored Add Report Installation feature that has been updated to work with the new database structure and supports multipart form data for file uploads.

## API Endpoint

### POST `/api/admin/customer-installation/report-installations`

Creates a new installation report with all related data using multipart form data.

#### Headers
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

#### Form Fields

##### Basic Installation Information
- `customer_id` (required): Customer ID
- `technician_id` (required): Technician ID  
- `status`: Installation status (default: "pending")
- `notes`: Installation notes
- `installation_type`: Type of installation (default: "new_installation")
- `on_air_date`: Date when service went live (YYYY-MM-DD format)
- `trial_end_date`: Trial period end date (YYYY-MM-DD format)
- `service_ready_date`: Service ready date (YYYY-MM-DD format)
- `installation_completed_at`: Installation completion datetime (YYYY-MM-DD HH:MM:SS format)

##### Document Information
- `document_type`: Document type (KTP, SIM, Paspor)
- `document_photo`: Document photo file (JPG/PNG only, max 5MB)

##### Network Device Information
- `assets_id` (required): Asset/Router ID
- `switch_id`: Switch ID
- `port_number`: Port number
- `remote_port`: Remote port
- `eth_port`: Ethernet port
- `mac_address`: MAC address (format: XX:XX:XX:XX:XX:XX)
- `ip_static`: Static IP address (format: XXX.XXX.XXX.XXX)
- `kepemilikan_perangkat`: Device ownership (owned, leased, customer)
- `status_perangkat`: Device status (active, inactive, maintenance, faulty)
- `last_ping_status`: Last ping status (up, down, unknown)

##### Cable Information
- `cable_type`: Cable type (e.g., UTP Cat6, Fiber)
- `cable_length`: Cable length in meters (decimal)

##### Customer Service Information
- `user_login`: Customer login username
- `password`: Customer login password
- `user_status`: User status (Active, Inactive, Suspended, Pending)
- `end_port_type`: End port type (e.g., RJ45, Fiber)

##### Installation Team Information
- `installation_team_name`: Installation team member name
- `installation_team_phone`: Installation team phone number

#### Response

##### Success Response (201 Created)
```json
{
  "success": true,
  "message": "Installation report created successfully",
  "data": {
    "id": "uuid",
    "customer_id": "customer-uuid",
    "technician_id": "technician-uuid",
    "status": "pending",
    "notes": "Installation notes",
    "document_type": "KTP",
    "document_photo": "uploads/installations/documents/document_20240101_120000.jpg",
    "installation_type": "new_installation",
    "on_air_date": "2024-01-01T00:00:00Z",
    "trial_end_date": "2024-01-31T00:00:00Z",
    "service_ready_date": "2024-01-01T00:00:00Z",
    "installation_completed_at": "2024-01-01T12:00:00Z",
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z",
    "customer": {
      "id": "customer-uuid",
      "name": "Customer Name",
      "address": "Customer Address",
      "phone": "08123456789"
    },
    "technician": {
      "id": "technician-uuid",
      "name": "Technician Name",
      "phone": "08123456789"
    },
    "network_devices": [
      {
        "id": "device-uuid",
        "switch_id": "switch-1",
        "port_number": "1",
        "remote_port": "1",
        "eth_port": "1",
        "mac_address": "00:11:22:33:44:55",
        "ip_static": "192.168.1.100",
        "kepemilikan_perangkat": "owned",
        "status_perangkat": "active",
        "last_ping_status": "unknown"
      }
    ],
    "customer_services": [
      {
        "id": "service-uuid",
        "user_login": "customer_login",
        "password": "encrypted_password",
        "user_status": "Active",
        "end_port_type": "RJ45"
      }
    ],
    "cables": [
      {
        "id": "cable-uuid",
        "name": "Installation Cable",
        "type": "UTP Cat6",
        "length": 100.5,
        "status": "in_use"
      }
    ]
  }
}
```

##### Error Response (400 Bad Request)
```json
{
  "success": false,
  "message": "Invalid form data",
  "data": "validation error details"
}
```

##### Error Response (500 Internal Server Error)
```json
{
  "success": false,
  "message": "Failed to create installation report",
  "data": "error details"
}
```

## Frontend Implementation

### Vue Component Usage

The refactored Vue component is located at:
`crm-fe/pages/dashboard/report/customer-installation/AddReportInstallationRefactored.vue`

#### Key Features:
1. **Form Validation**: Uses Yup schema validation
2. **File Upload**: Supports document photo upload with validation
3. **Multipart Form Data**: Uses FormData for submission
4. **Real-time Validation**: Client-side validation for IP and MAC addresses
5. **Responsive Design**: Mobile-friendly form layout

#### Usage Example:
```vue
<template>
  <AddReportInstallationRefactored />
</template>

<script setup>
import AddReportInstallationRefactored from './AddReportInstallationRefactored.vue'
</script>
```

## Database Schema

### Tables Involved:
1. `customer_installations` - Main installation record
2. `network_devices` - Network device configuration
3. `customer_services` - Customer service information
4. `cable` - Cable information
5. `assets` - Asset/Router information
6. `customer` - Customer information
7. `users` - Technician information

### Key Relationships:
- `customer_installations.customer_id` → `customer.id`
- `customer_installations.technician_id` → `users.id`
- `network_devices.customer_installation_id` → `customer_installations.id`
- `network_devices.assets_id` → `assets.id`
- `customer_services.customer_installation_id` → `customer_installations.id`
- `cable.customer_installation_id` → `customer_installations.id`

## Validation Rules

### File Upload:
- **Document Photo**: JPG/PNG only, max 5MB
- **File Naming**: Auto-generated with timestamp

### Data Validation:
- **IP Address**: Basic format validation (XXX.XXX.XXX.XXX)
- **MAC Address**: Basic format validation (XX:XX:XX:XX:XX:XX)
- **Required Fields**: customer_id, technician_id, assets_id
- **Date Formats**: YYYY-MM-DD for dates, YYYY-MM-DD HH:MM:SS for datetime

## Migration Notes

### From Legacy API:
1. **Endpoint Change**: `/report/complete` → `/report-installations`
2. **Content Type**: `application/json` → `multipart/form-data`
3. **File Handling**: Base64 encoding → Direct file upload
4. **Field Mapping**: Updated to match new database structure

### Backward Compatibility:
- Legacy endpoint `/report/complete` still available
- Old request format still supported
- Gradual migration recommended

## Testing

### cURL Example:
```bash
curl -X POST \
  http://localhost:8080/api/admin/customer-installation/report-installations \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "customer_id=customer-uuid" \
  -F "technician_id=technician-uuid" \
  -F "assets_id=asset-uuid" \
  -F "status=pending" \
  -F "installation_type=new_installation" \
  -F "document_type=KTP" \
  -F "document_photo=@/path/to/document.jpg" \
  -F "switch_id=switch-1" \
  -F "port_number=1" \
  -F "mac_address=00:11:22:33:44:55" \
  -F "ip_static=192.168.1.100" \
  -F "cable_type=UTP Cat6" \
  -F "cable_length=100.5" \
  -F "user_login=customer_login" \
  -F "password=customer_password" \
  -F "user_status=Active"
```

### Postman Collection:
Import the provided Postman collection for easy testing of all endpoints.

## Error Handling

### Common Errors:
1. **400 Bad Request**: Invalid form data or validation errors
2. **401 Unauthorized**: Missing or invalid token
3. **403 Forbidden**: Insufficient permissions
4. **500 Internal Server Error**: Database or server errors

### Error Response Format:
All errors follow the same JSON structure with `success: false`, `message`, and optional `data` fields.
