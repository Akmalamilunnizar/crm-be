# Check backend status and test endpoints

Write-Host "=== Backend Status Check ===" -ForegroundColor Green

# Check if backend is running
$backendUrl = "http://localhost:8080"
$testEndpoint = "$backendUrl/api/admin/customer-installation/report-complete"

Write-Host "Checking backend at: $backendUrl" -ForegroundColor Yellow

try {
    # Test basic connectivity
    $response = Invoke-WebRequest -Uri $backendUrl -Method GET -TimeoutSec 5
    Write-Host "✅ Backend is running" -ForegroundColor Green
    Write-Host "Status Code: $($response.StatusCode)" -ForegroundColor White
} catch {
    Write-Host "❌ Backend is not running or not accessible" -ForegroundColor Red
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "`nPlease start backend with: go run cmd/myapp/main.go" -ForegroundColor Yellow
    exit 1
}

# Test API endpoint (without auth for now)
Write-Host "`n=== Testing API Endpoints ===" -ForegroundColor Blue

try {
    $apiResponse = Invoke-WebRequest -Uri $testEndpoint -Method GET -TimeoutSec 5
    Write-Host "✅ API endpoint is accessible" -ForegroundColor Green
    Write-Host "Status Code: $($apiResponse.StatusCode)" -ForegroundColor White
} catch {
    if ($_.Exception.Response.StatusCode -eq 401) {
        Write-Host "⚠️  API endpoint requires authentication (401 Unauthorized)" -ForegroundColor Yellow
        Write-Host "This is expected - the endpoint is working but needs auth token" -ForegroundColor White
    } else {
        Write-Host "❌ API endpoint test failed" -ForegroundColor Red
        Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    }
}

# Check uploads directory
Write-Host "`n=== Checking Uploads Directory ===" -ForegroundColor Blue

$uploadsDir = "uploads/installations/documents"
if (Test-Path $uploadsDir) {
    Write-Host "✅ Uploads directory exists: $uploadsDir" -ForegroundColor Green
    
    $files = Get-ChildItem -Path $uploadsDir -File -ErrorAction SilentlyContinue
    if ($files) {
        Write-Host "📁 Found $($files.Count) files in uploads directory:" -ForegroundColor Yellow
        foreach ($file in $files) {
            Write-Host "  - $($file.Name) ($($file.Length) bytes)" -ForegroundColor White
        }
    } else {
        Write-Host "📁 Uploads directory is empty" -ForegroundColor Yellow
    }
} else {
    Write-Host "⚠️  Uploads directory does not exist: $uploadsDir" -ForegroundColor Yellow
    Write-Host "It will be created automatically when first file is uploaded" -ForegroundColor White
}

# Check if backend process is running
Write-Host "`n=== Checking Backend Process ===" -ForegroundColor Blue

$goProcesses = Get-Process -Name "go" -ErrorAction SilentlyContinue
if ($goProcesses) {
    Write-Host "✅ Go processes are running:" -ForegroundColor Green
    foreach ($proc in $goProcesses) {
        Write-Host "  - PID: $($proc.Id), CPU: $($proc.CPU), Memory: $([math]::Round($proc.WorkingSet64/1MB, 2)) MB" -ForegroundColor White
    }
} else {
    Write-Host "⚠️  No Go processes found" -ForegroundColor Yellow
    Write-Host "Backend might be running in a different way" -ForegroundColor White
}

Write-Host "`n=== Next Steps ===" -ForegroundColor Green
Write-Host "1. If backend is running, test document photo upload" -ForegroundColor Yellow
Write-Host "2. Use the curl command from test-simple-upload.ps1" -ForegroundColor Yellow
Write-Host "3. Check backend logs for upload messages" -ForegroundColor Yellow
Write-Host "4. Verify document_photo is saved to database" -ForegroundColor Yellow

Read-Host "Press Enter to continue"
