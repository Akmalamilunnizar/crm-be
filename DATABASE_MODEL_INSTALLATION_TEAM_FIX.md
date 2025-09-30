# Database Model Installation Team Fix

## Problem
Error saat membuat installation report:
```
Error 1054 (42S22): Unknown column 'installation_team_name' in 'field list'
```

## Root Cause
Field `InstallationTeamName` dan `InstallationTeamPhone` masih didefinisikan di database model `CustomerService` di file `customer_service_model.go`, sehingga GORM mencoba untuk insert ke kolom yang tidak ada di database.

## Solution
Menghapus field installation team dari database model:

### File Modified:
- ✅ `crm-be/internal/models/entities/customer_service_model.go`

### Removed Fields:
```go
// Removed these fields from CustomerService struct:
InstallationTeamPhone  *string `gorm:"column:installation_team_phone;type:varchar(20)" json:"installation_team_phone,omitempty"`
InstallationTeamName   *string `gorm:"column:installation_team_name;type:varchar(191)" json:"installation_team_name,omitempty"`
```

### Fixed CustomerService Model:
```go
type CustomerService struct {
    ID                     string                `gorm:"column:id;type:varchar;primaryKey" json:"id"`
    CustomerID             string                `gorm:"column:customer_id;type:varchar;index:idx_customer_services_customer_id" json:"customer_id"`
    CustomerInstallationID *string               `gorm:"column:customer_installation_id;type:varchar;index:idx_customer_services_customer_installation_id" json:"customer_installation_id,omitempty"`
    DeviceID               *string               `gorm:"column:device_id;type:varchar;index:idx_customer_services_device_id" json:"device_id,omitempty"`
    CableID                *string               `gorm:"column:cable_id;type:varchar" json:"cable_id,omitempty"`
    CableLength            *float64              `gorm:"column:cable_length;type:decimal(10,2)" json:"cable_length,omitempty"`
    EndPortType            *string               `gorm:"column:end_port_type;type:varchar(50)" json:"end_port_type,omitempty"`
    UserLogin              *string               `gorm:"column:user_login;type:varchar(191)" json:"user_login,omitempty"`
    ServiceActivationDate  *time.Time            `gorm:"column:service_activation_date;type:date" json:"service_activation_date,omitempty"`
    Password               *string               `gorm:"column:password;type:varchar(191)" json:"password,omitempty"`
    UserStatus             string                `gorm:"column:user_status;type:enum('Active','Inactive','Suspended','Pending');default:'Active'" json:"user_status,omitempty"`
    InstallationNotes      *string               `gorm:"column:installation_notes;type:text" json:"installation_notes,omitempty"`
    // ✅ InstallationTeamPhone and InstallationTeamName removed
    CreatedAt              *time.Time            `gorm:"column:created_at;type:timestamp;default:current_timestamp" json:"created_at,omitempty"`
    UpdatedAt              *time.Time            `gorm:"column:updated_at;type:timestamp;default:current_timestamp on update current_timestamp" json:"updated_at,omitempty"`
    // ... relations
}
```

## Complete Fix Summary:
✅ **Frontend**: Removed installation team fields from form
✅ **Backend Request Types**: Removed installation team fields
✅ **Backend Repository**: Removed installation team assignments
✅ **Database Model**: Removed installation team fields from CustomerService struct

## Status:
✅ **Database model fixed** - No more unknown column errors
✅ **GORM compatible** - Model matches database schema
✅ **Ready for testing** - Installation report should work now

## Test:
```bash
cd crm-be
go run cmd/myapp/main.go
# Should start without errors
```

## Expected Result:
Installation report creation should now work without database column errors.
