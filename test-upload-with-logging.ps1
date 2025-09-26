# Test upload with detailed logging

Write-Host "=== Test Upload with Detailed Logging ===" -ForegroundColor Green

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

Write-Host "`n=== Backend Logs to Watch ===" -ForegroundColor Blue
Write-Host "Watch the backend console for these log messages:" -ForegroundColor Yellow
Write-Host "1. 'Document photo upload started - filename: test-ktp.png, size: X'" -ForegroundColor White
Write-Host "2. 'Document photo uploaded successfully - filePath: uploads/installations/documents/document_YYYYMMDD_HHMMSS.png'" -ForegroundColor White
Write-Host "3. 'Installation report request data - customer_id: ..., document_photo: ...'" -ForegroundColor White
Write-Host "4. 'Creating installation record - DocumentPhoto: uploads/installations/documents/document_YYYYMMDD_HHMMSS.png, DocumentType: KTP'" -ForegroundColor White
Write-Host "5. 'Installation record created successfully with ID: ...'" -ForegroundColor White

Write-Host "`n=== Test with curl ===" -ForegroundColor Blue
Write-Host "Run this curl command (replace YOUR_TOKEN with actual token):" -ForegroundColor Yellow

$curlCommand = @"
curl -X POST "http://localhost:8080/api/admin/customer-installation/report-installations" \
  -H "Authorization: Bearer YOUR_TOKEN" \
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

Write-Host "`n=== After Upload Test ===" -ForegroundColor Blue
Write-Host "1. Check backend logs for all 5 log messages above" -ForegroundColor Yellow
Write-Host "2. Check database with: SELECT id, document_type, document_photo FROM customer_installations ORDER BY createdAt DESC LIMIT 1;" -ForegroundColor Yellow
Write-Host "3. Check uploads directory: uploads/installations/documents/" -ForegroundColor Yellow

# Clean up test file
if (Test-Path $testImagePath) {
    Remove-Item $testImagePath
    Write-Host "`n✅ Test image cleaned up" -ForegroundColor Green
}

Read-Host "Press Enter to continue"
