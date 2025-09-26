# Quick test without database access

Write-Host "=== Quick Upload Test ===" -ForegroundColor Green

# Step 1: Check if backend is running
Write-Host "`nStep 1: Checking if backend is running..." -ForegroundColor Blue
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080" -Method GET -TimeoutSec 5
    Write-Host "Backend is running" -ForegroundColor Green
} catch {
    Write-Host "Backend is not running. Please start with: go run cmd/myapp/main.go" -ForegroundColor Red
    Write-Host "Starting backend..." -ForegroundColor Yellow
    Start-Process -FilePath "go" -ArgumentList "run", "cmd/myapp/main.go" -WindowStyle Normal
    Write-Host "Backend started. Please wait a few seconds and try again." -ForegroundColor Yellow
    exit 1
}

# Step 2: Create test image
Write-Host "`nStep 2: Creating test image..." -ForegroundColor Blue
$testImagePath = "test-ktp.png"
$pngData = @"
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
"@

try {
    $bytes = [System.Convert]::FromBase64String($pngData)
    [System.IO.File]::WriteAllBytes($testImagePath, $bytes)
    Write-Host "Test image created: $testImagePath" -ForegroundColor Green
} catch {
    Write-Host "Failed to create test image: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Step 3: Check uploads directory
Write-Host "`nStep 3: Checking uploads directory..." -ForegroundColor Blue
$uploadsDir = "uploads/installations/documents"
if (Test-Path $uploadsDir) {
    Write-Host "Uploads directory exists: $uploadsDir" -ForegroundColor Green
    $files = Get-ChildItem -Path $uploadsDir -File -ErrorAction SilentlyContinue
    Write-Host "Found $($files.Count) files in uploads directory" -ForegroundColor Yellow
    if ($files.Count -gt 0) {
        Write-Host "Recent files:" -ForegroundColor Yellow
        $files | Sort-Object LastWriteTime -Descending | Select-Object -First 3 | ForEach-Object {
            Write-Host "  - $($_.Name) ($($_.LastWriteTime))" -ForegroundColor White
        }
    }
} else {
    Write-Host "Uploads directory does not exist: $uploadsDir" -ForegroundColor Yellow
    Write-Host "It will be created automatically when first file is uploaded" -ForegroundColor White
}

# Step 4: Show curl command
Write-Host "`nStep 4: Curl command for testing..." -ForegroundColor Blue
Write-Host "To test document upload, run this curl command:" -ForegroundColor Yellow
Write-Host "First, get your token from browser (F12 > Application > Cookies > token)" -ForegroundColor Yellow

$curlCommand = @"
curl -X POST "http://localhost:8080/api/admin/customer-installation/report-installations" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -F "customer_id=0ae7d407-2656-4fe3-878f-89c19abcbdac" \
  -F "technician_id=c13a6c87-ec28-47ba-84c2-58b5ace2af57" \
  -F "assets_id=5ca1606b-66b3-4958-b7af-f48d4cda800a" \
  -F "document_type=KTP" \
  -F "document_photo=@$testImagePath" \
  -F "status=completed" \
  -F "installation_type=new_installation" \
  -F "switch_id=SW-JBR-001" \
  -F "port_number=1" \
  -F "mac_address=00:11:12:13:14:15" \
  -F "ip_static=192.168.100.11" \
  -F "cable_type=UTP Cat6" \
  -F "cable_length=25" \
  -F "user_login=admin@test.com" \
  -F "password=password123" \
  -F "user_status=Active"
"@

Write-Host $curlCommand -ForegroundColor White

# Step 5: Show what to check in backend logs
Write-Host "`nStep 5: Backend logs to watch..." -ForegroundColor Blue
Write-Host "After running the curl command, check backend console for these logs:" -ForegroundColor Yellow
Write-Host "1. Document photo upload started - filename: test-ktp.png, size: 95" -ForegroundColor White
Write-Host "2. Document photo uploaded successfully - filePath: uploads/installations/documents/document_YYYYMMDD_HHMMSS.png" -ForegroundColor White
Write-Host "3. Installation report request data - customer_id: ..., document_photo: uploads/installations/documents/document_YYYYMMDD_HHMMSS.png" -ForegroundColor White
Write-Host "4. Creating installation record - DocumentPhoto: uploads/installations/documents/document_YYYYMMDD_HHMMSS.png, DocumentType: KTP" -ForegroundColor White
Write-Host "5. Installation record created successfully with ID: ..." -ForegroundColor White

# Step 6: Show database check
Write-Host "`nStep 6: Database check..." -ForegroundColor Blue
Write-Host "After successful upload, check database with:" -ForegroundColor Yellow
Write-Host "SELECT id, document_type, document_photo FROM customer_installations ORDER BY createdAt DESC LIMIT 1;" -ForegroundColor White

# Clean up test file
if (Test-Path $testImagePath) {
    Remove-Item $testImagePath
    Write-Host "`nTest image cleaned up" -ForegroundColor Green
}

Write-Host "`n=== Next Steps ===" -ForegroundColor Green
Write-Host "1. Get your authentication token from browser" -ForegroundColor Yellow
Write-Host "2. Replace YOUR_TOKEN_HERE in the curl command above" -ForegroundColor Yellow
Write-Host "3. Run the curl command" -ForegroundColor Yellow
Write-Host "4. Watch backend console for the 5 log messages" -ForegroundColor Yellow
Write-Host "5. Check database for document_photo value" -ForegroundColor Yellow

Read-Host "Press Enter to continue"
