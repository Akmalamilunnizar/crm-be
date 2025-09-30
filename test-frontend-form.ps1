# Test frontend form submission

Write-Host "=== Test Frontend Form Submission ===" -ForegroundColor Green

Write-Host "`n=== Frontend Form Test Steps ===" -ForegroundColor Blue
Write-Host "1. Open your browser and go to the frontend" -ForegroundColor Yellow
Write-Host "2. Navigate to Customer Actions > Add Report Installation" -ForegroundColor Yellow
Write-Host "3. Fill in the form with these test values:" -ForegroundColor Yellow

Write-Host "`n=== Test Form Data ===" -ForegroundColor Blue
Write-Host "Customer ID: 0ae7d407-2656-4fe3-878f-89c19abcbdac" -ForegroundColor White
Write-Host "Technician ID: c13a6c87-ec28-47ba-84c2-58b5ace2af57" -ForegroundColor White
Write-Host "Assets ID: 5ca1606b-66b3-4958-b7af-f48d4cda800a" -ForegroundColor White
Write-Host "Document Type: KTP" -ForegroundColor White
Write-Host "Document Photo: Upload any KTP image file" -ForegroundColor White
Write-Host "Status: completed" -ForegroundColor White
Write-Host "Installation Type: new_installation" -ForegroundColor White

Write-Host "`n=== What to Watch ===" -ForegroundColor Blue
Write-Host "4. Open browser Developer Tools (F12)" -ForegroundColor Yellow
Write-Host "5. Go to Network tab" -ForegroundColor Yellow
Write-Host "6. Submit the form" -ForegroundColor Yellow
Write-Host "7. Look for the POST request to '/api/admin/customer-installation/report-installations'" -ForegroundColor Yellow
Write-Host "8. Check the request payload - it should include the document_photo file" -ForegroundColor Yellow

Write-Host "`n=== Backend Console to Watch ===" -ForegroundColor Blue
Write-Host "Watch the backend console for these log messages:" -ForegroundColor Yellow
Write-Host "1. 'Document photo upload started - filename: [filename], size: [size]'" -ForegroundColor White
Write-Host "2. 'Document photo uploaded successfully - filePath: [path]'" -ForegroundColor White
Write-Host "3. 'Installation report request data - customer_id: ..., document_photo: [path]'" -ForegroundColor White
Write-Host "4. 'Creating installation record - DocumentPhoto: [path], DocumentType: KTP'" -ForegroundColor White
Write-Host "5. 'Installation record created successfully with ID: [id]'" -ForegroundColor White

Write-Host "`n=== If Form Submission Fails ===" -ForegroundColor Red
Write-Host "Check browser console for errors:" -ForegroundColor Yellow
Write-Host "- Network errors (404, 500, etc.)" -ForegroundColor White
Write-Host "- JavaScript errors" -ForegroundColor White
Write-Host "- Form validation errors" -ForegroundColor White

Write-Host "`n=== If Backend Logs Don't Show ===" -ForegroundColor Red
Write-Host "The issue might be:" -ForegroundColor Yellow
Write-Host "1. Backend not running" -ForegroundColor White
Write-Host "2. Wrong API endpoint being called" -ForegroundColor White
Write-Host "3. Authentication issues" -ForegroundColor White
Write-Host "4. Form data not being sent correctly" -ForegroundColor White

Write-Host "`n=== After Successful Form Submission ===" -ForegroundColor Green
Write-Host "1. Check database for new record" -ForegroundColor Yellow
Write-Host "2. Verify document_photo field has file path" -ForegroundColor Yellow
Write-Host "3. Check uploads/installations/documents/ directory" -ForegroundColor Yellow
Write-Host "4. Test document photo display in detail page" -ForegroundColor Yellow

Write-Host "`n=== Database Check After Form Submission ===" -ForegroundColor Blue
Write-Host "Run this SQL query to check the latest record:" -ForegroundColor Yellow
Write-Host "SELECT id, customer_id, document_type, document_photo, status, createdAt FROM customer_installations ORDER BY createdAt DESC LIMIT 1;" -ForegroundColor Cyan

Read-Host "Press Enter to continue"
