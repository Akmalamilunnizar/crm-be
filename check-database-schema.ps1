# Check database schema for customer_installations table

Write-Host "=== Check Database Schema ===" -ForegroundColor Green

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
    Write-Host "`nStep 1: Checking table structure..." -ForegroundColor Blue
    
    $structureQuery = "DESCRIBE customer_installations;"
    $structureResult = & $mysqlPath -u root -p iqgncnzy_skripsi -e $structureQuery 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Table structure retrieved successfully" -ForegroundColor Green
        Write-Host "Table Structure:" -ForegroundColor Yellow
        Write-Host $structureResult -ForegroundColor White
    } else {
        Write-Host "❌ Failed to get table structure:" -ForegroundColor Red
        Write-Host $structureResult -ForegroundColor Red
    }
    
    Write-Host "`nStep 2: Checking document_photo column..." -ForegroundColor Blue
    
    $columnQuery = @"
SELECT 
    COLUMN_NAME,
    DATA_TYPE,
    IS_NULLABLE,
    COLUMN_DEFAULT,
    COLUMN_COMMENT
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_SCHEMA = 'iqgncnzy_skripsi' 
AND TABLE_NAME = 'customer_installations'
AND COLUMN_NAME = 'document_photo';
"@
    
    $columnResult = & $mysqlPath -u root -p iqgncnzy_skripsi -e $columnQuery 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Document photo column info retrieved successfully" -ForegroundColor Green
        Write-Host "Document Photo Column Info:" -ForegroundColor Yellow
        Write-Host $columnResult -ForegroundColor White
    } else {
        Write-Host "❌ Failed to get document photo column info:" -ForegroundColor Red
        Write-Host $columnResult -ForegroundColor Red
    }
    
    Write-Host "`nStep 3: Checking recent installations..." -ForegroundColor Blue
    
    $recentQuery = @"
SELECT 
    id,
    customer_id,
    technician_id,
    document_type,
    document_photo,
    status,
    createdAt,
    updatedAt
FROM customer_installations 
ORDER BY createdAt DESC 
LIMIT 5;
"@
    
    $recentResult = & $mysqlPath -u root -p iqgncnzy_skripsi -e $recentQuery 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Recent installations retrieved successfully" -ForegroundColor Green
        Write-Host "Recent Installations:" -ForegroundColor Yellow
        Write-Host $recentResult -ForegroundColor White
    } else {
        Write-Host "❌ Failed to get recent installations:" -ForegroundColor Red
        Write-Host $recentResult -ForegroundColor Red
    }
    
    Write-Host "`nStep 4: Checking document photo statistics..." -ForegroundColor Blue
    
    $statsQuery = @"
SELECT 
    COUNT(*) as total_records,
    COUNT(document_photo) as records_with_document_photo,
    COUNT(CASE WHEN document_photo IS NOT NULL AND document_photo != '' THEN 1 END) as non_empty_document_photos
FROM customer_installations;
"@
    
    $statsResult = & $mysqlPath -u root -p iqgncnzy_skripsi -e $statsQuery 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Document photo statistics retrieved successfully" -ForegroundColor Green
        Write-Host "Document Photo Statistics:" -ForegroundColor Yellow
        Write-Host $statsResult -ForegroundColor White
    } else {
        Write-Host "❌ Failed to get document photo statistics:" -ForegroundColor Red
        Write-Host $statsResult -ForegroundColor Red
    }
    
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n=== Analysis ===" -ForegroundColor Green
Write-Host "1. Check if document_photo column exists in the table" -ForegroundColor Yellow
Write-Host "2. Check if document_photo column allows NULL values" -ForegroundColor Yellow
Write-Host "3. Check if recent installations have document_photo values" -ForegroundColor Yellow
Write-Host "4. Check if document_photo statistics show any data" -ForegroundColor Yellow

Read-Host "Press Enter to continue"
