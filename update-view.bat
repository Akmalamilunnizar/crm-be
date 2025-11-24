@echo off
REM Update Installation Report View
REM This script will help you update the database view

echo ========================================
echo Installation Report View Update Script
echo ========================================
echo.
echo This script will update the installation_report_complete view in your database.
echo.
echo Please follow these steps:
echo.
echo 1. Open your MySQL client (phpMyAdmin, MySQL Workbench, or command line)
echo 2. Connect to your database
echo 3. Execute the SQL file: update-installation-report-view.sql
echo.
echo Alternatively, you can run:
echo    mysql -h [host] -u [username] -p [database] ^< update-installation-report-view.sql
echo.
echo The SQL file is located in the same directory as this script.
echo.
echo ========================================
echo.
pause

REM Try to run the Go program to create the view
echo.
echo Attempting to run Go program to create view...
echo.
go run create_view.go

if %ERRORLEVEL% EQU 0 (
    echo.
    echo SUCCESS! View updated successfully.
    echo.
) else (
    echo.
    echo ERROR: Failed to update view using Go program.
    echo Please update the view manually using the SQL file.
    echo.
)

pause
