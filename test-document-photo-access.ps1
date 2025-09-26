# Test document photo access

Write-Host "=== Testing Document Photo Access ===" -ForegroundColor Green

# Check if uploads directory exists
$uploadsDir = "uploads/installations/documents"
if (Test-Path $uploadsDir) {
    Write-Host "✅ Uploads directory exists: $uploadsDir" -ForegroundColor Green
    
    # List files in uploads directory
    $files = Get-ChildItem -Path $uploadsDir -File
    if ($files.Count -gt 0) {
        Write-Host "📁 Found $($files.Count) files in uploads directory:" -ForegroundColor Yellow
        foreach ($file in $files) {
            Write-Host "  - $($file.Name) ($($file.Length) bytes)" -ForegroundColor White
        }
    } else {
        Write-Host "⚠️  No files found in uploads directory" -ForegroundColor Yellow
    }
} else {
    Write-Host "❌ Uploads directory does not exist: $uploadsDir" -ForegroundColor Red
    Write-Host "Creating directory..." -ForegroundColor Yellow
    New-Item -ItemType Directory -Path $uploadsDir -Force
}

# Test static file serving
Write-Host "`n=== Testing Static File Serving ===" -ForegroundColor Blue

$testUrl = "http://localhost:8080/uploads/installations/documents/"
Write-Host "Testing URL: $testUrl" -ForegroundColor Yellow

try {
    $response = Invoke-WebRequest -Uri $testUrl -Method GET -TimeoutSec 5
    Write-Host "✅ Static file serving is working" -ForegroundColor Green
    Write-Host "Response status: $($response.StatusCode)" -ForegroundColor White
} catch {
    Write-Host "❌ Static file serving test failed: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Make sure backend is running: go run cmd/myapp/main.go" -ForegroundColor Yellow
}

# Check database for document_photo paths
Write-Host "`n=== Checking Database Document Photo Paths ===" -ForegroundColor Blue

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

if ($mysqlPath) {
    Write-Host "Using MySQL at: $mysqlPath" -ForegroundColor Yellow
    
    $query = @"
SELECT 
    id,
    document_type,
    document_photo,
    customer_id
FROM customer_installations 
WHERE document_photo IS NOT NULL 
AND document_photo != ''
ORDER BY createdAt DESC
LIMIT 5;
"@
    
    try {
        $result = & $mysqlPath -u root -p iqgncnzy_skripsi -e $query 2>&1
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ Database query successful" -ForegroundColor Green
            Write-Host "Document photo records:" -ForegroundColor Yellow
            Write-Host $result -ForegroundColor White
        } else {
            Write-Host "❌ Database query failed:" -ForegroundColor Red
            Write-Host $result -ForegroundColor Red
        }
    } catch {
        Write-Host "❌ Database connection failed: $($_.Exception.Message)" -ForegroundColor Red
    }
} else {
    Write-Host "❌ MySQL not found. Please install MySQL or update the path." -ForegroundColor Red
}

Write-Host "`n=== Next Steps ===" -ForegroundColor Green
Write-Host "1. Make sure backend is running: go run cmd/myapp/main.go" -ForegroundColor Yellow
Write-Host "2. Check if document photos are uploaded correctly" -ForegroundColor Yellow
Write-Host "3. Test frontend document photo display" -ForegroundColor Yellow
Write-Host "4. Verify static file serving is working" -ForegroundColor Yellow

Read-Host "Press Enter to continue"
