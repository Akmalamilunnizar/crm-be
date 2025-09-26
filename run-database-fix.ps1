# PowerShell script to fix database and test installation reports

Write-Host "Running database fix for installation reports..." -ForegroundColor Green

# Try to find MySQL executable
$mysqlPath = "C:\laragon\bin\mysql\mysql-8.0.30-winx64\bin\mysql.exe"
if (-not (Test-Path $mysqlPath)) {
    $mysqlPath = "mysql"  # Try system PATH
}

Write-Host "Using MySQL at: $mysqlPath" -ForegroundColor Yellow

try {
    # Run the database fix script
    Write-Host "Running database fix script..." -ForegroundColor Blue
    & $mysqlPath -u root -p iqgncnzy_skripsi -e "source fix-installation-report-view.sql"
    
    Write-Host "Database fix completed. Now running debug script..." -ForegroundColor Green
    
    # Run the debug script
    & $mysqlPath -u root -p iqgncnzy_skripsi -e "source debug-installation-reports.sql"
    
    Write-Host "Debug completed. Check the output above for any issues." -ForegroundColor Green
} catch {
    Write-Host "Error running database scripts: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Please run the SQL scripts manually in your MySQL client." -ForegroundColor Yellow
}

Read-Host "Press Enter to continue"

