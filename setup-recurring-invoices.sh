#!/bin/bash

echo "🚀 Setting up Recurring Invoices Feature..."
echo "=========================================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

echo "✅ Go is installed"

# Check if MySQL is running
if ! mysqladmin ping -h localhost -u root --silent; then
    echo "❌ MySQL is not running. Please start MySQL first."
    exit 1
fi

echo "✅ MySQL is running"

# Run migrations
echo "📊 Running database migrations..."
go run migrate-recurring-invoices.go

if [ $? -eq 0 ]; then
    echo "✅ Database migrations completed"
else
    echo "❌ Database migrations failed"
    exit 1
fi

# Run seeder
echo "🌱 Seeding recurring invoice data..."
go run seed-recurring-invoices.go

if [ $? -eq 0 ]; then
    echo "✅ Recurring invoice seeding completed"
else
    echo "❌ Recurring invoice seeding failed"
    exit 1
fi

echo ""
echo "🎉 Recurring Invoices feature setup completed!"
echo "💡 You can now:"
echo "   1. Start your backend server"
echo "   2. Start your frontend server" 
echo "   3. Navigate to /dashboard/recurring-invoice"
echo ""
echo "📚 For more information, see RECURRING_INVOICE_SEEDER.md"
