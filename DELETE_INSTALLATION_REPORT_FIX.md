# DELETE INSTALLATION REPORT FIX

## Problem
Error saat menghapus Customer Installation Report:
```json
{
    "data": "",
    "message": "Area: unsupported relations for schema CustomerInstallation",
    "success": false
}
```

## Root Cause
Di repository `DeleteAdminCustomerInstallationRepository` dan `UpdateAdminCustomerInstallationRepository`, ada penggunaan `Preload("Area")` yang tidak valid karena:

1. **CustomerInstallation model** tidak memiliki relasi langsung dengan **Area**
2. **Area relation** hanya ada di **Customer model**, bukan di **CustomerInstallation**
3. GORM tidak dapat melakukan preload untuk relasi yang tidak ada

## Solution
Menghapus `Preload("Area")` dari fungsi-fungsi berikut:

### 1. Delete Function
**File:** `crm-be/internal/api/admin/customer/installation/repository.go`

**Before:**
```go
func (r AdminCustomerInstallationRepositoryStruct) DeleteAdminCustomerInstallationRepository(request IdAdminCustomerInstallationRequest) (entities.CustomerInstallation, error) {
	customer := entities.CustomerInstallation{}
	tx := r.db.Preload("Area").Find(&customer, "id = ?", request.Id) // ❌ ERROR
	// ...
}
```

**After:**
```go
func (r AdminCustomerInstallationRepositoryStruct) DeleteAdminCustomerInstallationRepository(request IdAdminCustomerInstallationRequest) (entities.CustomerInstallation, error) {
	customer := entities.CustomerInstallation{}
	tx := r.db.Find(&customer, "id = ?", request.Id) // ✅ FIXED
	// ...
}
```

### 2. Update Function
**File:** `crm-be/internal/api/admin/customer/installation/repository.go`

**Before:**
```go
func (r AdminCustomerInstallationRepositoryStruct) UpdateAdminCustomerInstallationRepository(request UpdateAdminCustomerInstallationRequest) (entities.CustomerInstallation, error) {
	customer := entities.CustomerInstallation{}
	tx := r.db.Preload("Area").Find(&customer, "id = ?", request.Id) // ❌ ERROR
	// ...
}
```

**After:**
```go
func (r AdminCustomerInstallationRepositoryStruct) UpdateAdminCustomerInstallationRepository(request UpdateAdminCustomerInstallationRequest) (entities.CustomerInstallation, error) {
	customer := entities.CustomerInstallation{}
	tx := r.db.Find(&customer, "id = ?", request.Id) // ✅ FIXED
	// ...
}
```

## Model Relations
### CustomerInstallation Model
```go
type CustomerInstallation struct {
    // ... fields ...
    
    // Relations
    Customer          *Customer          `gorm:"foreignKey:CustomerID;references:id"`
    Technician        *User              `gorm:"foreignKey:TechnicianID;references:id"`
    Images            []Image            `gorm:"foreignKey:ArchiveInstallationId;references:ID"`
    AssetTransactions []AssetTransaction `gorm:"foreignKey:CustomerInstallationID;references:ID"`
    NetworkDevices    []NetworkDevice    `gorm:"foreignKey:CustomerInstallationID;references:ID"`
    CustomerServices  []CustomerService  `gorm:"foreignKey:CustomerInstallationID;references:ID"`
    Cables            []Cable            `gorm:"foreignKey:CustomerInstallationID;references:ID"`
    // ❌ NO Area relation
}
```

### Customer Model
```go
type Customer struct {
    // ... fields ...
    AreaID string `gorm:"column:area_id" json:"area_id"`
    Area   *Areas  `gorm:"foreignKey:AreaID" json:"area"` // ✅ Area relation exists here
    // ...
}
```

## Frontend Implementation
### API Function
**File:** `crm-fe/api/admin/customer.ts`

```typescript
deleteInstallationReport: async (installationId: string) => {
  const response = await fetch(`${api}/api/admin/customer-installation/${installationId}`, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${useCookie("token").value}`,
    },
  });
  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.message || 'Failed to delete installation report');
  }
  return response.json();
},
```

### Delete Function with Confirmation
**File:** `crm-fe/pages/dashboard/report/customer-installation/reports.vue`

```typescript
async function deleteReport(installationId: string) {
  // Show confirmation dialog
  const confirmed = window.confirm(
    "Are you sure you want to delete this installation report? This action cannot be undone."
  );
  
  if (!confirmed) return;

  try {
    await customerAdminApi().deleteInstallationReport(installationId);
    
    // Show success notification
    useToast().add({
      title: "Success!",
      description: "Installation report deleted successfully",
      color: "green",
    });
    
    // Reload reports to reflect the changes
    await loadReports();
  } catch (error) {
    console.error("Error deleting installation report:", error);
    
    // Show error notification
    useToast().add({
      title: "Error",
      description: error instanceof Error ? error.message : "Failed to delete installation report",
      color: "red",
    });
  }
}
```

### Delete Button in Actions
```vue
<UButton @click="deleteReport(report.installation_id)" size="sm" color="red" variant="outline">
  Delete
</UButton>
```

## Testing
1. ✅ Backend compiles successfully
2. ✅ No more "Area: unsupported relations" error
3. ✅ Delete endpoint works properly
4. ✅ Frontend delete button functional
5. ✅ Confirmation dialog prevents accidental deletion
6. ✅ Success/error notifications work

## Files Modified
1. `crm-be/internal/api/admin/customer/installation/repository.go` - Fixed preload issue
2. `crm-fe/api/admin/customer.ts` - Added deleteInstallationReport function
3. `crm-fe/pages/dashboard/report/customer-installation/reports.vue` - Added delete button and function

## Result
✅ **FIXED**: Customer Installation Report delete functionality now works without errors!


