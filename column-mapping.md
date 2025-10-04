# MAPPING KOLOM INSTALLATION

## 📋 LOKASI KOLOM-KOLOM INSTALLATION

### ✅ **KOLOM YANG ADA DI DATABASE:**

| **Kolom** | **Tabel** | **Tipe Data** | **Keterangan** |
|-----------|-----------|---------------|----------------|
| `installation_date` | `customer` | `date` | Tanggal instalasi customer |
| `phone` | `customer` | `varchar(191)` | No. HP customer |
| `phone` | `users` | `varchar(191)` | No. HP technician (via technician_id) |
| `on_air_date` | `customer_installations` | `date` | Tanggal on-air |
| `date_trial` | `customer_installations` | `date` | Tanggal trial |
| `service_activation_date` | `customer_services` | `date` | Tanggal aktivasi layanan |
| `switch_id` | `network_devices` | `varchar(191)` | ID Switch |
| `port_number` | `network_devices` | `varchar(50)` | Nomor port |
| `type` | `assets` | `varchar(191)` | Tipe peralatan |
| `type` | `cable` | `varchar(100)` | Tipe kabel |
| `cable_length` | `customer_services` | `decimal(10,2)` | Panjang kabel |
| `end_port_type` | `customer_services` | `varchar(50)` | Tipe port akhir |
| `brand` | `assets` | `varchar(191)` | Merek peralatan |
| `model` | `assets` | `varchar(191)` | Model peralatan |
| `model` | `switch` | `varchar(255)` | Model switch |
| `mac_address` | `network_devices` | `varchar(191)` | MAC Address perangkat |
| `mac_address` | `assets` | `varchar(191)` | MAC Address aset |
| `eth_port` | `network_devices` | `varchar(50)` | Port ethernet |
| `ip_static` | `network_devices` | `varchar(191)` | IP Static |
| `remote_port` | `network_devices` | `varchar(50)` | Port remote |
| `user_login` | `customer_services` | `varchar(191)` | Login user |
| `password` | `customer` | `varchar(191)` | Password customer |
| `password` | `customer_services` | `varchar(191)` | Password layanan |
| `password` | `users` | `varchar(191)` | Password user |
| `status_perangkat` | `network_devices` | `enum` | Status perangkat |
| `kepemilikan_perangkat` | `network_devices` | `enum` | Kepemilikan perangkat |
| `last_ping_status` | `network_devices` | `enum` | Status ping terakhir |

### ❌ **KOLOM YANG TIDAK DITEMUKAN:**

| **Kolom** | **Status** | **Keterangan** |
|-----------|------------|----------------|
| `technician_id` | ❌ NOT FOUND | Sudah dihapus dari customer_installations |
| `status_user` | ❌ NOT FOUND | Tidak ada di database |

### 🔄 **KOLOM YANG BERADA DI BEBERAPA TABEL:**

| **Kolom** | **Tabel 1** | **Tabel 2** | **Keterangan** |
|-----------|-------------|-------------|----------------|
| `phone` | `customer` | `users` | Customer phone vs Technician phone |
| `type` | `assets` | `cable` | Asset type vs Cable type |
| `model` | `assets` | `switch` | Asset model vs Switch model |
| `mac_address` | `network_devices` | `assets` | Device MAC vs Asset MAC |
| `password` | `customer` | `customer_services` | Customer password vs Service password |

## 🎯 **REKOMENDASI UNTUK FORM INSTALLATION:**

### **Data yang Diambil dari Relasi:**
- **Technician Info** → `users` table (via technician_id)
- **Network Config** → `network_devices` table (via customer_id)
- **Service Config** → `customer_services` table (via customer_id)
- **Asset Info** → `assets` table (via assets_id di network_devices)
- **Cable Info** → `cable` table (via cable_id di customer_services)

### **Data yang Diinput Langsung:**
- `on_air_date` → `customer_installations.on_air_date`
- `date_trial` → `customer_installations.date_trial`
- `description` → `customer_installations.description`
- `equipment_used` → `customer_installations.equipment_used`
- `notes` → `customer_installations.notes`
- `completion_time` → `customer_installations.completion_time`
- `customer_signature` → `customer_installations.customer_signature`














