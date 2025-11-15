@echo off
REM Start WhatsApp service script for Windows

echo Starting WhatsApp Service...
cd /d "%~dp0"

REM Check if node_modules exists
if not exist "node_modules" (
    echo Installing dependencies...
    call npm install
)

REM Start the service
echo Starting service on port %WHATSAPP_SERVICE_PORT% (default: 3002)...
call npm start






