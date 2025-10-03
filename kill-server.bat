@echo off
echo 🔧 Killing all Go processes on port 3001...

REM Kill all processes using port 3001
for /f "tokens=5" %%a in ('netstat -ano ^| findstr :3001') do (
    echo Killing process %%a
    taskkill /F /PID %%a > nul 2>&1
)

echo ✅ All processes on port 3001 have been terminated
pause
