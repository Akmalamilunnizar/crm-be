# Recurring Invoice Seeder

This script seeds the database with sample recurring invoice data for testing and development purposes.

## Prerequisites

1. **Database Setup**: Ensure your MySQL database is running and accessible
2. **Tables Created**: Run the recurring invoice migrations first:
   ```bash
   go run migrate-recurring-invoices.go
   ```
3. **Existing Data**: The seeder requires existing customers and users in the database

## Usage

### Quick Setup (Recommended)

**For Linux/Mac:**
```bash
# From the crm-be directory
chmod +x setup-recurring-invoices.sh
./setup-recurring-invoices.sh
```

**For Windows:**
```cmd
# From the crm-be directory
setup-recurring-invoices.bat
```

### Manual Setup

**Step 1: Run Migrations**
```bash
go run migrate-recurring-invoices.go
```

**Step 2: Run the Seeder**
```bash
go run seed-recurring-invoices.go
```

### What it does

- ✅ Connects to the database (`iqgncnzy_skripsi`)
- ✅ Checks for existing recurring invoices (skips if already seeded)
- ✅ Automatically uses the first available customer and user
- ✅ Creates 5 sample recurring invoices with different frequencies and statuses
- ✅ Provides detailed output of what was created

### Sample Data Created

1. **rec_001** - Monthly internet service subscription (Rp 5,000,000) - Active
2. **rec_002** - Monthly hosting and maintenance (Rp 2,500,000) - Active  
3. **rec_003** - Quarterly software license (Rp 1,000,000) - Active
4. **rec_004** - Monthly cloud services and support (Rp 7,500,000) - Active
5. **rec_005** - Monthly digital marketing services (Rp 3,000,000) - Stopped

### Features

- **Smart Detection**: Automatically detects existing customers and users
- **Duplicate Prevention**: Skips seeding if recurring invoices already exist
- **Flexible**: Works with any existing customer/user data
- **Safe**: Uses proper error handling and validation

### Troubleshooting

**Error: "No customers found"**
- Run the main database seeder first: `go run database-seeder.go`

**Error: "No users found"**  
- Run the user setup: `go run setup-roles.go`

**Error: "Foreign key constraint fails"**
- Ensure all required tables exist and have data
- Check that customer and user IDs are valid

### Reseeding

To reseed the data (clear and recreate):

**Option 1: Using the clear script (Recommended)**
```bash
# Clear existing data
go run clear-recurring-invoices.go

# Then run the seeder again
go run seed-recurring-invoices.go
```

**Option 2: Manual SQL**
```sql
-- Clear existing data
DELETE FROM recurring_invoice_history;
DELETE FROM recurring_invoices;

-- Then run the seeder again
go run seed-recurring-invoices.go
```

## Integration

After seeding, you can:

1. **View the data** in the frontend at `/dashboard/recurring-invoice`
2. **Test CRUD operations** (Create, Read, Update, Delete)
3. **Test filtering and pagination**
4. **Test invoice generation** functionality

## Notes

- The seeder uses the first available customer and user from your database
- All amounts are in Indonesian Rupiah (IDR)
- Dates are set to 2024 for testing purposes
- Invoice items are stored as JSON strings
