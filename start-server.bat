@echo off
echo 🔧 Checking for existing Go processes on port 3001...

REM Check if port 3001 is in use
netstat -ano | findstr :3001 > nul
if %errorlevel% equ 0 (
    echo ⚠️  Port 3001 is already in use. Finding and killing process...
    for /f "tokens=5" %%a in ('netstat -ano ^| findstr :3001') do (
        echo Killing process %%a
        taskkill /F /PID %%a > nul 2>&1
    )
    echo ✅ Cleared port 3001
) else (
    echo ✅ Port 3001 is available
)

echo 🚀 Starting Go server...
go run cmd/myapp/main.go
