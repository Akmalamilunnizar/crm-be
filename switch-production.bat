@echo off
echo Switching backend to PRODUCTION (VPS) environment...
node switch-env.js production
echo.
echo ⚠️  WARNING: You are now in PRODUCTION mode!
echo.
echo Configuration:
echo   - Database: localhost on VPS (menarane_lilly)
echo   - MikroTik: 103.148.18.52:8252 (Production)
echo   - Port: 3001
echo.
echo ⚠️  IMPORTANT: Update YOUR_DB_PASSWORD_HERE in .env before deploying!
echo.
echo Deploy to VPS using:
echo   - Use deploy-to-vps.ps1 or deploy-to-vps.sh
echo   - Or manually: scp to VPS and rebuild
pause

