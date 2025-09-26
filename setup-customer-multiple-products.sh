#!/bin/bash

echo "========================================"
echo "Customer Multiple Products Migration"
echo "========================================"
echo

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed or not in PATH"
    echo "Please install Go first: https://golang.org/dl/"
    exit 1
fi

echo "✅ Go is installed"
echo

# Set environment variables
echo "Setting up environment variables..."
export DB_HOST="localhost"
export DB_PORT="3306"
export DB_USER="root"
export DB_PASSWORD=""
export DB_NAME="iqgncnzy_skripsi"

echo "✅ Environment variables set"
echo

# Check if migration file exists
if [ ! -f "database-migration-customer-multiple-products.sql" ]; then
    echo "❌ Migration file not found: database-migration-customer-multiple-products.sql"
    exit 1
fi

echo "✅ Migration file found"
echo

# Run migration
echo "🚀 Starting migration..."
echo
go run migrate-customer-multiple-products.go

if [ $? -eq 0 ]; then
    echo
    echo "✅ Migration completed successfully!"
    echo
    echo "📋 Next steps:"
    echo "1. Verify the migration worked correctly"
    echo "2. Test your application with the new schema"
    echo "3. Update your application code"
    echo "4. Uncomment the final DROP statements in the migration file"
    echo "5. Run the migration again to remove the old product_id column"
else
    echo
    echo "❌ Migration failed!"
    echo "Please check the error messages above."
    exit 1
fi

echo
