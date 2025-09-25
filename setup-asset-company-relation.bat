@echo off
echo Setting up Asset-Company Relationship...

echo Applying database migration...
mysql -u root -p iqgncnzy_skripsi < database-migration-asset-company-relation.sql

echo Migration completed!
echo.
echo Changes made:
echo - Added company_id column to assets table
echo - Removed site column from assets table  
echo - Added foreign key relationship to company table
echo - Added index for better performance
echo.
echo Please restart your backend server to use the updated code.
pause

