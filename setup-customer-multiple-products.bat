@echo off
echo ========================================
echo Customer Multiple Products Migration
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
echo 🚀 Starting migration...
echo.
go run migrate-customer-multiple-products.go

if %errorlevel% equ 0 (
    echo.
    echo ✅ Migration completed successfully!
    echo.
    echo 📋 Next steps:
    echo 1. Verify the migration worked correctly
    echo 2. Test your application with the new schema
    echo 3. Update your application code
    echo 4. Uncomment the final DROP statements in the migration file
    echo 5. Run the migration again to remove the old product_id column
) else (
    echo.
    echo ❌ Migration failed!
    echo Please check the error messages above.
)

echo.
pause
