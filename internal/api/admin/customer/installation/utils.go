package customerinstallation

import "strings"

// normalizeDocumentPhotoPath - Normalize document photo path by removing duplicated upload directories
func normalizeDocumentPhotoPath(path string) string {
	// Handle various path formats found in database:
	// 1. uploads\installations\documents\filename (Windows paths with backslashes)
	// 2. uploads/installations/documents/uploads/installations/documents/filename (duplicated)
	// 3. uploads/installations/documents/filename (correct)
	// 4. uploads/documents/filename (incorrect structure)

	// First, convert Windows backslashes to forward slashes for web URLs
	normalized := strings.ReplaceAll(path, "\\", "/")

	// Handle triple duplication: uploads/installations/documents/uploads/installations/documents/
	for strings.Contains(normalized, "uploads/installations/documents/uploads/installations/documents/") {
		normalized = strings.Replace(normalized, "uploads/installations/documents/uploads/installations/documents/", "uploads/installations/documents/", 1)
	}

	// Handle double duplication: uploads/installations/documents/uploads/installations/
	for strings.Contains(normalized, "uploads/installations/documents/uploads/installations/") {
		normalized = strings.Replace(normalized, "uploads/installations/documents/uploads/installations/", "uploads/installations/documents/", 1)
	}

	// Handle single duplication: uploads/installations/documents/uploads/
	if strings.HasPrefix(normalized, "uploads/installations/documents/uploads/") {
		normalized = strings.Replace(normalized, "uploads/installations/documents/uploads/", "uploads/installations/documents/", 1)
	}

	// Handle incorrect structure: uploads/documents/ -> uploads/installations/documents/
	if strings.HasPrefix(normalized, "uploads/documents/") {
		normalized = strings.Replace(normalized, "uploads/documents/", "uploads/installations/documents/", 1)
	}

	// Handle paths that are just filenames
	if !strings.Contains(normalized, "/") && strings.HasSuffix(normalized, ".jpg") {
		normalized = "uploads/installations/documents/" + normalized
	}

	return normalized
}
