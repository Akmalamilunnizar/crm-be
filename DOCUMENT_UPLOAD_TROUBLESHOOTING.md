# Document Upload Troubleshooting - Complete Guide

## Problem Description
Document photos (KTP) uploaded through the form are not being saved to the database in the `customer_installations.document_photo` column, despite the database schema being correct (VARCHAR(255), nullable).

## Root Cause Analysis
The issue could be in several places:

1. **Frontend Form Submission**: Form data not being sent correctly
2. **Backend File Upload**: File upload handling failing
3. **Backend Form Parsing**: Multipart form data not being parsed correctly
4. **Database Insertion**: Database insertion failing silently
5. **Authentication**: API endpoint requiring authentication

## Comprehensive Testing Approach

### 1. ✅ Database Schema Verification

#### **Check Database Schema**:
```sql
-- Check table structure
DESCRIBE customer_installations;

-- Check document_photo column
SELECT 
    COLUMN_NAME,
    DATA_TYPE,
    IS_NULLABLE,
    COLUMN_DEFAULT
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_SCHEMA = 'iqgncnzy_skripsi' 
AND TABLE_NAME = 'customer_installations'
AND COLUMN_NAME = 'document_photo';

-- Check recent installations
SELECT id, customer_id, document_type, document_photo, status, createdAt 
FROM customer_installations 
ORDER BY createdAt DESC 
LIMIT 5;

-- Check document photo statistics
SELECT 
    COUNT(*) as total_records,
    COUNT(document_photo) as records_with_document_photo,
    COUNT(CASE WHEN document_photo IS NOT NULL AND document_photo != '' THEN 1 END) as non_empty_document_photos
FROM customer_installations;
```

### 2. ✅ Backend Testing

#### **Test Backend Status**:
```powershell
cd crm-be
.\debug-upload-issue.ps1
```

#### **Test Direct Upload with curl**:
```powershell
cd crm-be
.\test-upload-direct.ps1
```

#### **Expected Backend Logs**:
When document upload works, you should see these logs in sequence:
```
1. Document photo upload started - filename: test-ktp.png, size: 95
2. Document photo uploaded successfully - filePath: uploads/installations/documents/document_20250101_120000.png
3. Installation report request data - customer_id: ..., document_photo: uploads/installations/documents/document_20250101_120000.png
4. Creating installation record - DocumentPhoto: uploads/installations/documents/document_20250101_120000.png, DocumentType: KTP
5. Installation record created successfully with ID: 12345678-1234-1234-1234-123456789012
```

### 3. ✅ Frontend Testing

#### **Test Frontend Form**:
```powershell
cd crm-be
.\test-frontend-form.ps1
```

#### **Frontend Form Test Steps**:
1. Open browser and go to frontend
2. Navigate to Customer Actions > Add Report Installation
3. Fill in form with test data
4. Upload a KTP image file
5. Submit form
6. Check browser Developer Tools > Network tab
7. Look for POST request to `/api/admin/customer-installation/report-installations`
8. Check request payload includes document_photo file

### 4. ✅ Database Verification

#### **Check Database After Upload**:
```powershell
cd crm-be
.\check-database-simple.ps1
```

#### **SQL Queries to Run**:
```sql
-- Check latest record
SELECT id, customer_id, document_type, document_photo, status, createdAt 
FROM customer_installations 
ORDER BY createdAt DESC 
LIMIT 1;

-- Check installations with document photos
SELECT id, document_type, document_photo, LENGTH(document_photo) as photo_path_length, createdAt 
FROM customer_installations 
WHERE document_photo IS NOT NULL AND document_photo != '' 
ORDER BY createdAt DESC;
```

## Troubleshooting Steps

### Step 1: Verify Backend is Running
```bash
cd crm-be
go run cmd/myapp/main.go
```

### Step 2: Test Backend Status
```powershell
cd crm-be
.\debug-upload-issue.ps1
```

### Step 3: Test Direct Upload
```powershell
cd crm-be
.\test-upload-direct.ps1
```

### Step 4: Test Frontend Form
```powershell
cd crm-be
.\test-frontend-form.ps1
```

### Step 5: Check Database
```powershell
cd crm-be
.\check-database-simple.ps1
```

## Common Issues and Solutions

### Issue 1: Backend Not Running
**Symptoms**: Connection refused, 404 errors
**Solution**: Start backend with `go run cmd/myapp/main.go`

### Issue 2: Authentication Required
**Symptoms**: 401 Unauthorized errors
**Solution**: Get valid authentication token from browser

### Issue 3: File Upload Failing
**Symptoms**: No "Document photo upload started" log
**Solution**: Check form data, verify file is being sent

### Issue 4: Database Insertion Failing
**Symptoms**: No "Installation record created successfully" log
**Solution**: Check database connection, verify table schema

### Issue 5: Form Data Not Parsed
**Symptoms**: Empty document_photo in request data log
**Solution**: Check multipart form parsing, verify form field names

## Expected Results

### ✅ Successful Upload Should Show:
1. **Backend Logs**: All 5 log messages in sequence
2. **Database**: New record with document_photo field populated
3. **File System**: File exists in uploads/installations/documents/
4. **Frontend**: Success message and redirect

### ❌ Failed Upload Indicators:
1. **Missing Logs**: Not all 5 log messages appear
2. **Empty Database**: document_photo field is NULL or empty
3. **No File**: File not found in uploads directory
4. **Frontend Errors**: Error messages or failed submission

## Files Created for Testing

### Testing Scripts:
- ✅ `crm-be/debug-upload-issue.ps1` - General debugging
- ✅ `crm-be/test-upload-direct.ps1` - Direct upload testing
- ✅ `crm-be/test-frontend-form.ps1` - Frontend form testing
- ✅ `crm-be/check-database-simple.ps1` - Database verification

### Documentation:
- ✅ `crm-be/DOCUMENT_UPLOAD_TROUBLESHOOTING.md` - This guide

## Status:
✅ **Comprehensive testing scripts created** - All aspects covered
✅ **Database verification tools** - Schema and data checking
✅ **Backend testing tools** - Direct upload testing
✅ **Frontend testing guide** - Form submission testing
✅ **Troubleshooting documentation** - Complete guide

## Next Steps:
1. **Run backend**: `go run cmd/myapp/main.go`
2. **Test backend**: `.\debug-upload-issue.ps1`
3. **Test upload**: `.\test-upload-direct.ps1`
4. **Test frontend**: `.\test-frontend-form.ps1`
5. **Check database**: `.\check-database-simple.ps1`

**Follow these steps systematically to identify and fix the document upload issue! 🎉**
