# NetworkDevice Assets Relation Fix

## Problem
Error saat membuat installation report:
```
Assets: unsupported relations for schema NetworkDevice
```

## Root Cause
Model `NetworkDevice` memiliki field `AssetsID` tapi tidak memiliki relasi `Assets` yang didefinisikan, sehingga GORM tidak bisa memproses relasi tersebut.

## Solution
Menambahkan relasi `Assets` yang hilang di model `NetworkDevice`:

### File Modified:
- ✅ `crm-be/internal/models/entities/network_device_model.go`

### Added Relation:
```go
// Added this relation to NetworkDevice struct:
Assets                 *Asset                `json:"assets,omitempty" gorm:"foreignKey:AssetsID;references:ID"`
```

### Fixed NetworkDevice Model:
```go
type NetworkDevice struct {
    ID                     string                `json:"id" gorm:"primaryKey;type:varchar(191)"`
    CustomerID             string                `json:"customer_id" gorm:"type:varchar(191);not null"`
    Customer               *Customer             `json:"customer,omitempty" gorm:"foreignKey:CustomerID;references:ID"`
    AssetsID               string                `json:"assets_id" gorm:"type:varchar(191);not null"`
    Assets                 *Asset                `json:"assets,omitempty" gorm:"foreignKey:AssetsID;references:ID"` // ✅ Added
    CustomerInstallationID *string               `json:"customer_installation_id,omitempty" gorm:"type:varchar(191);index:idx_network_devices_customer_installation_id"`
    CustomerInstallation   *CustomerInstallation `json:"customer_installation,omitempty" gorm:"foreignKey:CustomerInstallationID;references:ID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
    // ... other fields
}
```

## Database Schema
Relasi ini sesuai dengan database schema:
- `network_devices.assets_id` → `assets.id`
- Foreign key constraint: `network_devices_ibfk_2` FOREIGN KEY (`assets_id`) REFERENCES `assets` (`id`)

## Status:
✅ **Relation fixed** - NetworkDevice now has proper Assets relation
✅ **GORM compatible** - Relation properly defined with foreign key
✅ **Database aligned** - Matches actual database schema
✅ **Ready for testing** - Installation report should work now

## Test:
```bash
cd crm-be
go run cmd/myapp/main.go
# Should start without errors
```

## Expected Result:
Installation report creation should now work without "unsupported relations" error.
