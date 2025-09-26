# Test document photo in database

Write-Host "=== Testing Document Photo in Database ===" -ForegroundColor Green

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
    Write-Host "`nStep 1: Checking customer_installations table structure..." -ForegroundColor Blue
    
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
    
    Write-Host "`nStep 2: Checking recent installations..." -ForegroundColor Blue
    
    $recentQuery = @"
SELECT 
    id,
    customer_id,
    technician_id,
    document_type,
    document_photo,
    installation_type,
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
    
    Write-Host "`nStep 3: Checking document photo statistics..." -ForegroundColor Blue
    
    $statsQuery = @"
SELECT 
    COUNT(*) as total_installations,
    COUNT(document_photo) as installations_with_document_photo,
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
    
    Write-Host "`nStep 4: Checking installations with document photos..." -ForegroundColor Blue
    
    $photosQuery = @"
SELECT 
    id,
    document_type,
    document_photo,
    LENGTH(document_photo) as photo_path_length,
    createdAt
FROM customer_installations 
WHERE document_photo IS NOT NULL 
AND document_photo != ''
ORDER BY createdAt DESC;
"@
    
    $photosResult = & $mysqlPath -u root -p iqgncnzy_skripsi -e $photosQuery 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Installations with document photos retrieved successfully" -ForegroundColor Green
        Write-Host "Installations with Document Photos:" -ForegroundColor Yellow
        Write-Host $photosResult -ForegroundColor White
    } else {
        Write-Host "❌ Failed to get installations with document photos:" -ForegroundColor Red
        Write-Host $photosResult -ForegroundColor Red
    }
    
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n=== Next Steps ===" -ForegroundColor Green
Write-Host "1. Check if document_photo column exists in customer_installations table" -ForegroundColor Yellow
Write-Host "2. Verify document photos are being saved to database" -ForegroundColor Yellow
Write-Host "3. Check if files are being uploaded to uploads/installations/documents/" -ForegroundColor Yellow
Write-Host "4. Test the form submission again" -ForegroundColor Yellow

Read-Host "Press Enter to continue"
