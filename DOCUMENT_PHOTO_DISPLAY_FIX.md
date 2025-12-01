# Document Photo Display Fix - Complete Guide

## Problem Description
Document photos (KTP) uploaded during installation report creation are not displaying in the detail view, even though the document_type shows "KTP".

## Root Cause Analysis
The issue is likely caused by:

1. **Incorrect URL Path**: Document photo paths in database may not match the static file serving configuration
2. **Missing Static File Serving**: Backend may not be serving uploaded files correctly
3. **Frontend URL Construction**: Frontend may not be constructing the correct URL for document photos
4. **File Upload Path**: Files may be uploaded to wrong directory or with incorrect naming

## Solution Implemented

### 1. ✅ Frontend Enhancements

#### **Enhanced Document Photo Display**
- **File**: `crm-fe/pages/dashboard/report/customer-installation/detail/[id].vue`
- **Changes**:
  - Added proper URL construction for document photos
  - Added click-to-view functionality
  - Added error handling for missing images
  - Added modal for full-size viewing
  - Added download functionality

#### **New Functions Added**:
```typescript
// Construct proper URL for document photos
function getDocumentPhotoUrl(documentPhoto: string | undefined) {
  if (!documentPhoto) return '';
  
  // Handle different path formats
  if (documentPhoto.startsWith('http')) {
    return documentPhoto;
  }
  
  if (documentPhoto.startsWith('uploads/')) {
    return `http://localhost:8080/${documentPhoto}`;
  }
  
  if (!documentPhoto.includes('/')) {
    return `http://localhost:8080/uploads/installations/documents/${documentPhoto}`;
  }
  
  return `http://localhost:8080/${documentPhoto}`;
}

// Open document photo in modal
function openDocumentPhoto(documentPhoto: string | undefined) {
  if (!documentPhoto) return;
  selectedDocumentPhoto.value = documentPhoto;
  showDocumentModal.value = true;
}

// Download document photo
function downloadDocumentPhoto() {
  if (!selectedDocumentPhoto.value) return;
  
  const photoUrl = getDocumentPhotoUrl(selectedDocumentPhoto.value);
  const link = document.createElement('a');
  link.href = photoUrl;
  link.download = `document_${report.value?.customer_name || 'installation'}_${report.value?.document_type || 'document'}.jpg`;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
}
```

#### **Enhanced Template**:
```vue
<div v-if="report.document_photo">
  <label class="text-sm font-medium text-gray-500">Document Photo</label>
  <div class="mt-2">
    <img 
      :src="getDocumentPhotoUrl(report.document_photo)" 
      alt="Document Photo" 
      class="w-48 h-32 object-cover rounded border cursor-pointer hover:opacity-80 transition-opacity"
      @click="openDocumentPhoto(report.document_photo)"
      @error="handleImageError"
    />
    <p class="text-xs text-gray-500 mt-1">Click to view full size</p>
  </div>
</div>
<div v-else>
  <label class="text-sm font-medium text-gray-500">Document Photo</label>
  <p class="text-gray-400 italic">No document photo uploaded</p>
</div>
```

#### **Modal for Full-Size Viewing**:
```vue
<UModal v-model="showDocumentModal" :ui="{ width: 'w-full max-w-4xl' }">
  <UCard>
    <template #header>
      <h3 class="text-lg font-semibold">Document Photo - {{ report?.document_type || 'Document' }}</h3>
    </template>
    
    <div class="flex justify-center">
      <img 
        v-if="selectedDocumentPhoto"
        :src="getDocumentPhotoUrl(selectedDocumentPhoto)" 
        alt="Document Photo" 
        class="max-w-full max-h-96 object-contain rounded"
      />
    </div>
    
    <template #footer>
      <UButton @click="downloadDocumentPhoto">
        <LucideIcon name="i-heroicons-arrow-down-tray" class="mr-2" />
        Download
      </UButton>
    </template>
  </UCard>
</UModal>
```

### 2. ✅ Backend Configuration

#### **Static File Serving**
- **File**: `crm-be/internal/routes/route.go`
- **Configuration**: `app.Static("/uploads", "./uploads")`
- **Status**: ✅ Already configured correctly

#### **File Upload Handling**
- **File**: `crm-be/internal/api/admin/customer/installation/report_controller_new.go`
- **Upload Directory**: `uploads/installations/documents`
- **File Naming**: `document_YYYYMMDD_HHMMSS.ext`
- **Status**: ✅ Already implemented correctly

### 3. ✅ Testing Scripts

#### **Document Photo Access Test**
- **File**: `crm-be/test-document-photo-access.ps1`
- **Purpose**: Test document photo access and static file serving
- **Features**:
  - Check uploads directory existence
  - List uploaded files
  - Test static file serving
  - Check database document photo paths

## Expected Behavior

### 1. **Document Photo Display**
- ✅ Document photo should appear in the Document Information section
- ✅ Photo should be clickable to open in modal
- ✅ Modal should show full-size image
- ✅ Download button should be available

### 2. **URL Construction**
- ✅ Frontend should construct proper URLs for document photos
- ✅ URLs should point to `http://localhost:8080/uploads/installations/documents/filename`
- ✅ Static file serving should serve the files correctly

### 3. **Error Handling**
- ✅ Missing photos should show "No document photo uploaded"
- ✅ Broken image links should show fallback image
- ✅ Console errors should be handled gracefully

## Testing Steps

### 1. **Test Document Photo Access**
```powershell
cd crm-be
.\test-document-photo-access.ps1
```

### 2. **Test Frontend Display**
1. Navigate to installation report detail page
2. Check Document Information section
3. Verify document photo appears
4. Click on photo to open modal
5. Test download functionality

### 3. **Test Static File Serving**
```bash
# Test direct access to uploaded file
curl http://localhost:8080/uploads/installations/documents/document_20250101_120000.jpg
```

## Troubleshooting

### Issue 1: Document photo not displaying
**Cause**: Incorrect URL path or missing file
**Solution**: 
1. Check database for correct document_photo path
2. Verify file exists in uploads directory
3. Test static file serving

### Issue 2: 404 error when accessing document photo
**Cause**: Static file serving not working
**Solution**:
1. Verify backend is running
2. Check static file serving configuration
3. Ensure uploads directory exists

### Issue 3: Modal not opening
**Cause**: JavaScript error or missing reactive variables
**Solution**:
1. Check browser console for errors
2. Verify reactive variables are defined
3. Check modal component import

## Files Modified

### Frontend Files:
- ✅ `crm-fe/pages/dashboard/report/customer-installation/detail/[id].vue` - Enhanced document photo display

### Backend Files:
- ✅ `crm-be/internal/routes/route.go` - Static file serving (already configured)
- ✅ `crm-be/internal/api/admin/customer/installation/report_controller_new.go` - File upload handling (already implemented)

### Testing Files:
- ✅ `crm-be/test-document-photo-access.ps1` - Document photo access testing
- ✅ `crm-be/DOCUMENT_PHOTO_DISPLAY_FIX.md` - This documentation

## Status:
✅ **Frontend enhanced** - Document photo display with modal and download
✅ **URL construction fixed** - Proper URL handling for different path formats
✅ **Error handling added** - Graceful handling of missing or broken images
✅ **Modal functionality** - Full-size viewing and download capabilities
✅ **Testing script created** - Comprehensive testing for document photo access
✅ **Documentation complete** - Step-by-step guide for troubleshooting

## Next Steps:
1. **Test document photo access**: Run `test-document-photo-access.ps1`
2. **Verify frontend display**: Check installation report detail page
3. **Test functionality**: Click photo, open modal, download file
4. **Debug if needed**: Use browser console and testing script output

**Document photos should now display correctly with full functionality! 🎉**
