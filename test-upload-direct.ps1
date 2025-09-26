# Test upload directly with detailed debugging

Write-Host "=== Direct Upload Test with Debugging ===" -ForegroundColor Green

# Create a simple test image file
$testImagePath = "test-ktp.png"
$pngData = @"
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
"@

try {
    $bytes = [System.Convert]::FromBase64String($pngData)
    [System.IO.File]::WriteAllBytes($testImagePath, $bytes)
    Write-Host "✅ Test image created: $testImagePath" -ForegroundColor Green
} catch {
    Write-Host "❌ Failed to create test image: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

Write-Host "`n=== IMPORTANT: Get Your Token First ===" -ForegroundColor Red
Write-Host "1. Open your browser and login to the frontend" -ForegroundColor Yellow
Write-Host "2. Open Developer Tools (F12)" -ForegroundColor Yellow
Write-Host "3. Go to Application/Storage tab" -ForegroundColor Yellow
Write-Host "4. Find the 'token' cookie and copy its value" -ForegroundColor Yellow
Write-Host "5. Replace 'YOUR_TOKEN_HERE' in the curl command below" -ForegroundColor Yellow

Write-Host "`n=== Curl Command for Testing ===" -ForegroundColor Blue
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

Write-Host "`n=== What to Watch in Backend Console ===" -ForegroundColor Blue
Write-Host "After running the curl command, you should see these logs in order:" -ForegroundColor Yellow
Write-Host "1. 'Document photo upload started - filename: test-ktp.png, size: 95'" -ForegroundColor White
Write-Host "2. 'Document photo uploaded successfully - filePath: uploads/installations/documents/document_YYYYMMDD_HHMMSS.png'" -ForegroundColor White
Write-Host "3. 'Installation report request data - customer_id: ..., document_photo: uploads/installations/documents/document_YYYYMMDD_HHMMSS.png'" -ForegroundColor White
Write-Host "4. 'Creating installation record - DocumentPhoto: uploads/installations/documents/document_YYYYMMDD_HHMMSS.png, DocumentType: KTP'" -ForegroundColor White
Write-Host "5. 'Installation record created successfully with ID: ...'" -ForegroundColor White

Write-Host "`n=== If You See Errors ===" -ForegroundColor Red
Write-Host "If you see any errors, check:" -ForegroundColor Yellow
Write-Host "- Is backend running? (go run cmd/myapp/main.go)" -ForegroundColor White
Write-Host "- Is the token correct?" -ForegroundColor White
Write-Host "- Are the customer_id, technician_id, assets_id valid?" -ForegroundColor White
Write-Host "- Check backend console for error messages" -ForegroundColor White

Write-Host "`n=== After Successful Upload ===" -ForegroundColor Green
Write-Host "1. Check database: SELECT id, document_type, document_photo FROM customer_installations ORDER BY createdAt DESC LIMIT 1;" -ForegroundColor Yellow
Write-Host "2. Check uploads directory: uploads/installations/documents/" -ForegroundColor Yellow
Write-Host "3. Test frontend form submission" -ForegroundColor Yellow

# Clean up test file
if (Test-Path $testImagePath) {
    Remove-Item $testImagePath
    Write-Host "`n✅ Test image cleaned up" -ForegroundColor Green
}

Read-Host "Press Enter to continue"
