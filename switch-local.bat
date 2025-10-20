@echo off
echo Switching backend to LOCAL environment...
node switch-env.js local
echo.
echo ✅ Environment switched to LOCAL!
echo.
echo Configuration:
echo   - Database: localhost (iqgncnzy_skripsi)
echo   - MikroTik: 10.10.9.203:22 (Local)
echo   - Port: 3001
echo.
echo You can now start the server with start-server.bat
pause

