@echo off
echo Starting database migration...

"C:\laragon\bin\mysql\mysql-8.4.3-winx64\bin\mysql.exe" -u root iqgncnzy_skripsi < migration_script.sql

if %errorlevel% equ 0 (
    echo Migration completed successfully!
    echo Database structure has been updated to match lilly.sql
) else (
    echo Migration failed with error code: %errorlevel%
)

echo Migration process completed.
pause

