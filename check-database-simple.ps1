# Simple database check without password prompt

Write-Host "=== Simple Database Check ===" -ForegroundColor Green

Write-Host "`n=== Manual Database Check ===" -ForegroundColor Blue
Write-Host "Please run these SQL commands in your database tool (phpMyAdmin, MySQL Workbench, etc.):" -ForegroundColor Yellow

Write-Host "`n1. Check table structure:" -ForegroundColor White
Write-Host "DESCRIBE customer_installations;" -ForegroundColor Cyan

Write-Host "`n2. Check recent installations:" -ForegroundColor White
Write-Host "SELECT id, customer_id, technician_id, document_type, document_photo, status, createdAt FROM customer_installations ORDER BY createdAt DESC LIMIT 5;" -ForegroundColor Cyan

Write-Host "`n3. Check document photo statistics:" -ForegroundColor White
Write-Host "SELECT COUNT(*) as total_records, COUNT(document_photo) as records_with_document_photo, COUNT(CASE WHEN document_photo IS NOT NULL AND document_photo != '' THEN 1 END) as non_empty_document_photos FROM customer_installations;" -ForegroundColor Cyan

Write-Host "`n4. Check installations with document photos:" -ForegroundColor White
Write-Host "SELECT id, document_type, document_photo, LENGTH(document_photo) as photo_path_length, createdAt FROM customer_installations WHERE document_photo IS NOT NULL AND document_photo != '' ORDER BY createdAt DESC;" -ForegroundColor Cyan

Write-Host "`n=== Expected Results ===" -ForegroundColor Blue
Write-Host "If document upload is working:" -ForegroundColor Yellow
Write-Host "- document_photo column should contain file paths like 'uploads/installations/documents/document_YYYYMMDD_HHMMSS.png'" -ForegroundColor White
Write-Host "- non_empty_document_photos count should be > 0" -ForegroundColor White
Write-Host "- Recent installations should show document_photo values" -ForegroundColor White

Write-Host "`n=== If Document Photos Are Still Not Saving ===" -ForegroundColor Red
Write-Host "The issue might be:" -ForegroundColor Yellow
Write-Host "1. Backend not receiving the file upload" -ForegroundColor White
Write-Host "2. File upload failing" -ForegroundColor White
Write-Host "3. Database insertion failing" -ForegroundColor White
Write-Host "4. Form data not being parsed correctly" -ForegroundColor White

Write-Host "`n=== Next Steps ===" -ForegroundColor Green
Write-Host "1. Run the SQL queries above to check current state" -ForegroundColor Yellow
Write-Host "2. Test upload with: .\test-upload-direct.ps1" -ForegroundColor Yellow
Write-Host "3. Watch backend console for log messages" -ForegroundColor Yellow
Write-Host "4. Check if files are being uploaded to uploads/installations/documents/" -ForegroundColor Yellow

Read-Host "Press Enter to continue"
