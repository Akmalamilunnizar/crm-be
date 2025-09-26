# Example Request Payloads for Installation Report API

## 1. cURL Example (Multipart Form Data)

```bash
curl -X POST \
  http://localhost:8080/api/admin/customer-installation/report-installations \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -F "customer_id=24225552-b0c3-43d7-9a67-20e04f36fa5f" \
  -F "technician_id=c13a6c87-ec28-47ba-84c2-58b5ace2af57" \
  -F "assets_id=5ca1606b-66b3-4958-b7af-f48d4cda800a" \
  -F "status=pending" \
  -F "notes=Installation completed successfully" \
  -F "installation_type=new_installation" \
  -F "on_air_date=2024-01-15" \
  -F "trial_end_date=2024-02-15" \
  -F "service_ready_date=2024-01-15" \
  -F "installation_completed_at=2024-01-15 14:30:00" \
  -F "document_type=KTP" \
  -F "document_photo=@/path/to/ktp_document.jpg" \
  -F "switch_id=switch-1" \
  -F "port_number=1" \
  -F "remote_port=1" \
  -F "eth_port=1" \
  -F "mac_address=00:11:22:33:44:55" \
  -F "ip_static=192.168.1.100" \
  -F "kepemilikan_perangkat=owned" \
  -F "status_perangkat=active" \
  -F "last_ping_status=unknown" \
  -F "cable_type=UTP Cat6" \
  -F "cable_length=100.5" \
  -F "user_login=customer_login" \
  -F "password=customer_password123" \
  -F "user_status=Active" \
  -F "end_port_type=RJ45" \
  -F "installation_team_name=John Doe" \
  -F "installation_team_phone=08123456789"
```

## 2. JavaScript Fetch Example

```javascript
const formData = new FormData();

// Basic Information
formData.append('customer_id', '24225552-b0c3-43d7-9a67-20e04f36fa5f');
formData.append('technician_id', 'c13a6c87-ec28-47ba-84c2-58b5ace2af57');
formData.append('assets_id', '5ca1606b-66b3-4958-b7af-f48d4cda800a');
formData.append('status', 'pending');
formData.append('notes', 'Installation completed successfully');
formData.append('installation_type', 'new_installation');
formData.append('on_air_date', '2024-01-15');
formData.append('trial_end_date', '2024-02-15');
formData.append('service_ready_date', '2024-01-15');
formData.append('installation_completed_at', '2024-01-15 14:30:00');

// Document Information
formData.append('document_type', 'KTP');
formData.append('document_photo', documentPhotoFile); // File object

// Network Device Information
formData.append('switch_id', 'switch-1');
formData.append('port_number', '1');
formData.append('remote_port', '1');
formData.append('eth_port', '1');
formData.append('mac_address', '00:11:22:33:44:55');
formData.append('ip_static', '192.168.1.100');
formData.append('kepemilikan_perangkat', 'owned');
formData.append('status_perangkat', 'active');
formData.append('last_ping_status', 'unknown');

// Cable Information
formData.append('cable_type', 'UTP Cat6');
formData.append('cable_length', '100.5');

// Customer Service Information
formData.append('user_login', 'customer_login');
formData.append('password', 'customer_password123');
formData.append('user_status', 'Active');
formData.append('end_port_type', 'RJ45');

// Installation Team Information
formData.append('installation_team_name', 'John Doe');
formData.append('installation_team_phone', '08123456789');

// Submit request
fetch('/api/admin/customer-installation/report-installations', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`
  },
  body: formData
})
.then(response => response.json())
.then(data => console.log('Success:', data))
.catch(error => console.error('Error:', error));
```

## 3. Vue.js Component Example

```vue
<template>
  <form @submit.prevent="submitForm">
    <!-- Customer Selection -->
    <select v-model="formData.customer_id" required>
      <option value="">Select Customer</option>
      <option value="24225552-b0c3-43d7-9a67-20e04f36fa5f">Cha Eunwoo</option>
      <option value="55f03555-973d-425b-92bd-a0825c20b8ce">Rita Darmayanti</option>
    </select>

    <!-- Technician Selection -->
    <select v-model="formData.technician_id" required>
      <option value="">Select Technician</option>
      <option value="c13a6c87-ec28-47ba-84c2-58b5ace2af57">teknisi</option>
      <option value="ba73859d-8957-418a-8311-28a6b31cad95">cantiks</option>
    </select>

    <!-- Assets Selection -->
    <select v-model="formData.assets_id" required>
      <option value="">Select Asset</option>
      <option value="5ca1606b-66b3-4958-b7af-f48d4cda800a">TP-Link Router</option>
      <option value="842a48a4-6380-4280-96e1-a734f61a7d5b">Modem</option>
    </select>

    <!-- Document Photo Upload -->
    <input 
      type="file" 
      @change="handleFileUpload" 
      accept="image/*"
      ref="documentPhoto"
    />

    <!-- Network Configuration -->
    <input v-model="formData.switch_id" placeholder="Switch ID" />
    <input v-model="formData.port_number" placeholder="Port Number" />
    <input v-model="formData.mac_address" placeholder="MAC Address" />
    <input v-model="formData.ip_static" placeholder="IP Address" />

    <!-- Cable Information -->
    <input v-model="formData.cable_type" placeholder="Cable Type" />
    <input v-model.number="formData.cable_length" type="number" placeholder="Cable Length" />

    <!-- Customer Service -->
    <input v-model="formData.user_login" placeholder="User Login" />
    <input v-model="formData.password" type="password" placeholder="Password" />

    <button type="submit" :disabled="isSubmitting">
      {{ isSubmitting ? 'Creating...' : 'Create Installation Report' }}
    </button>
  </form>
</template>

<script setup>
import { ref, reactive } from 'vue'

const isSubmitting = ref(false)
const documentPhoto = ref(null)

const formData = reactive({
  customer_id: '',
  technician_id: '',
  assets_id: '',
  status: 'pending',
  notes: '',
  installation_type: 'new_installation',
  on_air_date: '',
  trial_end_date: '',
  service_ready_date: '',
  installation_completed_at: '',
  document_type: 'KTP',
  document_photo: null,
  switch_id: '',
  port_number: '',
  remote_port: '',
  eth_port: '',
  mac_address: '',
  ip_static: '',
  kepemilikan_perangkat: 'owned',
  status_perangkat: 'active',
  last_ping_status: 'unknown',
  cable_type: '',
  cable_length: 0,
  user_login: '',
  password: '',
  user_status: 'Active',
  end_port_type: '',
  installation_team_name: '',
  installation_team_phone: ''
})

function handleFileUpload(event) {
  const file = event.target.files[0]
  if (file) {
    formData.document_photo = file
  }
}

async function submitForm() {
  isSubmitting.value = true
  
  try {
    const formDataToSend = new FormData()
    
    // Add all form fields
    Object.keys(formData).forEach(key => {
      const value = formData[key]
      if (value !== null && value !== undefined && value !== '') {
        formDataToSend.append(key, value)
      }
    })
    
    const response = await fetch('/api/admin/customer-installation/report-installations', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${useCookie('token').value}`
      },
      body: formDataToSend
    })
    
    if (!response.ok) {
      throw new Error('Failed to create installation report')
    }
    
    const result = await response.json()
    console.log('Success:', result)
    
    // Redirect or show success message
    await navigateTo('/dashboard/report/customer-installation')
    
  } catch (error) {
    console.error('Error:', error)
    // Show error message
  } finally {
    isSubmitting.value = false
  }
}
</script>
```

## 4. Postman Collection

### Request Configuration:
- **Method**: POST
- **URL**: `{{base_url}}/api/admin/customer-installation/report-installations`
- **Headers**: 
  - `Authorization: Bearer {{token}}`
- **Body**: form-data

### Form Data Fields:
```
customer_id: 24225552-b0c3-43d7-9a67-20e04f36fa5f
technician_id: c13a6c87-ec28-47ba-84c2-58b5ace2af57
assets_id: 5ca1606b-66b3-4958-b7af-f48d4cda800a
status: pending
notes: Installation completed successfully
installation_type: new_installation
on_air_date: 2024-01-15
trial_end_date: 2024-02-15
service_ready_date: 2024-01-15
installation_completed_at: 2024-01-15 14:30:00
document_type: KTP
document_photo: [FILE] (select image file)
switch_id: switch-1
port_number: 1
remote_port: 1
eth_port: 1
mac_address: 00:11:22:33:44:55
ip_static: 192.168.1.100
kepemilikan_perangkat: owned
status_perangkat: active
last_ping_status: unknown
cable_type: UTP Cat6
cable_length: 100.5
user_login: customer_login
password: customer_password123
user_status: Active
end_port_type: RJ45
installation_team_name: John Doe
installation_team_phone: 08123456789
```

## 5. Expected Response

### Success Response (201 Created):
```json
{
  "success": true,
  "message": "Installation report created successfully",
  "data": {
    "id": "new-installation-uuid",
    "customer_id": "24225552-b0c3-43d7-9a67-20e04f36fa5f",
    "technician_id": "c13a6c87-ec28-47ba-84c2-58b5ace2af57",
    "status": "pending",
    "notes": "Installation completed successfully",
    "document_type": "KTP",
    "document_photo": "uploads/installations/documents/document_20240115_143000.jpg",
    "installation_type": "new_installation",
    "on_air_date": "2024-01-15T00:00:00Z",
    "trial_end_date": "2024-02-15T00:00:00Z",
    "service_ready_date": "2024-01-15T00:00:00Z",
    "installation_completed_at": "2024-01-15T14:30:00Z",
    "created_at": "2024-01-15T14:30:00Z",
    "updated_at": "2024-01-15T14:30:00Z",
    "customer": {
      "id": "24225552-b0c3-43d7-9a67-20e04f36fa5f",
      "name": "Cha Eunwoo",
      "address": "Jatirenggo, Turen, Talok, Kabupaten Malang",
      "phone": "081217706557"
    },
    "technician": {
      "id": "c13a6c87-ec28-47ba-84c2-58b5ace2af57",
      "name": "teknisi",
      "phone": "9764020363"
    },
    "network_devices": [
      {
        "id": "new-device-uuid",
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
        "id": "new-service-uuid",
        "user_login": "customer_login",
        "password": "encrypted_password",
        "user_status": "Active",
        "end_port_type": "RJ45"
      }
    ],
    "cables": [
      {
        "id": "new-cable-uuid",
        "name": "Installation Cable",
        "type": "UTP Cat6",
        "length": 100.5,
        "status": "in_use"
      }
    ]
  }
}
```

### Error Response (400 Bad Request):
```json
{
  "success": false,
  "message": "Invalid form data",
  "data": "Key: 'CreateReportInstallationRequest.CustomerID' Error:Field validation for 'CustomerID' failed on the 'required' tag"
}
```

## 6. Test Data

### Valid Test IDs:
- **Customer ID**: `24225552-b0c3-43d7-9a67-20e04f36fa5f` (Cha Eunwoo)
- **Technician ID**: `c13a6c87-ec28-47ba-84c2-58b5ace2af57` (teknisi)
- **Asset ID**: `5ca1606b-66b3-4958-b7af-f48d4cda800a` (TP-Link Router)

### Valid Test Values:
- **MAC Address**: `00:11:22:33:44:55`, `AA:BB:CC:DD:EE:FF`
- **IP Address**: `192.168.1.100`, `10.0.0.1`, `172.16.0.1`
- **Cable Types**: `UTP Cat6`, `UTP Cat5e`, `Single Mode Fiber`, `Multi Mode Fiber`
- **End Port Types**: `RJ45`, `Fiber`, `SC`, `LC`

### File Upload Test:
- **Valid**: JPG, PNG files under 5MB
- **Invalid**: PDF, DOC, files over 5MB, corrupted files
