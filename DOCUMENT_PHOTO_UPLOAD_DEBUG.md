# Document Photo Upload Debug Guide

## Problem
From the database screenshot, we can see that:
- `document_type` is "KTP" (correctly saved)
- `document_photo` field is empty (the problem)
- `notes` field has "frewferferf" (correctly saved)

This indicates that text fields are working but file upload is not.

## Enhanced Debug Logging

### Backend Controller Enhanced
Added comprehensive logging to `report_controller_new.go`:

```go
// Log all form files for debugging
form, formErr := ctx.MultipartForm()
if formErr == nil {
    log.Printf("Multipart form files found: %v", len(form.File))
    for fieldName, files := range form.File {
        log.Printf("Form field '%s' has %d files", fieldName, len(files))
        for i, file := range files {
            log.Printf("  File %d: %s (size: %d)", i, file.Filename, file.Size)
        }
    }
}

if file, err := ctx.FormFile("document_photo"); err == nil {
    // File upload success logging
} else {
    log.Printf("No document photo uploaded - error: %s", err.Error())
    
    // Check if document_photo field exists in form values
    if docPhotoValue := ctx.FormValue("document_photo"); docPhotoValue != "" {
        log.Printf("document_photo form value (not file): %s", docPhotoValue)
    } else {
        log.Printf("No document_photo field found in form values")
    }
}
```

## Debug Steps

### Step 1: Check Backend Logs
Start backend and watch for these logs when submitting form:

```bash
cd crm-be
go run cmd/myapp/main.go
```

### Step 2: Expected Logs (Success Case)
```
Multipart form files found: 1
Form field 'document_photo' has 1 files
  File 0: example.png (size: 12345)
Document photo upload started - filename: example.png, size: 12345
Document photo uploaded successfully - filePath: uploads/installations/documents/document_20250925_123456.png
Installation report request data - customer_id: ..., document_photo: uploads/installations/documents/document_20250925_123456.png, ...
Creating installation record - DocumentPhoto: uploads/installations/documents/document_20250925_123456.png, DocumentType: KTP
Installation record created successfully with ID: ...
```

### Step 3: Expected Logs (Failure Case)
```
Multipart form files found: 0
No document photo uploaded - error: request Content-Type isn't multipart/form-data
No document_photo field found in form values
Installation report request data - customer_id: ..., document_photo: , ...
Creating installation record - DocumentPhoto: nil, DocumentType: KTP
Installation record created successfully with ID: ...
```

## Troubleshooting Scenarios

### Scenario 1: No Files in Multipart Form
**Log**: `Multipart form files found: 0`
**Cause**: Frontend not sending file or Content-Type issue
**Solution**: Check frontend FormData implementation

### Scenario 2: Files Found but Wrong Field Name
**Log**: `Form field 'other_field' has 1 files` (not 'document_photo')
**Cause**: Frontend using wrong field name
**Solution**: Check FormData.append('document_photo', file)

### Scenario 3: File Found but Upload Fails
**Log**: `Document photo upload started` but no `uploaded successfully`
**Cause**: File validation or save error
**Solution**: Check file type validation and filesystem permissions

### Scenario 4: File Uploaded but Path Not Saved
**Log**: `uploaded successfully` but database still empty
**Cause**: Database insertion issue
**Solution**: Check entity field mapping

## Manual Test with cURL

Create test file and run:

```bash
# Create test image
echo "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==" | base64 -d > test-ktp.png

# Test upload
curl -X POST "http://localhost:8080/api/admin/customer-installation/report-installations" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -F "customer_id=0ae7d407-2656-4fe3-878f-89c19abcbdac" \
  -F "technician_id=c13a6c87-ec28-47ba-84c2-58b5ace2af57" \
  -F "assets_id=5ca1606b-66b3-4958-b7af-f48d4cda800a" \
  -F "document_type=KTP" \
  -F "document_photo=@test-ktp.png" \
  -F "status=completed" \
  -F "installation_type=new_installation" \
  -F "notes=Test document upload debugging"
```

## Frontend Debug

### Check FormData in Browser
Add this to frontend submit function:

```javascript
// Debug FormData before sending
for (let [key, value] of formData.entries()) {
  console.log(key, value);
  if (value instanceof File) {
    console.log(`File: ${value.name}, size: ${value.size}, type: ${value.type}`);
  }
}
```

### Check Network Tab
1. Open Browser DevTools
2. Go to Network tab
3. Submit form
4. Look for request to `/api/admin/customer-installation/report-installations`
5. Check Request Headers: `Content-Type: multipart/form-data; boundary=...`
6. Check Request Payload: Should show file data

## Database Verification

After successful upload, check:

```sql
SELECT id, document_type, document_photo, notes, createdAt 
FROM customer_installations 
ORDER BY createdAt DESC 
LIMIT 1;
```

Expected result:
- `document_type`: "KTP"
- `document_photo`: "uploads/installations/documents/document_20250925_123456.png"
- `notes`: Your test notes

## File System Check

Check if file actually exists:

```bash
ls -la uploads/installations/documents/
```

Should show uploaded files with timestamps.

## Next Steps

1. **Start Backend**: `go run cmd/myapp/main.go`
2. **Submit Form**: Use frontend form or cURL
3. **Watch Logs**: Monitor backend console for debug logs
4. **Identify Issue**: Based on log patterns above
5. **Fix Issue**: Apply appropriate solution
6. **Verify Database**: Check document_photo field

The enhanced logging will now show exactly where the file upload process is failing, making it much easier to identify and fix the issue.
