# Document Fields untuk Customer Installations

## 📋 **Fitur Dokumen yang Ditambahkan**

### ✅ **1. Kolom Document Type (Dropdown)**
- **Nama Kolom:** `document_type`
- **Tipe Data:** `ENUM('KTP', 'SIM', 'Paspor')`
- **Fungsi:** Dropdown pilihan jenis dokumen identitas
- **Validasi:** Hanya menerima nilai KTP, SIM, atau Paspor

### ✅ **2. Kolom Document Photo**
- **Nama Kolom:** `document_photo`
- **Tipe Data:** `VARCHAR(255)`
- **Fungsi:** Menyimpan path file foto dokumen identitas
- **Format:** Path relatif ke folder uploads

## 🗄️ **Struktur Database**

### **Tabel customer_installations (Updated)**
```sql
+----------------+--------------------------------+------+-----+----------------------+
| Field          | Type                           | Null | Key | Default              |
+----------------+--------------------------------+------+-----+----------------------+
| id             | varchar(191)                   | NO   | PRI | NULL                 | ✅ Primary Key
| customer_id    | varchar(191)                   | YES  | MUL | NULL                 | ✅ Foreign Key
| technician_id  | varchar(191)                   | YES  | MUL | NULL                 | ✅ Foreign Key
| status         | varchar(191)                   | YES  | MUL | pending              | ✅ Status
| notes          | text                           | YES  |     | NULL                 | ✅ Notes
| document_type  | enum('KTP','SIM','Paspor')     | YES  | MUL | NULL                 | 🆕 Document Type
| document_photo | varchar(255)                   | YES  |     | NULL                 | 🆕 Document Photo
| on_air_date    | date                           | YES  | MUL | NULL                 | ✅ On Air Date
| createdAt      | datetime(3)                    | NO   |     | CURRENT_TIMESTAMP(3) | ✅ Created At
| updatedAt      | datetime(3)                    | NO   |     | CURRENT_TIMESTAMP(3) | ✅ Updated At
+----------------+--------------------------------+------+-----+----------------------+
```

### **Indexes yang Ditambahkan**
```sql
-- Index untuk document_type untuk performa query
KEY `idx_customer_installations_document_type` (`document_type`)
```

## 🔧 **Perubahan Kode**

### **1. Entity Model (user_model.go)**
```go
type CustomerInstallation struct {
    ID            string     `gorm:"column:id;type:varchar;primaryKey" json:"id"`
    CustomerID    *string    `gorm:"column:customer_id;type:varchar;index:idx_customer_installations_customer_id" json:"customer_id,omitempty"`
    TechnicianID  *string    `gorm:"column:technician_id;type:varchar;index:idx_customer_installations_technician_id" json:"technician_id,omitempty"`
    Status        string     `gorm:"column:status;type:varchar;default:'pending'" json:"status,omitempty"`
    Notes         string     `gorm:"column:notes;type:text" json:"notes,omitempty"`
    DocumentType  *string    `gorm:"column:document_type;type:enum('KTP','SIM','Paspor')" json:"document_type,omitempty"`     // 🆕
    DocumentPhoto *string    `gorm:"column:document_photo;type:varchar(255)" json:"document_photo,omitempty"`               // 🆕
    OnAirDate     *time.Time `gorm:"column:on_air_date;type:date" json:"on_air_date,omitempty"`
    CreatedAt     time.Time  `gorm:"column:createdAt;default:current_timestamp" json:"createdAt"`
    UpdatedAt     time.Time  `gorm:"column:updatedAt;default:current_timestamp on update current_timestamp" json:"updatedAt"`
    Customer      *Customer  `gorm:"foreignKey:CustomerID;references:id;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"customer,omitempty"`
    Technician    *User      `gorm:"foreignKey:TechnicianID;references:id;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"technician,omitempty"`
    Images        []Image    `gorm:"foreignKey:ArchiveInstallationId;references:ID" json:"images,omitempty"`
}
```

### **2. Request Struct (request.go)**
```go
type CreateAdminCustomerInstallationRequest struct {
    CustomerId    string   `json:"customer_id" validate:"required"`
    TechnicianId  string   `json:"technician_id" validate:"required"`
    Status        string   `json:"status"`
    Notes         string   `json:"notes"`
    DocumentType  string   `json:"document_type" validate:"omitempty,oneof=KTP SIM Paspor"`  // 🆕
    DocumentPhoto string   `json:"document_photo"`                                            // 🆕
    ImageIds      []string `json:"image_ids" validate:"required"`
    OnAirDate     string   `json:"on_air_date"`
}
```

### **3. Repository Logic (repository.go)**
```go
// Handle document_type - convert to pointer if not empty
var documentType *string
if request.DocumentType != "" {
    documentType = &request.DocumentType
}

// Handle document_photo - convert to pointer if not empty
var documentPhoto *string
if request.DocumentPhoto != "" {
    documentPhoto = &request.DocumentPhoto
}

customerInstallation := entities.CustomerInstallation{
    ID:            "",
    CustomerID:    &request.CustomerId,
    TechnicianID:  &request.TechnicianId,
    Status:        request.Status,
    Notes:         request.Notes,
    DocumentType:  documentType,    // 🆕
    DocumentPhoto: documentPhoto,   // 🆕
    OnAirDate:     onAirDate,
}
```

## 🌐 **API Usage**

### **Request Body Example**
```json
{
    "customer_id": "24225552-b0c3-43d7-9a67-20e04f36fa5f",
    "technician_id": "c13a6c87-ec28-47ba-84c2-58b5ace2af57",
    "status": "completed",
    "notes": "Installation completed with KTP verification",
    "document_type": "KTP",
    "document_photo": "uploads/documents/ktp_24225552_20250922.jpg",
    "image_ids": ["image1", "image2"],
    "on_air_date": "2025-09-22"
}
```

### **Response Body Example**
```json
{
    "id": "949189b8-97d1-11f0-9fa2-d843ae0f1e06",
    "customer_id": "24225552-b0c3-43d7-9a67-20e04f36fa5f",
    "technician_id": "c13a6c87-ec28-47ba-84c2-58b5ace2af57",
    "status": "completed",
    "notes": "Installation completed with KTP verification",
    "document_type": "KTP",
    "document_photo": "uploads/documents/ktp_24225552_20250922.jpg",
    "on_air_date": "2025-09-22T00:00:00Z",
    "createdAt": "2025-09-22T23:31:24.000Z",
    "updatedAt": "2025-09-22T23:31:24.000Z"
}
```

## 📱 **Frontend Implementation**

### **Dropdown Document Type**
```html
<select name="document_type" v-model="form.document_type">
    <option value="">Pilih Jenis Dokumen</option>
    <option value="KTP">KTP</option>
    <option value="SIM">SIM</option>
    <option value="Paspor">Paspor</option>
</select>
```

### **File Upload Document Photo**
```html
<input 
    type="file" 
    name="document_photo" 
    accept="image/*"
    @change="handleDocumentPhotoUpload"
/>
```

### **Validation Rules**
```javascript
const validationRules = {
    document_type: {
        required: false,
        enum: ['KTP', 'SIM', 'Paspor']
    },
    document_photo: {
        required: false,
        type: 'string',
        maxLength: 255
    }
}
```

## 🎯 **Use Cases**

### **1. Customer Installation dengan KTP**
- Customer memilih "KTP" dari dropdown
- Upload foto KTP
- System menyimpan path foto dan jenis dokumen

### **2. Customer Installation dengan SIM**
- Customer memilih "SIM" dari dropdown
- Upload foto SIM
- System menyimpan path foto dan jenis dokumen

### **3. Customer Installation dengan Paspor**
- Customer memilih "Paspor" dari dropdown
- Upload foto Paspor
- System menyimpan path foto dan jenis dokumen

### **4. Installation tanpa Dokumen**
- Field document_type dan document_photo bisa kosong
- System tetap bisa menyimpan data instalasi

## 🔍 **Query Examples**

### **Cari Installation berdasarkan Document Type**
```sql
SELECT * FROM customer_installations 
WHERE document_type = 'KTP';
```

### **Cari Installation dengan Document Photo**
```sql
SELECT * FROM customer_installations 
WHERE document_photo IS NOT NULL;
```

### **Statistik Document Type**
```sql
SELECT 
    document_type,
    COUNT(*) as count
FROM customer_installations 
WHERE document_type IS NOT NULL
GROUP BY document_type;
```

## 📁 **File Structure untuk Document Photos**

```
uploads/
├── documents/
│   ├── ktp_24225552_20250922.jpg
│   ├── sim_55f03555_20250922.jpg
│   └── paspor_9b8c7e59_20250922.jpg
└── installations/
    ├── image1.jpg
    └── image2.jpg
```

## ✅ **Verifikasi**

- ✅ Kolom `document_type` berhasil ditambahkan dengan ENUM
- ✅ Kolom `document_photo` berhasil ditambahkan dengan VARCHAR(255)
- ✅ Index untuk `document_type` berhasil dibuat
- ✅ Entity model sesuai dengan database schema
- ✅ Request struct dengan validasi yang tepat
- ✅ Repository logic menangani field baru dengan benar
- ✅ Aplikasi build berhasil tanpa error

## 🚀 **Hasil Akhir**

Fitur dokumen untuk Customer Installations berhasil ditambahkan dengan:
- **Dropdown pilihan dokumen** (KTP, SIM, Paspor)
- **Upload foto dokumen** dengan path storage
- **Validasi yang tepat** untuk input data
- **Database structure yang optimal** dengan indexes
- **API yang siap digunakan** untuk frontend integration

Fitur dokumen siap digunakan! 🎉
