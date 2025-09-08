# PDF View Tracking Update

## Overview
Implemented PDF view tracking system to prevent duplicate payments by restricting PDF access to only once per invoice. This security feature ensures that once an invoice PDF is viewed, it cannot be accessed again to prevent employees from printing multiple copies and creating duplicate payments.

## Changes Made

### Backend Changes

#### 1. **Database Model Update** (`crm-be/internal/models/entities/user_model.go`)
- **Added**: `PdfViewed` boolean field to track if PDF has been viewed
- **Added**: `PdfViewedAt` timestamp field to record when PDF was viewed

```go
type Invoice struct {
    // ... existing fields ...
    PdfViewed    bool       `gorm:"column:pdf_viewed;type:boolean;default:false" json:"pdf_viewed"`
    PdfViewedAt  *time.Time `gorm:"column:pdf_viewed_at;type:timestamp" json:"pdf_viewed_at"`
    // ... rest of fields ...
}
```

#### 2. **New API Endpoint** (`crm-be/internal/api/admin/invoice/`)
- **Handler**: `MarkPdfViewedHandler` - Marks PDF as viewed
- **Service**: `MarkPdfViewedService` - Business logic for PDF tracking
- **Repository**: `MarkPdfViewedRepository` - Database operations
- **Route**: `POST /api/admin/invoice/:id/mark-pdf-viewed`

```go
func (h AdminInvoiceHandlerStruct) MarkPdfViewedHandler(c *fiber.Ctx) error {
    request := &IdAdminInvoiceRequest{}
    request.Id = c.Params("id")
    
    invoice, err := h.service.MarkPdfViewedService(*request)
    if err != nil {
        return helpers.ResponseUtils(c, fiber.StatusBadRequest, false, err.Error(), nil)
    }
    
    return helpers.ResponseUtils(c, fiber.StatusOK, true, "PDF marked as viewed", invoice)
}
```

#### 3. **Security Logic** (`crm-be/internal/api/admin/invoice/repository.go`)
- **Prevention**: Checks if PDF has already been viewed
- **Error Handling**: Returns error if PDF already viewed
- **Timestamp**: Records exact time when PDF was viewed

```go
func (r AdminInvoiceRepositoryStruct) MarkPdfViewedRepository(request IdAdminInvoiceRequest) (entities.Invoice, error) {
    // Check if PDF has already been viewed
    if invoice.PdfViewed {
        return invoice, fmt.Errorf("PDF has already been viewed and cannot be accessed again")
    }
    
    // Mark PDF as viewed with current timestamp
    now := time.Now()
    err = r.db.Model(&invoice).Updates(map[string]interface{}{
        "pdf_viewed":    true,
        "pdf_viewed_at": &now,
    }).Error
}
```

### Frontend Changes

#### 1. **API Integration** (`crm-fe/api/admin/invoice.ts`)
- **Added**: `markPdfViewed` method to call backend endpoint

```typescript
markPdfViewed: async (invoiceId: string) => {
  const response = await fetch(`${api}/api/admin/invoice/${invoiceId}/mark-pdf-viewed`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${useCookie("token").value}`,
    },
  });
  return response.json();
}
```

#### 2. **PDF View Handler** (`crm-fe/pages/dashboard/invoice/index.vue`)
- **Security Check**: Prevents viewing if PDF already viewed
- **API Call**: Marks PDF as viewed before opening
- **User Feedback**: Shows appropriate messages
- **State Management**: Updates local invoice data
- **Auto-Open Support**: Different messages for manual vs automatic opening

```typescript
async function handlePdfView(invoiceId: string, isAutoOpen: boolean = false) {
  try {
    // Check if PDF has already been viewed
    const invoice = customer.value.find(inv => inv.id === invoiceId);
    if (invoice?.pdf_viewed) {
      useToast().add({
        title: 'PDF Sudah Dilihat',
        description: 'PDF invoice ini sudah pernah dilihat dan tidak dapat diakses lagi untuk mencegah duplikasi pembayaran.',
        color: 'red'
      });
      return;
    }

    // Mark PDF as viewed before opening
    await invoiceAdminApi().markPdfViewed(invoiceId);
    
    // Navigate to PDF view
    navigateTo(`/invoice/${invoiceId}`);
    
    // Different messages for manual vs auto open
    if (isAutoOpen) {
      useToast().add({
        title: 'Status Diubah ke PAID',
        description: 'PDF invoice dibuka otomatis. PDF ini tidak dapat dibuka lagi untuk mencegah duplikasi pembayaran.',
        color: 'green'
      });
    } else {
      useToast().add({
        title: 'PDF Dibuka',
        description: 'PDF invoice telah dibuka. PDF ini tidak dapat dibuka lagi untuk mencegah duplikasi pembayaran.',
        color: 'yellow'
      });
    }
  } catch (error) {
    // Error handling
  }
}
```

#### 3. **Auto-Open PDF on Status Change**
- **Automatic Opening**: PDF opens automatically when status changes to "paid"
- **User Notification**: Clear message in confirmation modal
- **Seamless Experience**: No need for manual PDF download

```typescript
async function confirmStatusChange() {
  if (statusConfirmationData.value) {
    try {
      await proceedWithStatusUpdate(
        statusConfirmationData.value.invoiceId,
        statusConfirmationData.value.newStatus,
        statusConfirmationData.value.currentStatus
      );
      
      // If status was successfully changed to 'paid', automatically open PDF
      if (statusConfirmationData.value.newStatus === 'paid') {
        // Small delay to ensure status update is reflected in UI
        setTimeout(() => {
          handlePdfView(statusConfirmationData.value!.invoiceId, true);
        }, 500);
      }
    } catch (error) {
      console.error('Error updating status:', error);
    }
  }
  closeStatusConfirmationModal();
}
```

#### 4. **Visual Indicators**
- **Dropdown Menu**: Shows "PDF Sudah Dilihat" for viewed PDFs
- **Status Column**: Red indicator with "PDF Sudah Dilihat" text
- **Disabled State**: PDF button becomes disabled after viewing
- **Confirmation Modal**: Shows hint that PDF will open automatically

```vue
<!-- Dropdown Menu -->
{
  label: isPdfViewed(row.id) ? "PDF Sudah Dilihat" : "Download PDF",
  icon: isPdfViewed(row.id) ? "i-heroicons-eye-slash-20-solid" : "i-heroicons-arrow-down-on-square-20-solid",
  disabled: isPdfViewed(row.id),
  click: () => handlePdfView(row.id),
}

<!-- Status Column Indicator -->
<div v-if="isPdfViewed(row.id)" class="flex items-center text-xs text-red-600">
  <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
    <path fill-rule="evenodd" d="M13.477 14.89A6 6 0 015.11 6.524l8.367 8.368zm1.414-1.414L6.524 5.11a6 6 0 018.367 8.367zM18 10a8 8 0 11-16 0 8 8 0 0116 0z" clip-rule="evenodd"></path>
  </svg>
  PDF Sudah Dilihat
</div>
```

## Security Features

### 1. **One-Time Access**
- PDF can only be viewed once per invoice
- Subsequent attempts are blocked with error message
- Timestamp recorded for audit trail

### 2. **Visual Feedback**
- Clear indicators when PDF has been viewed
- Disabled buttons prevent accidental clicks
- Warning messages explain the restriction

### 3. **Audit Trail**
- `pdf_viewed_at` timestamp tracks when PDF was accessed
- Database record of all PDF view events
- Prevents duplicate payment attempts

## User Experience

### 1. **Manual PDF Viewing**
- User clicks "Download PDF"
- PDF opens normally
- Toast notification warns about one-time access
- PDF is marked as viewed in database

### 2. **Auto PDF Opening (Status Change to PAID)**
- User changes status to "paid" and confirms
- Status is updated successfully
- PDF opens automatically after 500ms delay
- Toast notification: "Status Diubah ke PAID - PDF invoice dibuka otomatis"
- PDF is marked as viewed in database
- User doesn't need to manually click "Download PDF"

### 3. **Subsequent Attempts**
- Button shows "PDF Sudah Dilihat" and is disabled
- Red indicator appears in status column
- Error message explains the restriction
- No PDF access allowed

### 4. **Visual Indicators**
- **Green**: PDF can be viewed
- **Red**: PDF has been viewed (restricted)
- **Disabled**: Button is not clickable
- **Icon Change**: Eye-slash icon for viewed PDFs

## Database Migration Required

To implement this feature, the following database migration is needed:

```sql
ALTER TABLE invoices 
ADD COLUMN pdf_viewed BOOLEAN DEFAULT FALSE,
ADD COLUMN pdf_viewed_at TIMESTAMP NULL;
```

## Benefits

1. **Prevents Duplicate Payments**: Employees cannot print multiple copies
2. **Audit Trail**: Complete record of PDF access
3. **Security**: One-time access prevents fraud
4. **User Awareness**: Clear visual indicators
5. **Compliance**: Helps with financial audit requirements

## Testing

### Test Cases
1. **First PDF View**: Should work normally and mark as viewed
2. **Second PDF View**: Should be blocked with error message
3. **Visual Indicators**: Should show correct status
4. **Database Updates**: Should record timestamp correctly
5. **Error Handling**: Should handle API failures gracefully

### Test Scenarios
- View PDF for first time ✅
- Try to view same PDF again ❌
- Check visual indicators ✅
- Verify database records ✅
- Test error handling ✅
