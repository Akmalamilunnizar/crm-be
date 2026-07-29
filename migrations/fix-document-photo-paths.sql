-- Fix Document Photo Paths in Database
-- This script normalizes document photo paths in the customer_installations table

-- First, let's see what paths currently exist
SELECT
    id,
    document_type,
    document_photo,
    CASE
        WHEN document_photo LIKE '%uploads/installations/documents/uploads/installations/documents/%' THEN 'TRIPLE_DUPLICATED'
        WHEN document_photo LIKE '%uploads/installations/documents/uploads/installations/%' THEN 'DOUBLE_DUPLICATED'
        WHEN document_photo LIKE '%uploads/installations/documents/uploads/%' THEN 'SINGLE_DUPLICATED'
        WHEN document_photo LIKE '%uploads/documents/%' THEN 'WRONG_STRUCTURE'
        WHEN document_photo LIKE '%uploads/installations/documents/%' THEN 'CORRECT'
        ELSE 'OTHER'
    END as path_status,
    CASE
        WHEN document_photo LIKE '%uploads/installations/documents/uploads/installations/documents/%'
            THEN REPLACE(document_photo, 'uploads/installations/documents/uploads/installations/documents/', 'uploads/installations/documents/')
        WHEN document_photo LIKE '%uploads/installations/documents/uploads/installations/%'
            THEN REPLACE(document_photo, 'uploads/installations/documents/uploads/installations/', 'uploads/installations/documents/')
        WHEN document_photo LIKE '%uploads/installations/documents/uploads/%'
            THEN REPLACE(document_photo, 'uploads/installations/documents/uploads/', 'uploads/installations/documents/')
        WHEN document_photo LIKE '%uploads/documents/%'
            THEN REPLACE(document_photo, 'uploads/documents/', 'uploads/installations/documents/')
        ELSE document_photo
    END as normalized_path
FROM customer_installations
WHERE document_photo IS NOT NULL
AND document_photo != ''
ORDER BY createdAt DESC;

-- Update the paths (uncomment the line below after reviewing the results above)
/*
UPDATE customer_installations
SET document_photo = CASE
    WHEN document_photo LIKE '%uploads/installations/documents/uploads/installations/documents/%'
        THEN REPLACE(document_photo, 'uploads/installations/documents/uploads/installations/documents/', 'uploads/installations/documents/')
    WHEN document_photo LIKE '%uploads/installations/documents/uploads/installations/%'
        THEN REPLACE(document_photo, 'uploads/installations/documents/uploads/installations/', 'uploads/installations/documents/')
    WHEN document_photo LIKE '%uploads/installations/documents/uploads/%'
        THEN REPLACE(document_photo, 'uploads/installations/documents/uploads/', 'uploads/installations/documents/')
    WHEN document_photo LIKE '%uploads/documents/%'
        THEN REPLACE(document_photo, 'uploads/documents/', 'uploads/installations/documents/')
    ELSE document_photo
END
WHERE document_photo IS NOT NULL
AND document_photo != ''
AND (
    document_photo LIKE '%uploads/installations/documents/uploads/installations/documents/%' OR
    document_photo LIKE '%uploads/installations/documents/uploads/installations/%' OR
    document_photo LIKE '%uploads/installations/documents/uploads/%' OR
    document_photo LIKE '%uploads/documents/%'
);
*/

-- Verify the fix worked
-- SELECT id, document_type, document_photo FROM customer_installations WHERE document_photo IS NOT NULL ORDER BY createdAt DESC LIMIT 5;

