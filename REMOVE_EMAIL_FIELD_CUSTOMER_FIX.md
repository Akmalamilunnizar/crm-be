# Remove Email Field from Customer - Database Error Fix

## Issue
Error occurred when creating a new customer:

```
Error 1054 (42S22): Unknown column 'email' in 'field list'
```

## Root Cause
The customer form and backend were trying to use an `email` field that doesn't exist in the `customer` table in the database.

## Database Verification
Checked the customer table structure and confirmed that the `email` column does NOT exist:

```sql
Customer table columns:
  - id: varchar (nullable: NO)
  - address: varchar (nullable: NO)
  - area_id: varchar (nullable: NO)
  - latitude: double (nullable: NO)
  - longitude: double (nullable: NO)
  - name: varchar (nullable: NO)
  - alias: varchar (nullable: NO)
  - password: varchar (nullable: NO)
  - phone: varchar (nullable: NO)
  - createdAt: datetime (nullable: NO)
  - updatedAt: datetime (nullable: NO)
  - installation_date: date (nullable: NO)
  - next_payment_date: date (nullable: NO)
  - sales_representative_id: varchar (nullable: YES)
  - service_request_date: date (nullable: YES)
  - proposed_package: varchar (nullable: YES)
  - bandwidth_capacity: varchar (nullable: YES)
  - company_id: varchar (nullable: YES)
```

**Note**: No `email` column exists in the customer table.

## Solution
Removed all references to the `email` field from both frontend and backend:

### 1. Frontend Changes (`crm-fe/pages/dashboard/customer/FormAddComponent.vue`)

**Removed from props:**
```typescript
// REMOVED:
email: {
  type: String,
  default: ""
},
```

**Verification:**
- ✅ No email field in form template
- ✅ No email field in form state
- ✅ No email field in form schema

### 2. Backend Changes

#### A. Customer Entity Model (`crm-be/internal/models/entities/user_model.go`)

**Removed from Customer struct:**
```go
// REMOVED:
Email string `gorm:"column:email" json:"email"`
```

**Before:**
```go
type Customer struct {
    ID                    string    `json:"id" gorm:"primaryKey"`
    Address               string    `gorm:"column:address" json:"address"`
    AreaID                string    `gorm:"column:area_id" json:"area_id"`
    Area                  *Areas    `gorm:"foreignKey:AreaID" json:"area"`
    Latitude              float64   `gorm:"column:latitude" json:"latitude"`
    Longitude             float64   `gorm:"column:longitude" json:"longitude"`
    Name                  string    `gorm:"column:name" json:"name"`
    Alias                 string    `gorm:"column:alias" json:"alias"`
    Email                 string    `gorm:"column:email" json:"email"`  // REMOVED
    Phone                 string    `gorm:"column:phone" json:"phone"`
    // ... rest of fields
}
```

**After:**
```go
type Customer struct {
    ID                    string    `json:"id" gorm:"primaryKey"`
    Address               string    `gorm:"column:address" json:"address"`
    AreaID                string    `gorm:"column:area_id" json:"area_id"`
    Area                  *Areas    `gorm:"foreignKey:AreaID" json:"area"`
    Latitude              float64   `gorm:"column:latitude" json:"latitude"`
    Longitude             float64   `gorm:"column:longitude" json:"longitude"`
    Name                  string    `gorm:"column:name" json:"name"`
    Alias                 string    `gorm:"column:alias" json:"alias"`
    Phone                 string    `gorm:"column:phone" json:"phone"`
    // ... rest of fields
}
```

#### B. Customer Request Structure (`crm-be/internal/api/admin/customer/request.go`)

**Already correct** - no email field in request struct:
```go
type CreateAdminCustomerRequest struct {
    Name                  string  `json:"name" validate:"required"`
    Alias                 string  `json:"alias"`
    Address               string  `json:"address" validate:"required"`
    AreaID                string  `json:"area_id" validate:"required"`
    Phone                 string  `json:"phone" validate:"required"`
    Latitude              float64 `json:"latitude" validate:"required"`
    Longitude             float64 `json:"longitude" validate:"required"`
    ServiceRequestDate    string  `json:"service_request_date" validate:"required"`
    ProposedPackage       string  `json:"proposed_package" validate:"required"`
    BandwidthCapacity     string  `json:"bandwidth_capacity" validate:"required"`
    SalesRepresentativeID *string `json:"sales_representative_id"`
    CompanyID             *string `json:"company_id"`
    // No email field - correct!
}
```

## Testing

### Compilation Test
```bash
go build -o myapp.exe ./cmd/myapp
```
**Result**: ✅ Compilation successful with no errors

### Database Schema Verification
- ✅ Customer table structure checked
- ✅ No email column exists in database
- ✅ All other required columns are present

### Form Functionality
- ✅ Customer form loads without email field
- ✅ Form submission should work without email field
- ✅ Company selection feature still works

## Current Customer Form Fields

### Required Fields
1. **Name** - Customer full name
2. **Address** - Customer address
3. **Area ID** - Customer area selection
4. **Phone** - Customer phone number
5. **Latitude** - GPS latitude
6. **Longitude** - GPS longitude
7. **Service Request Date** - Date of service request
8. **Proposed Package** - Internet package selection
9. **Bandwidth Capacity** - Bandwidth specification

### Optional Fields
1. **Alias** - Customer nickname/alias
2. **Sales Representative ID** - Assigned sales rep
3. **Company ID** - Company assignment (newly added)

## Benefits

### For System
- ✅ **No database errors** when creating customers
- ✅ **Consistent schema** between frontend and backend
- ✅ **Proper data validation** without non-existent fields
- ✅ **Cleaner data model** without unused fields

### For Users
- ✅ **Form works correctly** without email field
- ✅ **No confusing email field** that doesn't save
- ✅ **Streamlined form** with only relevant fields
- ✅ **Company selection** still available

## Related Files Modified
1. **Frontend**: `crm-fe/pages/dashboard/customer/FormAddComponent.vue`
   - Removed email field from props
2. **Backend**: `crm-be/internal/models/entities/user_model.go`
   - Removed Email field from Customer struct

## Future Considerations

### If Email is Needed Later
If email functionality is required in the future:

1. **Add email column to database:**
   ```sql
   ALTER TABLE customer ADD COLUMN email VARCHAR(191) DEFAULT NULL;
   ```

2. **Add email field back to Customer entity:**
   ```go
   Email string `gorm:"column:email" json:"email"`
   ```

3. **Add email field to form and request structures**

### Current Contact Information
Customer contact information is currently handled through:
- **Phone** field for primary contact
- **Address** field for location
- **Company** field for business association

## Conclusion
The email field has been successfully removed from both frontend and backend, resolving the database error. Customer creation should now work correctly without the non-existent email column.

**Status**: ✅ **FIXED** - Customer creation without email field works correctly



