package customerinstallation

import "strings"

// normalizeDocumentPhotoPath - Normalize document photo path by removing duplicated upload directories
func normalizeDocumentPhotoPath(path string) string {
	// Replace multiple occurrences of the pattern
	for strings.Contains(path, "uploads/installations/documents/uploads/installations/documents/") {
		path = strings.Replace(path, "uploads/installations/documents/uploads/installations/documents/", "uploads/installations/documents/", 1)
	}

	// Also handle cases with different duplication patterns
	for strings.Contains(path, "uploads/installations/documents/uploads/installations/") {
		path = strings.Replace(path, "uploads/installations/documents/uploads/installations/", "uploads/installations/documents/", 1)
	}

	// Handle cases where it starts with uploads/installations/documents/uploads/
	if strings.HasPrefix(path, "uploads/installations/documents/uploads/") {
		path = strings.Replace(path, "uploads/installations/documents/uploads/", "uploads/installations/documents/", 1)
	}

	return path
}
