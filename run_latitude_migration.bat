@echo off
echo Applying latitude/longitude migration...

"C:\laragon\bin\mysql\mysql-8.4.3-winx64\bin\mysql.exe" -u root iqgncnzy_skripsi < migrations/add_latitude_longitude_to_customer_installations.sql

if %errorlevel% equ 0 (
    echo Migration completed successfully!
) else (
    echo Migration failed with error code: %errorlevel%
)

pause
