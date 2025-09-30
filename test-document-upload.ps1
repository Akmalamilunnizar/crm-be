# Test document photo upload functionality

Write-Host "=== Testing Document Photo Upload ===" -ForegroundColor Green

# Create a test image file
$testImagePath = "test-document.jpg"
$testImageContent = @"
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
"@

# Convert base64 to bytes and create test image
try {
    $bytes = [System.Convert]::FromBase64String($testImageContent)
    [System.IO.File]::WriteAllBytes($testImagePath, $bytes)
    Write-Host "✅ Test image created: $testImagePath" -ForegroundColor Green
} catch {
    Write-Host "❌ Failed to create test image: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Test API endpoint with multipart form data
$apiUrl = "http://localhost:8080/api/admin/customer-installation/report-installations"

Write-Host "`n=== Testing API Endpoint ===" -ForegroundColor Blue
Write-Host "API URL: $apiUrl" -ForegroundColor Yellow

# Create multipart form data
$boundary = [System.Guid]::NewGuid().ToString()
$LF = "`r`n"

$bodyLines = @(
    "--$boundary",
    "Content-Disposition: form-data; name=`"customer_id`"",
    "",
    "test-customer-id",
    "--$boundary",
    "Content-Disposition: form-data; name=`"technician_id`"",
    "",
    "test-technician-id",
    "--$boundary",
    "Content-Disposition: form-data; name=`"assets_id`"",
    "",
    "test-assets-id",
    "--$boundary",
    "Content-Disposition: form-data; name=`"document_type`"",
    "",
    "KTP",
    "--$boundary",
    "Content-Disposition: form-data; name=`"document_photo`"; filename=`"test-document.jpg`"",
    "Content-Type: image/jpeg",
    "",
    [System.Text.Encoding]::UTF8.GetString($bytes),
    "--$boundary--"
)

$body = $bodyLines -join $LF
$bodyBytes = [System.Text.Encoding]::UTF8.GetBytes($body)

try {
    $response = Invoke-WebRequest -Uri $apiUrl -Method POST -Body $bodyBytes -ContentType "multipart/form-data; boundary=$boundary" -TimeoutSec 30
    
    Write-Host "✅ API request successful" -ForegroundColor Green
    Write-Host "Status Code: $($response.StatusCode)" -ForegroundColor White
    Write-Host "Response:" -ForegroundColor Yellow
    Write-Host $response.Content -ForegroundColor White
    
} catch {
    Write-Host "❌ API request failed: $($_.Exception.Message)" -ForegroundColor Red
    
    if ($_.Exception.Response) {
        $statusCode = $_.Exception.Response.StatusCode
        Write-Host "Status Code: $statusCode" -ForegroundColor Red
        
        try {
            $errorStream = $_.Exception.Response.GetResponseStream()
            $reader = New-Object System.IO.StreamReader($errorStream)
            $errorContent = $reader.ReadToEnd()
            Write-Host "Error Response: $errorContent" -ForegroundColor Red
        } catch {
            Write-Host "Could not read error response" -ForegroundColor Red
        }
    }
}

# Clean up test file
if (Test-Path $testImagePath) {
    Remove-Item $testImagePath
    Write-Host "`n✅ Test image cleaned up" -ForegroundColor Green
}

Write-Host "`n=== Next Steps ===" -ForegroundColor Green
Write-Host "1. Make sure backend is running: go run cmd/myapp/main.go" -ForegroundColor Yellow
Write-Host "2. Check if document photo is saved to database" -ForegroundColor Yellow
Write-Host "3. Verify file is uploaded to uploads/installations/documents/" -ForegroundColor Yellow

Read-Host "Press Enter to continue"
