# Quick fix for installation reports

Write-Host "=== Installation Reports Quick Fix ===" -ForegroundColor Green

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
    Write-Host "Running database fix..." -ForegroundColor Blue
    
    # Run the fix script
    $result = & $mysqlPath -u root -p iqgncnzy_skripsi -e "source fix-and-test.sql" 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Database fix completed successfully!" -ForegroundColor Green
        Write-Host "Output:" -ForegroundColor Yellow
        Write-Host $result -ForegroundColor White
    } else {
        Write-Host "Error running database fix:" -ForegroundColor Red
        Write-Host $result -ForegroundColor Red
    }
    
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n=== Next Steps ===" -ForegroundColor Green
Write-Host "1. Start backend: cd crm-be && go run cmd/myapp/main.go" -ForegroundColor Yellow
Write-Host "2. Refresh frontend page" -ForegroundColor Yellow
Write-Host "3. Check browser console for any remaining errors" -ForegroundColor Yellow

Read-Host "Press Enter to continue"
