# 🏗️ **STRUKTUR DATABASE UNTUK FITUR ADD REPORT INSTALLATION**

## 🎯 **Overview**

Struktur database telah disesuaikan untuk fitur "Add Report Installation" dengan membuat relasi yang tepat antara `customer_installations` dengan tabel-tabel lain yang berisi data instalasi. Tidak semua kolom dipindahkan ke `customer_installations`, melainkan dibuat relasi yang logis sesuai dengan kebutuhan.

---

## 🔗 **RELASI DATABASE YANG DIBUAT**

### **1. Relasi `customer_installations` → `network_devices`**
```sql
-- Tambahkan kolom customer_installation_id ke network_devices
ALTER TABLE network_devices 
ADD COLUMN `customer_installation_id` varchar(191) DEFAULT NULL 
COMMENT 'Reference to customer_installations table';

-- Foreign Key Constraint
ALTER TABLE network_devices 
ADD CONSTRAINT `fk_network_devices_customer_installation` 
FOREIGN KEY (`customer_installation_id`) REFERENCES `customer_installations` (`id`) 
ON DELETE SET NULL ON UPDATE CASCADE;
```

### **2. Relasi `customer_installations` → `customer_services`**
```sql
-- Tambahkan kolom customer_installation_id ke customer_services
ALTER TABLE customer_services 
ADD COLUMN `customer_installation_id` varchar(191) DEFAULT NULL 
COMMENT 'Reference to customer_installations table';

-- Foreign Key Constraint
ALTER TABLE customer_services 
ADD CONSTRAINT `fk_customer_services_customer_installation` 
FOREIGN KEY (`customer_installation_id`) REFERENCES `customer_installations` (`id`) 
ON DELETE SET NULL ON UPDATE CASCADE;
```

### **3. Relasi `customer_installations` → `cable`**
```sql
-- Tambahkan kolom customer_installation_id ke cable
ALTER TABLE cable 
ADD COLUMN `customer_installation_id` varchar(191) DEFAULT NULL 
COMMENT 'Reference to customer_installations table';

-- Foreign Key Constraint
ALTER TABLE cable 
ADD CONSTRAINT `fk_cable_customer_installation` 
FOREIGN KEY (`customer_installation_id`) REFERENCES `customer_installations` (`id`) 
ON DELETE SET NULL ON UPDATE CASCADE;
```

---

## 📊 **MAPPING KOLOM UNTUK FITUR ADD REPORT INSTALLATION**

| **No** | **Istilah** | **Lokasi Database** | **Tabel** | **Relasi** |
|--------|-------------|-------------------|-----------|------------|
| 1 | **Tgl. Pemasangan** | `customer.installation_date` | `customer` | ✅ Direct |
| 2 | **Nama Anggota Tim Install** | `users.name` | `users` | ✅ Via `customer_installations.technician_id` |
| 3 | **No.HP Tim Install** | `users.phone` | `users` | ✅ Via `customer_installations.technician_id` |
| 4 | **Tgl. On Air** | `customer_installations.on_air_date` | `customer_installations` | ✅ Direct |
| 5 | **Tgl. Batas Percobaan** | `customer_installations.trial_end_date` | `customer_installations` | ✅ Direct |
| 6 | **Tgl. Siap Layanan** | `customer_installations.service_ready_date` | `customer_installations` | ✅ Direct |
| 7 | **Switch ID** | `network_devices.switch_id` | `network_devices` | ✅ Via `customer_installation_id` |
| 8 | **Port Number** | `network_devices.port_number` | `network_devices` | ✅ Via `customer_installation_id` |
| 9 | **Cable Type** | `cable.type` | `cable` | ✅ Via `customer_installation_id` |
| 10 | **Length** | `customer_services.cable_length` | `customer_services` | ✅ Via `customer_installation_id` |
| 11 | **End Port Type** | `customer_services.end_port_type` | `customer_services` | ✅ Via `customer_installation_id` |
| 12 | **Router Brand** | `assets.brand` | `assets` | ✅ Via `network_devices.assets_id` |
| 13 | **Type/Series** | `assets.type` + `assets.model` | `assets` | ✅ Via `network_devices.assets_id` |
| 14 | **MacAddr** | `network_devices.mac_address` | `network_devices` | ✅ Via `customer_installation_id` |
| 15 | **EthPort** | `network_devices.eth_port` | `network_devices` | ✅ Via `customer_installation_id` |
| 16 | **IPAddr** | `network_devices.ip_static` | `network_devices` | ✅ Via `customer_installation_id` |
| 17 | **RemotePort** | `network_devices.remote_port` | `network_devices` | ✅ Via `customer_installation_id` |
| 18 | **UserLogin** | `customer_services.user_login` | `customer_services` | ✅ Via `customer_installation_id` |
| 19 | **Password** | `customer_services.password` | `customer_services` | ✅ Via `customer_installation_id` |
| 20 | **Status User** | `customer_services.user_status` | `customer_services` | ✅ Via `customer_installation_id` |
| 21 | **Status Perangkat** | `network_devices.status_perangkat` | `network_devices` | ✅ Via `customer_installation_id` |
| 22 | **Kepemilikan Perangkat** | `network_devices.kepemilikan_perangkat` | `network_devices` | ✅ Via `customer_installation_id` |
| 23 | **Ping** | `network_devices.last_ping_status` | `network_devices` | ✅ Via `customer_installation_id` |

---

## 🗄️ **STRUKTUR TABEL YANG DISESUAIKAN**

### **1. Tabel `customer_installations` (Central Hub)**
```sql
+---------------------------+--------------------------------------------------------------+------+-----+----------------------+
| Field                     | Type                                                         | Null | Key | Default              |
+---------------------------+--------------------------------------------------------------+------+-----+----------------------+
| id                        | varchar(191)                                                 | NO   | PRI | NULL                 | ✅ Primary Key
| customer_id               | varchar(191)                                                 | YES  | MUL | NULL                 | ✅ Foreign Key
| technician_id             | varchar(191)                                                 | YES  | MUL | NULL                 | ✅ Foreign Key
| status                    | varchar(191)                                                 | YES  | MUL | pending              | ✅ Status
| notes                     | text                                                         | YES  |     | NULL                 | ✅ Notes
| document_type             | enum('KTP','SIM','Paspor')                                  | YES  | MUL | NULL                 | ✅ Document Type
| document_photo            | varchar(255)                                                 | YES  |     | NULL                 | ✅ Document Photo
| installation_type         | enum('new_installation','maintenance','upgrade','downgrade') | YES  | MUL | new_installation     | ✅ Installation Type
| total_assets_out          | int                                                          | YES  |     | 0                    | ✅ Total Assets Out
| total_assets_in           | int                                                          | YES  |     | 0                    | ✅ Total Assets In
| installation_completed_at | datetime(3)                                                  | YES  | MUL | NULL                 | ✅ Completion Date
| trial_end_date            | date                                                         | YES  | MUL | NULL                 | ✅ Trial End Date
| service_ready_date        | date                                                         | YES  | MUL | NULL                 | ✅ Service Ready Date
| on_air_date               | date                                                         | YES  | MUL | NULL                 | ✅ On Air Date
| createdAt                 | datetime(3)                                                  | NO   |     | CURRENT_TIMESTAMP(3) | ✅ Created At
| updatedAt                 | datetime(3)                                                  | NO   |     | CURRENT_TIMESTAMP(3) | ✅ Updated At
+---------------------------+--------------------------------------------------------------+------+-----+----------------------+
```

### **2. Tabel `network_devices` (Updated)**
```sql
+---------------------------+--------------------------------------------------+------+-----+-------------------+-----------------------------------------------+
| Field                     | Type                                             | Null | Key | Default           | Extra                                         |
+---------------------------+--------------------------------------------------+------+-----+-------------------+-----------------------------------------------+
| id                        | varchar(191)                                     | NO   | PRI | NULL              | ✅ Primary Key
| customer_id               | varchar(191)                                     | NO   | MUL | NULL              | ✅ Foreign Key
| assets_id                 | varchar(191)                                     | NO   | MUL | NULL              | ✅ Foreign Key
| customer_installation_id  | varchar(191)                                     | YES  | MUL | NULL              | 🆕 New Foreign Key
| switch_id                 | varchar(191)                                     | YES  |     | NULL              | ✅ Switch ID
| port_number               | varchar(50)                                      | YES  |     | NULL              | ✅ Port Number
| remote_port               | varchar(50)                                      | YES  |     | NULL              | ✅ Remote Port
| eth_port                  | varchar(50)                                      | YES  |     | NULL              | ✅ Eth Port
| kepemilikan_perangkat     | enum('owned','leased','customer')                | YES  |     | owned             | ✅ Kepemilikan Perangkat
| status_perangkat          | enum('active','inactive','maintenance','faulty') | YES  |     | active            | ✅ Status Perangkat
| last_ping_status          | enum('up','down','unknown')                      | YES  |     | unknown           | ✅ Ping Status
| last_ping_timestamp       | timestamp                                        | YES  |     | NULL              | ✅ Ping Timestamp
| mac_address               | varchar(191)                                     | YES  |     | NULL              | ✅ Mac Address
| ip_static                 | varchar(191)                                     | YES  |     | NULL              | ✅ IP Static
| product_id                | varchar(191)                                     | YES  | MUL | NULL              | ✅ Product ID
| created_at                | timestamp                                        | YES  |     | CURRENT_TIMESTAMP | ✅ Created At
| updated_at                | timestamp                                        | YES  |     | CURRENT_TIMESTAMP | ✅ Updated At
+---------------------------+--------------------------------------------------+------+-----+-------------------+-----------------------------------------------+
```

### **3. Tabel `customer_services` (Updated)**
```sql
+---------------------------+-------------------------------------------------+------+-----+-------------------+-----------------------------------------------+
| Field                     | Type                                            | Null | Key | Default           | Extra                                         |
+---------------------------+-------------------------------------------------+------+-----+-------------------+-----------------------------------------------+
| id                        | varchar(191)                                    | NO   | PRI | NULL              | ✅ Primary Key
| customer_id               | varchar(191)                                    | NO   | MUL | NULL              | ✅ Foreign Key
| customer_installation_id  | varchar(191)                                    | YES  | MUL | NULL              | 🆕 New Foreign Key
| device_id                 | varchar(191)                                    | YES  | MUL | NULL              | ✅ Device ID
| cable_id                  | varchar(191)                                    | YES  |     | NULL              | ✅ Cable ID
| cable_length              | decimal(10,2)                                   | YES  |     | NULL              | ✅ Cable Length
| end_port_type             | varchar(50)                                     | YES  |     | NULL              | ✅ End Port Type
| user_login                | varchar(191)                                    | YES  |     | NULL              | ✅ User Login
| service_activation_date   | date                                            | YES  |     | NULL              | ✅ Service Activation Date
| password                  | varchar(191)                                    | YES  |     | NULL              | ✅ Password
| user_status               | enum('Active','Inactive','Suspended','Pending') | YES  | MUL | Active            | ✅ User Status
| installation_notes        | text                                            | YES  |     | NULL              | ✅ Installation Notes
| installation_team_phone   | varchar(20)                                     | YES  | MUL | NULL              | ✅ Installation Team Phone
| created_at                | timestamp                                       | YES  |     | CURRENT_TIMESTAMP | ✅ Created At
| updated_at                | timestamp                                       | YES  |     | CURRENT_TIMESTAMP | ✅ Updated At
+---------------------------+-------------------------------------------------+------+-----+-------------------+-----------------------------------------------+
```

### **4. Tabel `cable` (Updated)**
```sql
+---------------------------+------------------------------------------------+------+-----+-------------------+-----------------------------------------------+
| Field                     | Type                                           | Null | Key | Default           | Extra                                         |
+---------------------------+------------------------------------------------+------+-----+-------------------+-----------------------------------------------+
| id                        | varchar(191)                                   | NO   | PRI | NULL              | ✅ Primary Key
| name                      | varchar(255)                                   | NO   |     | NULL              | ✅ Name
| type                      | varchar(100)                                   | YES  |     | NULL              | ✅ Type
| length                    | decimal(10,2)                                  | YES  |     | NULL              | ✅ Length
| status                    | enum('available','in_use','damaged','retired') | YES  |     | available         | ✅ Status
| customer_installation_id  | varchar(191)                                   | YES  | MUL | NULL              | 🆕 New Foreign Key
| created_at                | timestamp                                      | YES  |     | CURRENT_TIMESTAMP | ✅ Created At
| updated_at                | timestamp                                      | YES  |     | CURRENT_TIMESTAMP | ✅ Updated At
+---------------------------+------------------------------------------------+------+-----+-------------------+-----------------------------------------------+
```

---

## 🔧 **ENTITY MODELS YANG DISESUAIKAN**

### **1. CustomerInstallation (Updated)**
```go
type CustomerInstallation struct {
    ID                      string              `gorm:"column:id;type:varchar;primaryKey" json:"id"`
    CustomerID              *string             `gorm:"column:customer_id;type:varchar;index:idx_customer_installations_customer_id" json:"customer_id,omitempty"`
    TechnicianID            *string             `gorm:"column:technician_id;type:varchar;index:idx_customer_installations_technician_id" json:"technician_id,omitempty"`
    Status                  string              `gorm:"column:status;type:varchar;default:'pending'" json:"status,omitempty"`
    Notes                   string              `gorm:"column:notes;type:text" json:"notes,omitempty"`
    DocumentType            *string             `gorm:"column:document_type;type:enum('KTP','SIM','Paspor')" json:"document_type,omitempty"`
    DocumentPhoto           *string             `gorm:"column:document_photo;type:varchar(255)" json:"document_photo,omitempty"`
    InstallationType        string              `gorm:"column:installation_type;type:enum('new_installation','maintenance','upgrade','downgrade');default:'new_installation'" json:"installation_type,omitempty"`
    TotalAssetsOut          int                 `gorm:"column:total_assets_out;type:int;default:0" json:"total_assets_out,omitempty"`
    TotalAssetsIn           int                 `gorm:"column:total_assets_in;type:int;default:0" json:"total_assets_in,omitempty"`
    InstallationCompletedAt *time.Time          `gorm:"column:installation_completed_at;type:datetime(3)" json:"installation_completed_at,omitempty"`
    TrialEndDate            *time.Time          `gorm:"column:trial_end_date;type:date" json:"trial_end_date,omitempty"`
    ServiceReadyDate        *time.Time          `gorm:"column:service_ready_date;type:date" json:"service_ready_date,omitempty"`
    OnAirDate               *time.Time          `gorm:"column:on_air_date;type:date" json:"on_air_date,omitempty"`
    CreatedAt               time.Time           `gorm:"column:createdAt;default:current_timestamp" json:"createdAt"`
    UpdatedAt               time.Time           `gorm:"column:updatedAt;default:current_timestamp on update current_timestamp" json:"updatedAt"`
    Customer                *Customer           `gorm:"foreignKey:CustomerID;references:id;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"customer,omitempty"`
    Technician              *User               `gorm:"foreignKey:TechnicianID;references:id;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"technician,omitempty"`
    Images                  []Image             `gorm:"foreignKey:ArchiveInstallationId;references:ID" json:"images,omitempty"`
    AssetTransactions       []AssetTransaction  `gorm:"foreignKey:CustomerInstallationID;references:ID" json:"asset_transactions,omitempty"`
    NetworkDevices          []NetworkDevice     `gorm:"foreignKey:CustomerInstallationID;references:ID" json:"network_devices,omitempty"`
    CustomerServices        []CustomerService   `gorm:"foreignKey:CustomerInstallationID;references:ID" json:"customer_services,omitempty"`
    Cables                  []Cable             `gorm:"foreignKey:CustomerInstallationID;references:ID" json:"cables,omitempty"`
}
```

### **2. NetworkDevice (New)**
```go
type NetworkDevice struct {
    ID                      string             `gorm:"column:id;type:varchar;primaryKey" json:"id"`
    CustomerID              string             `gorm:"column:customer_id;type:varchar;index:idx_network_devices_customer_id" json:"customer_id"`
    AssetsID                string             `gorm:"column:assets_id;type:varchar;index:idx_network_devices_assets_id" json:"assets_id"`
    CustomerInstallationID  *string            `gorm:"column:customer_installation_id;type:varchar;index:idx_network_devices_customer_installation_id" json:"customer_installation_id,omitempty"`
    SwitchID                *string            `gorm:"column:switch_id;type:varchar" json:"switch_id,omitempty"`
    PortNumber              *string            `gorm:"column:port_number;type:varchar(50)" json:"port_number,omitempty"`
    RemotePort              *string            `gorm:"column:remote_port;type:varchar(50)" json:"remote_port,omitempty"`
    EthPort                 *string            `gorm:"column:eth_port;type:varchar(50)" json:"eth_port,omitempty"`
    KepemilikanPerangkat    string             `gorm:"column:kepemilikan_perangkat;type:enum('owned','leased','customer');default:'owned'" json:"kepemilikan_perangkat,omitempty"`
    StatusPerangkat         string             `gorm:"column:status_perangkat;type:enum('active','inactive','maintenance','faulty');default:'active'" json:"status_perangkat,omitempty"`
    LastPingStatus          string             `gorm:"column:last_ping_status;type:enum('up','down','unknown');default:'unknown'" json:"last_ping_status,omitempty"`
    LastPingTimestamp       *time.Time         `gorm:"column:last_ping_timestamp;type:timestamp" json:"last_ping_timestamp,omitempty"`
    MacAddress              *string            `gorm:"column:mac_address;type:varchar(191)" json:"mac_address,omitempty"`
    IPStatic                *string            `gorm:"column:ip_static;type:varchar(191)" json:"ip_static,omitempty"`
    ProductID               *string            `gorm:"column:product_id;type:varchar;index:idx_network_devices_product_id" json:"product_id,omitempty"`
    CreatedAt               *time.Time         `gorm:"column:created_at;type:timestamp;default:current_timestamp" json:"created_at,omitempty"`
    UpdatedAt               *time.Time         `gorm:"column:updated_at;type:timestamp;default:current_timestamp on update current_timestamp" json:"updated_at,omitempty"`
    Customer                *Customer          `gorm:"foreignKey:CustomerID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"customer,omitempty"`
    Asset                   *Asset             `gorm:"foreignKey:AssetsID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"asset,omitempty"`
    CustomerInstallation    *CustomerInstallation `gorm:"foreignKey:CustomerInstallationID;references:ID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"customer_installation,omitempty"`
    Product                 *Product           `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"product,omitempty"`
}
```

### **3. CustomerService (New)**
```go
type CustomerService struct {
    ID                      string             `gorm:"column:id;type:varchar;primaryKey" json:"id"`
    CustomerID              string             `gorm:"column:customer_id;type:varchar;index:idx_customer_services_customer_id" json:"customer_id"`
    CustomerInstallationID  *string            `gorm:"column:customer_installation_id;type:varchar;index:idx_customer_services_customer_installation_id" json:"customer_installation_id,omitempty"`
    DeviceID                *string            `gorm:"column:device_id;type:varchar;index:idx_customer_services_device_id" json:"device_id,omitempty"`
    CableID                 *string            `gorm:"column:cable_id;type:varchar" json:"cable_id,omitempty"`
    CableLength             *float64           `gorm:"column:cable_length;type:decimal(10,2)" json:"cable_length,omitempty"`
    EndPortType             *string            `gorm:"column:end_port_type;type:varchar(50)" json:"end_port_type,omitempty"`
    UserLogin               *string            `gorm:"column:user_login;type:varchar(191)" json:"user_login,omitempty"`
    ServiceActivationDate   *time.Time         `gorm:"column:service_activation_date;type:date" json:"service_activation_date,omitempty"`
    Password                *string            `gorm:"column:password;type:varchar(191)" json:"password,omitempty"`
    UserStatus              string             `gorm:"column:user_status;type:enum('Active','Inactive','Suspended','Pending');default:'Active'" json:"user_status,omitempty"`
    InstallationNotes       *string            `gorm:"column:installation_notes;type:text" json:"installation_notes,omitempty"`
    InstallationTeamPhone   *string            `gorm:"column:installation_team_phone;type:varchar(20)" json:"installation_team_phone,omitempty"`
    CreatedAt               *time.Time         `gorm:"column:created_at;type:timestamp;default:current_timestamp" json:"created_at,omitempty"`
    UpdatedAt               *time.Time         `gorm:"column:updated_at;type:timestamp;default:current_timestamp on update current_timestamp" json:"updated_at,omitempty"`
    Customer                *Customer          `gorm:"foreignKey:CustomerID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"customer,omitempty"`
    CustomerInstallation    *CustomerInstallation `gorm:"foreignKey:CustomerInstallationID;references:ID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"customer_installation,omitempty"`
    NetworkDevice           *NetworkDevice     `gorm:"foreignKey:DeviceID;references:ID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"network_device,omitempty"`
}
```

### **4. Cable (New)**
```go
type Cable struct {
    ID                      string             `gorm:"column:id;type:varchar;primaryKey" json:"id"`
    Name                    string             `gorm:"column:name;type:varchar(255)" json:"name"`
    Type                    *string            `gorm:"column:type;type:varchar(100)" json:"type,omitempty"`
    Length                  *float64           `gorm:"column:length;type:decimal(10,2)" json:"length,omitempty"`
    Status                  string             `gorm:"column:status;type:enum('available','in_use','damaged','retired');default:'available'" json:"status,omitempty"`
    CustomerInstallationID  *string            `gorm:"column:customer_installation_id;type:varchar;index:idx_cable_customer_installation_id" json:"customer_installation_id,omitempty"`
    CreatedAt               *time.Time         `gorm:"column:created_at;type:timestamp;default:current_timestamp" json:"created_at,omitempty"`
    UpdatedAt               *time.Time         `gorm:"column:updated_at;type:timestamp;default:current_timestamp on update current_timestamp" json:"updated_at,omitempty"`
    CustomerInstallation    *CustomerInstallation `gorm:"foreignKey:CustomerInstallationID;references:ID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"customer_installation,omitempty"`
}
```

---

## 📊 **VIEWS UNTUK LAPORAN INSTALASI**

### **1. View `installation_report_complete`**
```sql
CREATE OR REPLACE VIEW `installation_report_complete` AS
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
    cs.installation_team_phone,
    
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

### **2. View `installation_summary_per_customer`**
```sql
CREATE OR REPLACE VIEW `installation_summary_per_customer` AS
SELECT 
    c.id as customer_id,
    c.name as customer_name,
    c.address as customer_address,
    c.phone as customer_phone,
    COUNT(ci.id) as total_installations,
    COUNT(CASE WHEN ci.status = 'completed' THEN 1 END) as completed_installations,
    COUNT(CASE WHEN ci.status = 'pending' THEN 1 END) as pending_installations,
    COUNT(CASE WHEN ci.status = 'in_progress' THEN 1 END) as in_progress_installations,
    MAX(ci.on_air_date) as latest_on_air_date,
    MAX(ci.installation_completed_at) as latest_completion_date
FROM customer c
LEFT JOIN customer_installations ci ON c.id = ci.customer_id
GROUP BY c.id, c.name, c.address, c.phone
ORDER BY total_installations DESC;
```

### **3. View `installation_asset_report`**
```sql
CREATE OR REPLACE VIEW `installation_asset_report` AS
SELECT 
    ci.id as installation_id,
    c.name as customer_name,
    ci.installation_type,
    ci.status as installation_status,
    ci.on_air_date,
    ci.installation_completed_at,
    
    -- Asset Out Summary
    COUNT(CASE WHEN at.transaction_type = 'out' THEN 1 END) as total_assets_out,
    SUM(CASE WHEN at.transaction_type = 'out' THEN at.quantity ELSE 0 END) as total_quantity_out,
    
    -- Asset In Summary
    COUNT(CASE WHEN at.transaction_type = 'in' THEN 1 END) as total_assets_in,
    SUM(CASE WHEN at.transaction_type = 'in' THEN at.quantity ELSE 0 END) as total_quantity_in,
    
    -- Asset Details
    GROUP_CONCAT(DISTINCT 
        CASE WHEN at.transaction_type = 'out' THEN 
            CONCAT(a.brand, ' ', a.model, ' (', at.quantity, ' pcs)')
        END 
        SEPARATOR ', '
    ) as assets_out_details,
    
    GROUP_CONCAT(DISTINCT 
        CASE WHEN at.transaction_type = 'in' THEN 
            CONCAT(a.brand, ' ', a.model, ' (', at.quantity, ' pcs)')
        END 
        SEPARATOR ', '
    ) as assets_in_details
    
FROM customer_installations ci
LEFT JOIN customer c ON ci.customer_id = c.id
LEFT JOIN asset_transactions at ON ci.id = at.customer_installation_id
LEFT JOIN assets a ON at.asset_id = a.id
GROUP BY ci.id, c.name, ci.installation_type, ci.status, ci.on_air_date, ci.installation_completed_at
ORDER BY ci.createdAt DESC;
```

### **4. View `installation_technician_report`**
```sql
CREATE OR REPLACE VIEW `installation_technician_report` AS
SELECT 
    u.id as technician_id,
    u.name as technician_name,
    u.phone as technician_phone,
    COUNT(ci.id) as total_installations,
    COUNT(CASE WHEN ci.status = 'completed' THEN 1 END) as completed_installations,
    COUNT(CASE WHEN ci.status = 'pending' THEN 1 END) as pending_installations,
    COUNT(CASE WHEN ci.status = 'in_progress' THEN 1 END) as in_progress_installations,
    AVG(DATEDIFF(ci.installation_completed_at, ci.createdAt)) as avg_completion_days,
    MAX(ci.installation_completed_at) as latest_completion_date
FROM users u
LEFT JOIN customer_installations ci ON u.id = ci.technician_id
WHERE u.role_id = (SELECT id FROM roles WHERE name = 'TECHNICIAN')
GROUP BY u.id, u.name, u.phone
ORDER BY total_installations DESC;
```

---

## 🌐 **API USAGE EXAMPLES**

### **1. Get Complete Installation Report**
```json
GET /api/installation-report/complete
Response:
{
    "data": [
        {
            "installation_id": "949189b8-97d1-11f0-9fa2-d843ae0f1e06",
            "customer_name": "Cha Eunwoo",
            "technician_name": "teknisi",
            "technician_phone": "9764020363",
            "installation_status": "completed",
            "installation_type": "new_installation",
            "on_air_date": "2025-09-17",
            "trial_end_date": "2025-10-17",
            "service_ready_date": "2025-09-17",
            "network_device_id": "sadvwabafbsdgd",
            "switch_id": null,
            "port_number": null,
            "mac_address": "14:E2:O2:19:22:90",
            "ip_static": "172.168.92.12",
            "router_brand": "ads",
            "router_type": "asdzxc",
            "router_model": "asd",
            "customer_service_id": "uuid-here",
            "user_login": "cha_eunwoo",
            "user_status": "Active",
            "cable_id": "cable-installation-1",
            "cable_type": "UTP Cat6",
            "cable_length": 100.00,
            "end_port_type": "RJ45"
        }
    ]
}
```

### **2. Get Installation Summary per Customer**
```json
GET /api/installation-report/summary-per-customer
Response:
{
    "data": [
        {
            "customer_id": "24225552-b0c3-43d7-9a67-20e04f36fa5f",
            "customer_name": "Cha Eunwoo",
            "total_installations": 1,
            "completed_installations": 1,
            "pending_installations": 0,
            "in_progress_installations": 0,
            "latest_on_air_date": "2025-09-17",
            "latest_completion_date": "2025-09-23T10:49:22.000Z"
        }
    ]
}
```

### **3. Get Installation Asset Report**
```json
GET /api/installation-report/asset-report
Response:
{
    "data": [
        {
            "installation_id": "949189b8-97d1-11f0-9fa2-d843ae0f1e06",
            "customer_name": "Cha Eunwoo",
            "installation_type": "new_installation",
            "installation_status": "completed",
            "total_assets_out": 2,
            "total_quantity_out": 2,
            "total_assets_in": 0,
            "total_quantity_in": 0,
            "assets_out_details": "TP-Link Router Edit (1 pcs), ads asdzxc (1 pcs)",
            "assets_in_details": null
        }
    ]
}
```

---

## 🔍 **QUERY EXAMPLES**

### **1. Get Installation Report dengan Semua Data Terkait**
```sql
SELECT * FROM installation_report_complete 
WHERE installation_id = '949189b8-97d1-11f0-9fa2-d843ae0f1e06';
```

### **2. Get Installation Summary per Customer**
```sql
SELECT * FROM installation_summary_per_customer 
WHERE customer_id = '24225552-b0c3-43d7-9a67-20e04f36fa5f';
```

### **3. Get Asset Report per Installation**
```sql
SELECT * FROM installation_asset_report 
WHERE installation_id = '949189b8-97d1-11f0-9fa2-d843ae0f1e06';
```

### **4. Get Technician Performance Report**
```sql
SELECT * FROM installation_technician_report 
WHERE technician_id = 'c13a6c87-ec28-47ba-84c2-58b5ace2af57';
```

---

## ✅ **VERIFIKASI SISTEM**

- ✅ Relasi `customer_installations` → `network_devices` berhasil dibuat
- ✅ Relasi `customer_installations` → `customer_services` berhasil dibuat
- ✅ Relasi `customer_installations` → `cable` berhasil dibuat
- ✅ Foreign key constraints berfungsi dengan baik
- ✅ Indexes untuk performa query optimal
- ✅ Views untuk laporan berhasil dibuat
- ✅ Sample data berhasil dimasukkan
- ✅ Entity models sesuai dengan database schema
- ✅ Relasi antar tabel berfungsi dengan baik

---

## 🚀 **KESIMPULAN**

**Struktur database telah disesuaikan dengan sempurna untuk fitur "Add Report Installation"!**

### **Keuntungan Struktur Baru:**
- **🎯 Central Hub**: `customer_installations` sebagai pusat data instalasi
- **🔗 Relasi Logis**: Tidak memindahkan semua kolom, melainkan membuat relasi yang tepat
- **📊 Data Integrity**: Foreign key constraints untuk menjaga konsistensi data
- **⚡ Performance**: Indexes untuk query yang optimal
- **📈 Reporting**: Views untuk laporan yang komprehensif
- **🛡️ Scalability**: Struktur yang dapat dikembangkan untuk kebutuhan masa depan

### **Fitur yang Tersedia:**
- **📋 Complete Installation Report**: Laporan lengkap dengan semua data terkait
- **👥 Customer Summary**: Summary instalasi per customer
- **📦 Asset Report**: Laporan aset per instalasi
- **👨‍🔧 Technician Report**: Laporan performa teknisi

**Sistem database siap untuk fitur "Add Report Installation" dengan struktur yang logis dan efisien!** 🎉
