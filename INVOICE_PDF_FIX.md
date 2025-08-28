# Invoice PDF Amount Fix Documentation

## 🔍 Masalah yang Ditemukan

**Problem**: PDF invoice yang di-download selalu menampilkan nilai `Rp 0` untuk:
- Invoice Total
- Total Paid  
- Total/Remaining Amount

**Root Cause**:
1. **Backend**: Field `Amount` di entity Invoice hanya menyimpan total dari item terakhir, bukan total keseluruhan
2. **Frontend**: Menggunakan `props.invoice.total` yang hardcoded, bukan data dari API `invoiceDetail`

## 🛠️ Solusi yang Diterapkan

### 1. Backend Fixes

#### A. Repository Invoice (`crm-be/internal/api/admin/invoice/repository.go`)

**Before (Buggy)**:
```go
for _, invoiceItem := range invoice.InvoiceItems {
    invoiceItem.InvoiceID = invoice.ID
    invoiceItem.Total = invoiceItem.Price * invoiceItem.Qty
    invoice.Amount = invoiceItem.Total  // ❌ Hanya item terakhir
}
```

**After (Fixed)**:
```go
// Calculate total amount from all invoice items
var totalAmount int64 = 0
for _, invoiceItem := range invoice.InvoiceItems {
    invoiceItem.InvoiceID = invoice.ID
    invoiceItem.Total = invoiceItem.Price * invoiceItem.Qty
    totalAmount += invoiceItem.Total
}

// Set total amount to sum of all items
invoice.Amount = totalAmount
```

#### B. Enhanced FindByIdAdminInvoiceRepository

**Added auto-recalculation**:
```go
// Calculate total amount from invoice items if not set or incorrect
if len(invoice.InvoiceItems) > 0 {
    var totalAmount int64 = 0
    for _, item := range invoice.InvoiceItems {
        totalAmount += item.Total
    }
    
    // Update invoice amount if it's different from calculated total
    if invoice.Amount != totalAmount {
        invoice.Amount = totalAmount
        // Update in database
        r.db.Model(&invoice).Update("amount", totalAmount)
    }
}
```

### 2. Frontend Fixes

#### A. Invoice Detail Page (`crm-fe/pages/invoice/[id].vue`)

**Before (Hardcoded)**:
```vue
<p>Invoice Total: Rp {{ props.invoice.total.toLocaleString("id-ID") }}</p>
<p>Total Paid: Rp {{ props.invoice.totalPaid.toLocaleString("id-ID") }}</p>
```

**After (Dynamic from API)**:
```vue
<p>Invoice Total: {{ formatIDR(invoiceTotal) }}</p>
<p>Total Paid: {{ formatIDR(totalPaid) }}</p>
```

#### B. Computed Properties

**Added smart calculations**:
```javascript
const invoiceTotal = computed(() => {
  if (!invoiceDetail.value.invoice_items || invoiceDetail.value.invoice_items.length === 0) {
    return invoiceDetail.value.amount || 0;
  }
  return invoiceDetail.value.invoice_items.reduce((sum: number, item: any) => sum + (item.total || 0), 0);
});

const totalPaid = computed(() => {
  if (!invoiceDetail.value.transaction) return 0;
  return invoiceDetail.value.transaction.amount || 0;
});

const remainingAmount = computed(() => {
  return invoiceTotal.value - totalPaid.value;
});
```

### 3. Data Migration Script

**Created**: `crm-be/fix-invoice-amounts.go`

Script ini memperbaiki data invoice yang sudah ada dengan:
- Menghitung ulang total amount dari invoice items
- Update database jika ada perbedaan
- Log detail perbaikan

## 🚀 Cara Test

### 1. Backend Test
```bash
cd crm-be
go run cmd/myapp/main.go
```

### 2. Frontend Test  
```bash
cd crm-fe
npm run dev
```

### 3. Test API
```bash
# Test invoice endpoint
curl http://localhost:3001/api/admin/invoice/1
```

### 4. Test PDF Generation
- Buka halaman invoice detail
- Klik "Download PDF"
- Verifikasi nilai total sudah benar (bukan Rp 0)

## 📊 Expected Results

**Before Fix**:
- Invoice Total: Rp 0
- Total Paid: Rp 0  
- Amount Due: Rp 0

**After Fix**:
- Invoice Total: Rp 100.000 (sesuai data)
- Total Paid: Rp 50.000 (sesuai transaction)
- Amount Due: Rp 50.000 (calculated)

## 🔧 Maintenance

### Auto-fix untuk Invoice Baru
- Setiap invoice baru akan otomatis menghitung total dari items
- Field `Amount` akan selalu konsisten dengan sum dari `InvoiceItems.Total`

### Auto-fix untuk Invoice Existing
- Setiap kali invoice di-fetch, total akan di-recalculate otomatis
- Database akan di-update jika ada perbedaan

## 📝 Notes

1. **Performance**: Recalculation hanya dilakukan saat diperlukan
2. **Data Integrity**: Total amount selalu konsisten dengan invoice items
3. **Backward Compatibility**: Tidak ada breaking changes pada API
4. **Error Handling**: Graceful fallback jika data tidak lengkap

## 🎯 Next Steps

1. **Monitor**: Pastikan PDF generation berfungsi dengan benar
2. **Test**: Test dengan berbagai skenario invoice
3. **Optimize**: Jika diperlukan, tambahkan caching untuk total calculation
4. **Document**: Update API documentation jika diperlukan
