# Test API endpoint for installation reports

Write-Host "=== Testing API Endpoint ===" -ForegroundColor Green

# Test the API endpoint
$apiUrl = "http://localhost:8080/api/admin/customer-installation/report-complete"

Write-Host "Testing API endpoint: $apiUrl" -ForegroundColor Yellow

try {
    # Test with curl if available
    $curlResult = & curl -X GET $apiUrl -H "Content-Type: application/json" 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "API Response:" -ForegroundColor Green
        Write-Host $curlResult -ForegroundColor White
    } else {
        Write-Host "Curl not available or API not responding" -ForegroundColor Red
        Write-Host "Make sure backend is running: go run cmd/myapp/main.go" -ForegroundColor Yellow
    }
    
} catch {
    Write-Host "Error testing API: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n=== Manual Test Steps ===" -ForegroundColor Green
Write-Host "1. Start backend: go run cmd/myapp/main.go" -ForegroundColor Yellow
Write-Host "2. Open browser and go to: http://localhost:8080/api/admin/customer-installation/report-complete" -ForegroundColor Yellow
Write-Host "3. Check if you get JSON response with installation reports" -ForegroundColor Yellow
Write-Host "4. If you get authentication error, add Authorization header with Bearer token" -ForegroundColor Yellow

Read-Host "Press Enter to continue"
