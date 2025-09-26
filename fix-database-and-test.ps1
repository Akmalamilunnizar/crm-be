# Fix database view and test installation reports

Write-Host "=== Installation Reports Database Fix ===" -ForegroundColor Green

# Try to find MySQL executable
$mysqlPaths = @(
    "C:\laragon\bin\mysql\mysql-8.0.30-winx64\bin\mysql.exe",
    "C:\laragon\bin\mysql\mysql-8.4.3-winx64\bin\mysql.exe",
    "mysql"
)

$mysqlPath = $null
foreach ($path in $mysqlPaths) {
    if ($path -eq "mysql") {
        $mysqlPath = "mysql"
        break
    } elseif (Test-Path $path) {
        $mysqlPath = $path
        break
    }
}

if (-not $mysqlPath) {
    Write-Host "MySQL not found. Please install MySQL or update the path in this script." -ForegroundColor Red
    exit 1
}

Write-Host "Using MySQL at: $mysqlPath" -ForegroundColor Yellow

try {
    Write-Host "Step 1: Fixing database view..." -ForegroundColor Blue
    
    # Run the fix script
    $result = & $mysqlPath -u root -p iqgncnzy_skripsi -e "source fix-and-test.sql" 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Database view fixed successfully!" -ForegroundColor Green
    } else {
        Write-Host "Error fixing database view:" -ForegroundColor Red
        Write-Host $result -ForegroundColor Red
    }
    
    Write-Host "`nStep 2: Testing database view..." -ForegroundColor Blue
    
    # Test the view
    $testResult = & $mysqlPath -u root -p iqgncnzy_skripsi -e "source test-database-view.sql" 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Database test completed!" -ForegroundColor Green
        Write-Host "Test Results:" -ForegroundColor Yellow
        Write-Host $testResult -ForegroundColor White
    } else {
        Write-Host "Error testing database:" -ForegroundColor Red
        Write-Host $testResult -ForegroundColor Red
    }
    
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n=== Next Steps ===" -ForegroundColor Green
Write-Host "1. Restart backend: cd crm-be && go run cmd/myapp/main.go" -ForegroundColor Yellow
Write-Host "2. Test API endpoint: http://localhost:8080/api/admin/customer-installation/report-complete" -ForegroundColor Yellow
Write-Host "3. Refresh frontend page" -ForegroundColor Yellow
Write-Host "4. Check browser console for API response" -ForegroundColor Yellow

Read-Host "Press Enter to continue"
