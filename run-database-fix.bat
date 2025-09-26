@echo off
echo Running database fix for installation reports...

REM Run the database fix script
mysql -u root -p iqgncnzy_skripsi < fix-installation-report-view.sql

echo.
echo Database fix completed. Now running debug script...

REM Run the debug script
mysql -u root -p iqgncnzy_skripsi < debug-installation-reports.sql

echo.
echo Debug completed. Check the output above for any issues.
pause

