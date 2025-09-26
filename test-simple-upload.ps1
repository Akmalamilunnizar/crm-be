# Simple test for document photo upload

Write-Host "=== Simple Document Photo Upload Test ===" -ForegroundColor Green

# Create a simple test image file (1x1 pixel PNG)
$testImagePath = "test-ktp.png"
$pngData = @"
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
"@

try {
    # Convert base64 to bytes and create test image
    $bytes = [System.Convert]::FromBase64String($pngData)
    [System.IO.File]::WriteAllBytes($testImagePath, $bytes)
    Write-Host "✅ Test image created: $testImagePath" -ForegroundColor Green
} catch {
    Write-Host "❌ Failed to create test image: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Test API endpoint
$apiUrl = "http://localhost:8080/api/admin/customer-installation/report-installations"

Write-Host "`n=== Testing API Endpoint ===" -ForegroundColor Blue
Write-Host "API URL: $apiUrl" -ForegroundColor Yellow

# Create multipart form data using curl
$curlCommand = @"
curl -X POST "$apiUrl" \
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

Write-Host "Curl command to test:" -ForegroundColor Yellow
Write-Host $curlCommand -ForegroundColor White

Write-Host "`n=== Manual Test Steps ===" -ForegroundColor Green
Write-Host "1. Make sure backend is running: go run cmd/myapp/main.go" -ForegroundColor Yellow
Write-Host "2. Get your authentication token from browser" -ForegroundColor Yellow
Write-Host "3. Replace YOUR_TOKEN_HERE with your actual token" -ForegroundColor Yellow
Write-Host "4. Run the curl command above" -ForegroundColor Yellow
Write-Host "5. Check backend logs for document photo upload messages" -ForegroundColor Yellow
Write-Host "6. Check database for document_photo value" -ForegroundColor Yellow

# Clean up test file
if (Test-Path $testImagePath) {
    Remove-Item $testImagePath
    Write-Host "`n✅ Test image cleaned up" -ForegroundColor Green
}

Write-Host "`n=== Expected Backend Logs ===" -ForegroundColor Blue
Write-Host "You should see these log messages in backend:" -ForegroundColor Yellow
Write-Host "- 'Document photo upload started - filename: test-ktp.png, size: X'" -ForegroundColor White
Write-Host "- 'Document photo uploaded successfully - filePath: uploads/installations/documents/document_YYYYMMDD_HHMMSS.png'" -ForegroundColor White
Write-Host "- 'Installation report request data - customer_id: ..., document_photo: ...'" -ForegroundColor White

Read-Host "Press Enter to continue"
