@echo off
echo ========================================
echo Auto-Populate Invoice Items Setup
echo ========================================
echo.

REM Check if Go is installed
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Go is not installed or not in PATH
    echo Please install Go first: https://golang.org/dl/
    pause
    exit /b 1
)

echo ✅ Go is installed
echo.

REM Set environment variables
echo Setting up environment variables...
set DB_HOST=localhost
set DB_PORT=3306
set DB_USER=root
set DB_PASSWORD=
set DB_NAME=iqgncnzy_skripsi

echo ✅ Environment variables set
echo.

REM Check if migration file exists
if not exist "database-migration-customer-multiple-products.sql" (
    echo ❌ Migration file not found: database-migration-customer-multiple-products.sql
    pause
    exit /b 1
)

echo ✅ Migration file found
echo.

REM Run migration
echo 🚀 Starting database migration...
echo This will add product_id to network_devices table
echo.
go run migrate-customer-multiple-products.go

if %errorlevel% equ 0 (
    echo.
    echo ✅ Database migration completed successfully!
    echo.
    echo 📋 Next steps:
    echo 1. Restart your backend server
    echo 2. Test the auto-populate feature in recurring invoices
    echo 3. Assign products to network devices if needed
    echo 4. Create a recurring invoice to test the feature
    echo.
    echo 🎯 How to use:
    echo 1. Go to Recurring Invoices page
    echo 2. Click "Create Recurring Invoice"
    echo 3. Select a customer
    echo 4. Invoice items will auto-populate from their network devices
    echo 5. Use "Auto-fill from Devices" button to refresh items
) else (
    echo.
    echo ❌ Migration failed!
    echo Please check the error messages above.
)

echo.
pause
