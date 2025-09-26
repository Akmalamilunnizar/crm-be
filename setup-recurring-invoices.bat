@echo off
echo 🚀 Setting up Recurring Invoices Feature...
echo ==========================================

REM Check if Go is installed
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Go is not installed. Please install Go first.
    pause
    exit /b 1
)

echo ✅ Go is installed

REM Check if MySQL is running (basic check)
mysql --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ MySQL is not accessible. Please ensure MySQL is running.
    pause
    exit /b 1
)

echo ✅ MySQL is accessible

REM Run migrations
echo 📊 Running database migrations...
go run migrate-recurring-invoices.go

if %errorlevel% neq 0 (
    echo ❌ Database migrations failed
    pause
    exit /b 1
)

echo ✅ Database migrations completed

REM Run seeder
echo 🌱 Seeding recurring invoice data...
go run seed-recurring-invoices.go

if %errorlevel% neq 0 (
    echo ❌ Recurring invoice seeding failed
    pause
    exit /b 1
)

echo ✅ Recurring invoice seeding completed

echo.
echo 🎉 Recurring Invoices feature setup completed!
echo 💡 You can now:
echo    1. Start your backend server
echo    2. Start your frontend server
echo    3. Navigate to /dashboard/recurring-invoice
echo.
echo 📚 For more information, see RECURRING_INVOICE_SEEDER.md
pause
