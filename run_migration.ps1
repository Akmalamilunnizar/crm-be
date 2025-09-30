# PowerShell script to run database migration
# This script will execute the migration_script.sql file

Write-Host "Starting database migration..." -ForegroundColor Green

# Try to find MySQL executable in common Laragon locations
$mysqlPaths = @(
    "C:\laragon\bin\mysql\mysql-8.4.3-winx64\bin\mysql.exe",
    "C:\laragon\bin\mysql\mysql-8.0.30-winx64\bin\mysql.exe",
    "C:\laragon\bin\mysql\mysql-8.0.31-winx64\bin\mysql.exe",
    "C:\laragon\bin\mysql\mysql-8.0.32-winx64\bin\mysql.exe",
    "C:\laragon\bin\mysql\mysql-8.0.33-winx64\bin\mysql.exe",
    "C:\laragon\bin\mysql\mysql-8.0.34-winx64\bin\mysql.exe",
    "C:\laragon\bin\mysql\mysql-8.0.35-winx64\bin\mysql.exe",
    "C:\laragon\bin\mysql\mysql-8.0.36-winx64\bin\mysql.exe",
    "C:\laragon\bin\mysql\mysql-8.0.37-winx64\bin\mysql.exe",
    "C:\laragon\bin\mysql\mysql-8.0.38-winx64\bin\mysql.exe",
    "C:\laragon\bin\mysql\mysql-8.0.39-winx64\bin\mysql.exe",
    "C:\laragon\bin\mysql\mysql-8.0.40-winx64\bin\mysql.exe"
)

$mysqlExe = $null
foreach ($path in $mysqlPaths) {
    if (Test-Path $path) {
        $mysqlExe = $path
        Write-Host "Found MySQL at: $path" -ForegroundColor Yellow
        break
    }
}

if (-not $mysqlExe) {
    Write-Host "MySQL executable not found in common Laragon locations." -ForegroundColor Red
    Write-Host "Please check if Laragon is installed and MySQL is available." -ForegroundColor Red
    exit 1
}

# Check if migration script exists
if (-not (Test-Path "migration_script.sql")) {
    Write-Host "Migration script not found: migration_script.sql" -ForegroundColor Red
    exit 1
}

Write-Host "Executing migration script..." -ForegroundColor Yellow

# Execute the migration script
try {
    # Read SQL content and execute it
    $sqlContent = Get-Content "migration_script.sql" -Raw
    
    # Execute using cmd to handle redirection properly
    $cmdArgs = "/c `"$mysqlExe`" -u root iqgncnzy_skripsi < migration_script.sql"
    $result = cmd.exe $cmdArgs
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Migration completed successfully!" -ForegroundColor Green
        Write-Host "Database structure has been updated to match lilly.sql" -ForegroundColor Green
        Write-Host "Output: $result" -ForegroundColor Cyan
    } else {
        Write-Host "Migration failed with exit code: $LASTEXITCODE" -ForegroundColor Red
        Write-Host "Output: $result" -ForegroundColor Red
    }
} catch {
    Write-Host "Error executing migration: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "Migration process completed." -ForegroundColor Cyan
