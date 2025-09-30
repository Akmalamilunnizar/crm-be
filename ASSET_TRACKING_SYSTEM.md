# Asset Tracking System untuk Customer Installations

## 🎯 **Overview Sistem**

Sistem tracking aset ini dirancang untuk fitur **"Add Report Installation"** yang memungkinkan:
- **Aset Keluar** - Tracking aset yang dikeluarkan untuk instalasi customer
- **Aset Masuk** - Tracking aset yang dikembalikan dari instalasi customer
- **Relasi dengan Tabel Assets** - Integrasi penuh dengan sistem manajemen aset
- **Laporan Komprehensif** - View dan query untuk monitoring aset

## 🗄️ **Struktur Database**

### **1. Tabel `asset_transactions` (Baru)**
```sql
+---------------------------+--------------------------------+------+-----+----------------------+
| Field                     | Type                           | Null | Key | Default              |
+---------------------------+--------------------------------+------+-----+----------------------+
| id                        | varchar(191)                   | NO   | PRI | NULL                 | ✅ Primary Key
| customer_installation_id  | varchar(191)                   | NO   | MUL | NULL                 | ✅ Foreign Key
| asset_id                  | varchar(191)                   | NO   | MUL | NULL                 | ✅ Foreign Key
| transaction_type          | enum('out','in')               | NO   | MUL | NULL                 | ✅ Type Transaction
| quantity                  | int                            | NO   |     | 1                    | ✅ Jumlah Aset
| notes                     | text                           | YES  |     | NULL                 | ✅ Catatan
| transaction_date          | datetime(3)                    | NO   | MUL | CURRENT_TIMESTAMP(3) | ✅ Tanggal Transaksi
| created_by                | varchar(191)                   | NO   | MUL | NULL                 | ✅ User Creator
| createdAt                 | datetime(3)                    | NO   |     | CURRENT_TIMESTAMP(3) | ✅ Created At
| updatedAt                 | datetime(3)                    | NO   |     | CURRENT_TIMESTAMP(3) | ✅ Updated At
+---------------------------+--------------------------------+------+-----+----------------------+
```

### **2. Tabel `customer_installations` (Updated)**
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
| installation_type         | enum('new_installation','maintenance','upgrade','downgrade') | YES  | MUL | new_installation     | 🆕 Installation Type
| total_assets_out          | int                                                          | YES  |     | 0                    | 🆕 Total Assets Out
| total_assets_in           | int                                                          | YES  |     | 0                    | 🆕 Total Assets In
| installation_completed_at | datetime(3)                                                  | YES  | MUL | NULL                 | 🆕 Completion Date
| on_air_date               | date                                                         | YES  | MUL | NULL                 | ✅ On Air Date
| createdAt                 | datetime(3)                                                  | NO   |     | CURRENT_TIMESTAMP(3) | ✅ Created At
| updatedAt                 | datetime(3)                                                  | NO   |     | CURRENT_TIMESTAMP(3) | ✅ Updated At
+---------------------------+--------------------------------------------------------------+------+-----+----------------------+
```

## 🔗 **Relasi Database**

### **Foreign Key Constraints**
```sql
-- Asset Transactions
CONSTRAINT `fk_asset_transactions_customer_installation` 
FOREIGN KEY (`customer_installation_id`) REFERENCES `customer_installations` (`id`) 
ON DELETE CASCADE ON UPDATE CASCADE

CONSTRAINT `fk_asset_transactions_asset` 
FOREIGN KEY (`asset_id`) REFERENCES `assets` (`id`) 
ON DELETE CASCADE ON UPDATE CASCADE

CONSTRAINT `fk_asset_transactions_user` 
FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) 
ON DELETE CASCADE ON UPDATE CASCADE
```

### **Indexes untuk Performance**
```sql
-- Asset Transactions Indexes
KEY `idx_asset_transactions_customer_installation_id` (`customer_installation_id`)
KEY `idx_asset_transactions_asset_id` (`asset_id`)
KEY `idx_asset_transactions_transaction_type` (`transaction_type`)
KEY `idx_asset_transactions_transaction_date` (`transaction_date`)
KEY `idx_asset_transactions_created_by` (`created_by`)

-- Customer Installations New Indexes
KEY `idx_customer_installations_installation_type` (`installation_type`)
KEY `idx_customer_installations_installation_completed_at` (`installation_completed_at`)
```

## 📊 **Views untuk Laporan**

### **1. View `asset_transaction_report`**
```sql
CREATE OR REPLACE VIEW `asset_transaction_report` AS
SELECT 
    at.id as transaction_id,
    at.customer_installation_id,
    ci.customer_id,
    c.name as customer_name,
    at.asset_id,
    a.brand,
    a.type as asset_type,
    a.model,
    a.serial_number,
    at.transaction_type,
    at.quantity,
    at.notes as transaction_notes,
    at.transaction_date,
    at.created_by,
    u.name as created_by_name,
    ci.installation_type,
    ci.status as installation_status,
    ci.on_air_date
FROM asset_transactions at
JOIN customer_installations ci ON at.customer_installation_id = ci.id
JOIN customer c ON ci.customer_id = c.id
JOIN assets a ON at.asset_id = a.id
JOIN users u ON at.created_by = u.id
ORDER BY at.transaction_date DESC;
```

### **2. View `installation_asset_summary`**
```sql
CREATE OR REPLACE VIEW `installation_asset_summary` AS
SELECT 
    ci.id as installation_id,
    ci.customer_id,
    c.name as customer_name,
    ci.installation_type,
    ci.status,
    ci.on_air_date,
    ci.installation_completed_at,
    COUNT(CASE WHEN at.transaction_type = 'out' THEN 1 END) as total_assets_out,
    COUNT(CASE WHEN at.transaction_type = 'in' THEN 1 END) as total_assets_in,
    SUM(CASE WHEN at.transaction_type = 'out' THEN at.quantity ELSE 0 END) as total_quantity_out,
    SUM(CASE WHEN at.transaction_type = 'in' THEN at.quantity ELSE 0 END) as total_quantity_in
FROM customer_installations ci
LEFT JOIN asset_transactions at ON ci.id = at.customer_installation_id
LEFT JOIN customer c ON ci.customer_id = c.id
GROUP BY ci.id, ci.customer_id, c.name, ci.installation_type, ci.status, ci.on_air_date, ci.installation_completed_at;
```

## 🔧 **Entity Models**

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
    OnAirDate               *time.Time          `gorm:"column:on_air_date;type:date" json:"on_air_date,omitempty"`
    CreatedAt               time.Time           `gorm:"column:createdAt;default:current_timestamp" json:"createdAt"`
    UpdatedAt               time.Time           `gorm:"column:updatedAt;default:current_timestamp on update current_timestamp" json:"updatedAt"`
    Customer                *Customer           `gorm:"foreignKey:CustomerID;references:id;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"customer,omitempty"`
    Technician              *User               `gorm:"foreignKey:TechnicianID;references:id;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"technician,omitempty"`
    Images                  []Image             `gorm:"foreignKey:ArchiveInstallationId;references:ID" json:"images,omitempty"`
    AssetTransactions       []AssetTransaction  `gorm:"foreignKey:CustomerInstallationID;references:ID" json:"asset_transactions,omitempty"`
}
```

### **2. AssetTransaction (New)**
```go
type AssetTransaction struct {
    ID                      string             `gorm:"column:id;type:varchar;primaryKey" json:"id"`
    CustomerInstallationID  string             `gorm:"column:customer_installation_id;type:varchar;index:idx_asset_transactions_customer_installation_id" json:"customer_installation_id"`
    AssetID                 string             `gorm:"column:asset_id;type:varchar;index:idx_asset_transactions_asset_id" json:"asset_id"`
    TransactionType         string             `gorm:"column:transaction_type;type:enum('out','in')" json:"transaction_type"`
    Quantity                int                `gorm:"column:quantity;type:int;default:1" json:"quantity"`
    Notes                   string             `gorm:"column:notes;type:text" json:"notes,omitempty"`
    TransactionDate         time.Time          `gorm:"column:transaction_date;type:datetime(3);default:current_timestamp" json:"transaction_date"`
    CreatedBy               string             `gorm:"column:created_by;type:varchar;index:idx_asset_transactions_created_by" json:"created_by"`
    CreatedAt               time.Time          `gorm:"column:createdAt;default:current_timestamp" json:"createdAt"`
    UpdatedAt               time.Time          `gorm:"column:updatedAt;default:current_timestamp on update current_timestamp" json:"updatedAt"`
    CustomerInstallation    *CustomerInstallation `gorm:"foreignKey:CustomerInstallationID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"customer_installation,omitempty"`
    Asset                   *Asset             `gorm:"foreignKey:AssetID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"asset,omitempty"`
    User                    *User              `gorm:"foreignKey:CreatedBy;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"user,omitempty"`
}
```

## 🌐 **API Usage Examples**

### **1. Create Asset Transaction (Aset Keluar)**
```json
POST /api/asset-transactions
{
    "customer_installation_id": "949189b8-97d1-11f0-9fa2-d843ae0f1e06",
    "asset_id": "5ca1606b-66b3-4958-b7af-f48d4cda800a",
    "transaction_type": "out",
    "quantity": 1,
    "notes": "Router keluar untuk instalasi customer baru",
    "created_by": "c13a6c87-ec28-47ba-84c2-58b5ace2af57"
}
```

### **2. Create Asset Transaction (Aset Masuk)**
```json
POST /api/asset-transactions
{
    "customer_installation_id": "949189b8-97d1-11f0-9fa2-d843ae0f1e06",
    "asset_id": "5ca1606b-66b3-4958-b7af-f48d4cda800a",
    "transaction_type": "in",
    "quantity": 1,
    "notes": "Router dikembalikan setelah maintenance",
    "created_by": "c13a6c87-ec28-47ba-84c2-58b5ace2af57"
}
```

### **3. Get Asset Transaction Report**
```json
GET /api/asset-transaction-report
Response:
{
    "data": [
        {
            "transaction_id": "4a77cc83-9830-11f0-9fa2-d843ae0f1e06",
            "customer_installation_id": "949189b8-97d1-11f0-9fa2-d843ae0f1e06",
            "customer_name": "Cha Eunwoo",
            "asset_id": "5ca1606b-66b3-4958-b7af-f48d4cda800a",
            "brand": "TP-Link",
            "asset_type": "Router Edit",
            "model": "RBG128",
            "transaction_type": "out",
            "quantity": 1,
            "transaction_notes": "Router keluar untuk instalasi customer baru",
            "transaction_date": "2025-09-23T10:49:22.000Z",
            "created_by_name": "teknisi",
            "installation_type": "new_installation",
            "installation_status": "completed"
        }
    ]
}
```

## 📱 **Frontend Implementation**

### **1. Form Add Report Installation**
```html
<form @submit="submitInstallationReport">
    <!-- Basic Installation Info -->
    <div class="form-group">
        <label>Customer</label>
        <select v-model="form.customer_id" required>
            <option v-for="customer in customers" :key="customer.id" :value="customer.id">
                {{ customer.name }}
            </option>
        </select>
    </div>

    <div class="form-group">
        <label>Installation Type</label>
        <select v-model="form.installation_type" required>
            <option value="new_installation">New Installation</option>
            <option value="maintenance">Maintenance</option>
            <option value="upgrade">Upgrade</option>
            <option value="downgrade">Downgrade</option>
        </select>
    </div>

    <!-- Asset Out Section -->
    <div class="asset-section">
        <h3>Assets Out</h3>
        <div v-for="(asset, index) in form.assets_out" :key="index" class="asset-item">
            <select v-model="asset.asset_id" required>
                <option v-for="availableAsset in availableAssets" :key="availableAsset.id" :value="availableAsset.id">
                    {{ availableAsset.brand }} {{ availableAsset.model }} ({{ availableAsset.serial_number }})
                </option>
            </select>
            <input type="number" v-model="asset.quantity" min="1" required>
            <textarea v-model="asset.notes" placeholder="Notes"></textarea>
            <button type="button" @click="removeAssetOut(index)">Remove</button>
        </div>
        <button type="button" @click="addAssetOut">Add Asset Out</button>
    </div>

    <!-- Asset In Section -->
    <div class="asset-section">
        <h3>Assets In</h3>
        <div v-for="(asset, index) in form.assets_in" :key="index" class="asset-item">
            <select v-model="asset.asset_id" required>
                <option v-for="outAsset in form.assets_out" :key="outAsset.asset_id" :value="outAsset.asset_id">
                    {{ getAssetName(outAsset.asset_id) }}
                </option>
            </select>
            <input type="number" v-model="asset.quantity" min="1" required>
            <textarea v-model="asset.notes" placeholder="Notes"></textarea>
            <button type="button" @click="removeAssetIn(index)">Remove</button>
        </div>
        <button type="button" @click="addAssetIn">Add Asset In</button>
    </div>

    <button type="submit">Submit Report</button>
</form>
```

### **2. Asset Transaction List**
```html
<template>
    <div class="asset-transaction-list">
        <h2>Asset Transactions</h2>
        <table>
            <thead>
                <tr>
                    <th>Date</th>
                    <th>Customer</th>
                    <th>Asset</th>
                    <th>Type</th>
                    <th>Quantity</th>
                    <th>Notes</th>
                    <th>Created By</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="transaction in transactions" :key="transaction.transaction_id">
                    <td>{{ formatDate(transaction.transaction_date) }}</td>
                    <td>{{ transaction.customer_name }}</td>
                    <td>{{ transaction.brand }} {{ transaction.model }}</td>
                    <td>
                        <span :class="transaction.transaction_type === 'out' ? 'badge-out' : 'badge-in'">
                            {{ transaction.transaction_type.toUpperCase() }}
                        </span>
                    </td>
                    <td>{{ transaction.quantity }}</td>
                    <td>{{ transaction.transaction_notes }}</td>
                    <td>{{ transaction.created_by_name }}</td>
                </tr>
            </tbody>
        </table>
    </div>
</template>
```

## 🔍 **Query Examples**

### **1. Get Assets Out for Installation**
```sql
SELECT 
    at.*,
    a.brand,
    a.model,
    a.serial_number
FROM asset_transactions at
JOIN assets a ON at.asset_id = a.id
WHERE at.customer_installation_id = '949189b8-97d1-11f0-9fa2-d843ae0f1e06'
AND at.transaction_type = 'out';
```

### **2. Get Assets In for Installation**
```sql
SELECT 
    at.*,
    a.brand,
    a.model,
    a.serial_number
FROM asset_transactions at
JOIN assets a ON at.asset_id = a.id
WHERE at.customer_installation_id = '949189b8-97d1-11f0-9fa2-d843ae0f1e06'
AND at.transaction_type = 'in';
```

### **3. Get Installation Asset Summary**
```sql
SELECT * FROM installation_asset_summary 
WHERE installation_id = '949189b8-97d1-11f0-9fa2-d843ae0f1e06';
```

### **4. Get Asset Transaction Report**
```sql
SELECT * FROM asset_transaction_report 
WHERE customer_installation_id = '949189b8-97d1-11f0-9fa2-d843ae0f1e06'
ORDER BY transaction_date DESC;
```

## ✅ **Verifikasi Sistem**

- ✅ Tabel `asset_transactions` berhasil dibuat dengan struktur lengkap
- ✅ Tabel `customer_installations` berhasil diupdate dengan field baru
- ✅ Foreign key constraints berfungsi dengan baik
- ✅ Indexes untuk performa query optimal
- ✅ Views untuk laporan berhasil dibuat
- ✅ Sample data berhasil dimasukkan
- ✅ Entity models sesuai dengan database schema
- ✅ Relasi antar tabel berfungsi dengan baik

## 🚀 **Hasil Akhir**

Sistem Asset Tracking untuk Customer Installations berhasil dibuat dengan:
- **📊 Tracking Aset Keluar/Masuk** - Sistem lengkap untuk monitoring aset
- **🔗 Relasi dengan Assets** - Integrasi penuh dengan tabel assets
- **📈 Laporan Komprehensif** - Views untuk monitoring dan reporting
- **🎯 Fitur Add Report Installation** - Siap untuk implementasi frontend
- **⚡ Performance Optimal** - Indexes dan struktur database yang efisien

Sistem Asset Tracking siap digunakan untuk fitur Add Report Installation! 🎉
