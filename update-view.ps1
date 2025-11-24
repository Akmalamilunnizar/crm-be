# Update Installation Report View
# This script updates the installation_report_complete view to include all necessary fields

Write-Host "Updating installation_report_complete view..." -ForegroundColor Cyan

# Read database credentials from .env file
$envFile = ".env"
if (Test-Path $envFile) {
    Get-Content $envFile | ForEach-Object {
        if ($_ -match '^\s*([^#][^=]+?)\s*=\s*(.+?)\s*$') {
            $name = $matches[1]
            $value = $matches[2]
            Set-Variable -Name $name -Value $value -Scope Script
        }
    }
}

# Extract database credentials
if ($DB_DSN -match 'mysql://([^:]+):([^@]+)@([^:]+):(\d+)/(.+)') {
    $dbUser = $matches[1]
    $dbPassword = $matches[2]
    $dbHost = $matches[3]
    $dbPort = $matches[4]
    $dbName = $matches[5]
} else {
    Write-Host "Error: Could not parse DB_DSN from .env file" -ForegroundColor Red
    exit 1
}

Write-Host "Database: $dbName" -ForegroundColor Yellow
Write-Host "Host: $dbHost" -ForegroundColor Yellow

# Execute the SQL file
$sqlFile = "update-installation-report-view.sql"
if (Test-Path $sqlFile) {
    Write-Host "Executing SQL file: $sqlFile" -ForegroundColor Green
    
    # Use mysql command to execute the SQL file
    $mysqlCmd = "mysql -h $dbHost -P $dbPort -u $dbUser -p$dbPassword $dbName"
    
    Get-Content $sqlFile | & mysql -h $dbHost -P $dbPort -u $dbUser "-p$dbPassword" $dbName
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "View updated successfully!" -ForegroundColor Green
    } else {
        Write-Host "Error updating view. Exit code: $LASTEXITCODE" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "Error: SQL file not found: $sqlFile" -ForegroundColor Red
    exit 1
}

Write-Host "`nDone!" -ForegroundColor Cyan
