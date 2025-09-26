# Debug document photo upload

Write-Host "=== Debug Document Photo Upload ===" -ForegroundColor Green

# Create a simple test image file
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

# Test with curl (you'll need to replace YOUR_TOKEN with actual token)
$curlCommand = @"
curl -X POST "$apiUrl" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "customer_id=0ae7d407-2656-4fe3-878f-89c19abcbdac" \
  -F "technician_id=c13a6c87-ec28-47ba-84c2-58b5ace2af57" \
  -F "assets_id=5ca1606b-66b3-4958-b7af-f48d4cda800a" \
  -F "document_type=KTP" \
  -F "document_photo=@$testImagePath" \
  -F "status=completed" \
  -F "installation_type=new_installation"
"@

Write-Host "Curl command:" -ForegroundColor Yellow
Write-Host $curlCommand -ForegroundColor White

Write-Host "`n=== Check Backend Logs ===" -ForegroundColor Blue
Write-Host "After running the curl command, check backend console for these logs:" -ForegroundColor Yellow
Write-Host "1. 'Document photo upload started - filename: test-ktp.png, size: X'" -ForegroundColor White
Write-Host "2. 'Document photo uploaded successfully - filePath: uploads/installations/documents/document_YYYYMMDD_HHMMSS.png'" -ForegroundColor White
Write-Host "3. 'Installation report request data - customer_id: ..., document_photo: ...'" -ForegroundColor White

Write-Host "`n=== Check Database ===" -ForegroundColor Blue
Write-Host "After successful upload, check database with:" -ForegroundColor Yellow
Write-Host "SELECT id, document_type, document_photo FROM customer_installations ORDER BY createdAt DESC LIMIT 1;" -ForegroundColor White

Write-Host "`n=== Check Uploads Directory ===" -ForegroundColor Blue
Write-Host "Check if file exists in: uploads/installations/documents/" -ForegroundColor Yellow

# Clean up test file
if (Test-Path $testImagePath) {
    Remove-Item $testImagePath
    Write-Host "`n✅ Test image cleaned up" -ForegroundColor Green
}

Read-Host "Press Enter to continue"
