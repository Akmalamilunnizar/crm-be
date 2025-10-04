# PowerShell script to start Go server with port conflict handling
Write-Host "🔧 Checking for existing Go processes on port 3001..." -ForegroundColor Yellow

# Check if port 3001 is in use
$portInUse = Get-NetTCPConnection -LocalPort 3001 -ErrorAction SilentlyContinue

if ($portInUse) {
    Write-Host "⚠️  Port 3001 is already in use. Finding and killing process..." -ForegroundColor Red
    
    foreach ($connection in $portInUse) {
        $pid = $connection.OwningProcess
        Write-Host "Killing process $pid" -ForegroundColor Yellow
        Stop-Process -Id $pid -Force -ErrorAction SilentlyContinue
    }
    
    Write-Host "✅ Cleared port 3001" -ForegroundColor Green
    Start-Sleep -Seconds 2
} else {
    Write-Host "✅ Port 3001 is available" -ForegroundColor Green
}

Write-Host "🚀 Starting Go server..." -ForegroundColor Cyan
go run cmd/myapp/main.go
