# 🔄 Recurring Invoices Feature

This feature allows you to create and manage recurring invoices that automatically generate regular invoices for customers.

## 🚀 Quick Start

After pulling the latest code from GitHub, run one of these commands to set up the recurring invoices feature:

### Windows
```cmd
setup-recurring-invoices.bat
```

### Linux/Mac
```bash
chmod +x setup-recurring-invoices.sh
./setup-recurring-invoices.sh
```

## 📁 Files Added

### Backend Files
- `internal/api/admin/recurring_invoice/` - Complete API implementation
- `internal/models/entities/recurring_invoice_model.go` - Data models
- `database-migration-recurring-invoices.sql` - Main table
- `database-migration-recurring-invoice-history.sql` - History table
- `database-migration-recurring-invoice-indexes*.sql` - Database indexes
- `migrate-recurring-invoices.go` - Migration runner
- `seed-recurring-invoices.go` - Data seeder
- `clear-recurring-invoices.go` - Data cleaner

### Frontend Files
- `pages/dashboard/recurring-invoice/index.vue` - Main list page
- `pages/dashboard/recurring-invoice/FormAddRecurringInvoice.vue` - Add/Edit form
- `api/admin/recurring-invoice.ts` - API client
- Updated `utilities/rolePermissions.ts` - Added menu item

## 🎯 Features

- ✅ **CRUD Operations**: Create, Read, Update, Delete recurring invoices
- ✅ **Multiple Frequencies**: Monthly, Quarterly, Yearly
- ✅ **Status Management**: Active, Stopped, Completed
- ✅ **Invoice Items**: JSON-based item management
- ✅ **History Tracking**: Track generated invoices
- ✅ **Filtering & Pagination**: Advanced list management
- ✅ **Responsive Design**: Works on all devices

## 📊 Sample Data

The seeder creates 5 sample recurring invoices:
1. Monthly internet service (Rp 5,000,000)
2. Monthly hosting & maintenance (Rp 2,500,000)
3. Quarterly software license (Rp 1,000,000)
4. Monthly cloud services (Rp 7,500,000)
5. Monthly digital marketing (Rp 3,000,000) - Stopped

## 🔧 Manual Setup

If the quick setup doesn't work:

1. **Run Migrations**
   ```bash
   go run migrate-recurring-invoices.go
   ```

2. **Seed Data**
   ```bash
   go run seed-recurring-invoices.go
   ```

3. **Start Backend**
   ```bash
   go run cmd/myapp/main.go
   ```

4. **Start Frontend**
   ```bash
   cd ../crm-fe
   npm run dev
   ```

## 🌐 Access

Navigate to: `http://localhost:3000/dashboard/recurring-invoice`

## 🛠️ Troubleshooting

### Common Issues

**"No customers found"**
- Run: `go run database-seeder.go`

**"No users found"**
- Run: `go run setup-roles.go`

**"Foreign key constraint fails"**
- Ensure all tables exist and have data
- Check database connection

**"404 Not Found"**
- Ensure backend is running on port 3001
- Check API routes are registered

### Reseeding Data

To clear and reseed:
```bash
go run clear-recurring-invoices.go
go run seed-recurring-invoices.go
```

## 📚 Documentation

- `RECURRING_INVOICE_SEEDER.md` - Detailed seeder documentation
- `TECHNICIAN_NOTE_FEATURE.md` - Related feature docs

## 🎉 Success!

Once set up, you'll have a fully functional recurring invoices system that matches the reference implementation at `https://www.rndbill.lilly.net.id/?ng=invoices/list-recurring`.

Happy coding! 🚀
