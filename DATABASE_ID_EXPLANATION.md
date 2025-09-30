# Penjelasan Fungsi ID di Database CRM

## Fungsi ID (Primary Key) di Setiap Tabel Database

### 1. **Uniqueness (Keunikan)**
- Setiap record memiliki identifier yang unik
- Mencegah duplikasi data
- Memastikan integritas data

### 2. **Performance (Kinerja)**
- Indexing otomatis untuk pencarian yang cepat
- Optimasi query database
- Meningkatkan kecepatan operasi CRUD

### 3. **Relationships (Relasi)**
- Sebagai foreign key untuk menghubungkan antar tabel
- Memungkinkan JOIN operations yang efisien
- Menjaga referential integrity

### 4. **API Operations**
- Endpoint untuk operasi CRUD (Create, Read, Update, Delete)
- Identifikasi unik untuk operasi REST API
- Memudahkan frontend untuk mengelola data

### 5. **Data Integrity**
- Mencegah data yang tidak konsisten
- Memudahkan maintenance dan debugging
- Standar database design yang baik

## Masalah yang Ditemukan dan Solusi

### Masalah: Tabel `customer_installations` Tidak Memiliki ID

**Sebelum:**
```sql
CREATE TABLE IF NOT EXISTS `customer_installations` (
  `date_trial` date NOT NULL,
  `on_air_date` date DEFAULT NULL,
  KEY `idx_customer_installations_on_air_date` (`on_air_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

**Masalah:**
- ❌ Tidak ada Primary Key (ID)
- ❌ Tidak ada relasi ke tabel lain
- ❌ Tidak ada timestamp fields
- ❌ Struktur tidak sesuai dengan entity model

**Sesudah (Solusi):**
```sql
CREATE TABLE IF NOT EXISTS `customer_installations` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `customer_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `technician_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `date_trial` date NOT NULL,
  `status` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'pending',
  `equipment_used` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `completion_time` int DEFAULT 0,
  `customer_signature` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `on_air_date` date DEFAULT NULL,
  `createdAt` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updatedAt` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_customer_installations_customer_id` (`customer_id`),
  KEY `idx_customer_installations_technician_id` (`technician_id`),
  KEY `idx_customer_installations_on_air_date` (`on_air_date`),
  KEY `idx_customer_installations_status` (`status`),
  CONSTRAINT `fk_customer_installations_customer` FOREIGN KEY (`customer_id`) REFERENCES `customer` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_customer_installations_technician` FOREIGN KEY (`technician_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

**Keuntungan:**
- ✅ Primary Key (ID) untuk identifikasi unik
- ✅ Relasi ke tabel `customer` dan `users`
- ✅ Timestamp fields untuk audit trail
- ✅ Indexes untuk performa optimal
- ✅ Foreign key constraints untuk data integrity
- ✅ Struktur sesuai dengan entity model

## Perubahan Entity Model

### CustomerInstallation Struct (Updated)

```go
type CustomerInstallation struct {
    ID                string     `gorm:"column:id;type:varchar;primaryKey" json:"id"`
    CustomerID        *string    `gorm:"column:customer_id;type:varchar;index:idx_customer_installations_customer_id" json:"customer_id,omitempty"`
    TechnicianID      *string    `gorm:"column:technician_id;type:varchar;index:idx_customer_installations_technician_id" json:"technician_id,omitempty"`
    Description       string     `gorm:"column:description;type:text" json:"description,omitempty"`
    DateTrial         time.Time  `gorm:"column:date_trial;type:date" json:"date_trial"`
    Status            string     `gorm:"column:status;type:varchar;default:'pending'" json:"status,omitempty"`
    EquipmentUsed     string     `gorm:"column:equipment_used;type:text" json:"equipment_used,omitempty"`
    Notes             string     `gorm:"column:notes;type:text" json:"notes,omitempty"`
    CompletionTime    int        `gorm:"column:completion_time;type:int;default:0" json:"completion_time,omitempty"`
    CustomerSignature string     `gorm:"column:customer_signature;type:text" json:"customer_signature,omitempty"`
    OnAirDate         *time.Time `gorm:"column:on_air_date;type:date" json:"on_air_date,omitempty"`
    CreatedAt         time.Time  `gorm:"column:createdAt;default:current_timestamp" json:"createdAt"`
    UpdatedAt         time.Time  `gorm:"column:updatedAt;default:current_timestamp on update current_timestamp" json:"updatedAt"`
    Customer          *Customer  `gorm:"foreignKey:CustomerID;references:id;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"customer,omitempty"`
    Technician        *User      `gorm:"foreignKey:TechnicianID;references:id;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"technician,omitempty"`
    Images            []Image    `gorm:"foreignKey:ArchiveInstallationId;references:ID" json:"images,omitempty"`
}
```

## Cara Menjalankan Migrasi

1. **Backup Database** (Opsional tapi disarankan)
2. **Jalankan Script Migrasi:**
   ```bash
   mysql -u username -p database_name < fix-customer-installations-table.sql
   ```
3. **Verifikasi:**
   ```sql
   DESCRIBE customer_installations;
   SELECT COUNT(*) FROM customer_installations;
   ```

## Relasi yang Dibuat

### 1. CustomerInstallation → Customer
- **Foreign Key:** `customer_id` → `customer.id`
- **Constraint:** `ON DELETE SET NULL, ON UPDATE CASCADE`
- **Fungsi:** Menghubungkan instalasi dengan customer

### 2. CustomerInstallation → User (Technician)
- **Foreign Key:** `technician_id` → `users.id`
- **Constraint:** `ON DELETE SET NULL, ON UPDATE CASCADE`
- **Fungsi:** Menghubungkan instalasi dengan teknisi yang bertugas

### 3. CustomerInstallation → Images
- **Foreign Key:** `images.archive_installation_id` → `customer_installations.id`
- **Fungsi:** Menyimpan foto dokumentasi instalasi

## Indexes yang Dibuat

1. `idx_customer_installations_customer_id` - Untuk query berdasarkan customer
2. `idx_customer_installations_technician_id` - Untuk query berdasarkan teknisi
3. `idx_customer_installations_on_air_date` - Untuk query berdasarkan tanggal on-air
4. `idx_customer_installations_status` - Untuk query berdasarkan status instalasi

## Kesimpulan

Dengan menambahkan ID dan relasi yang proper ke tabel `customer_installations`, sekarang:

1. **Data Integrity** terjaga dengan foreign key constraints
2. **Performance** meningkat dengan indexes yang tepat
3. **API Operations** bisa dilakukan dengan mudah
4. **Relasi** antar tabel berfungsi dengan baik
5. **Audit Trail** tersedia dengan timestamp fields

Struktur database sekarang sudah sesuai dengan best practices dan entity model yang ada.
