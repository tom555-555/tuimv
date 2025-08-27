package model

import (
	"fmt"
	"os"
	"time"
)

// FileInfo represents a file or directory with its metadata
type FileInfo struct {
	Name        string      // File/directory name
	Path        string      // Relative path from current directory
	AbsPath     string      // Absolute path
	IsDirectory bool        // Whether this is a directory
	Size        int64       // File size in bytes
	ModTime     time.Time   // Last modification time
	IsSelected  bool        // Whether this file is selected for moving
	Permission  os.FileMode // File permission bits
}

// GetDisplayName returns a formatted name for display
func (f *FileInfo) GetDisplayName() string {
	if f.IsDirectory {
		return f.Name + "/"
	}
	return f.Name
}

// GetSizeString returns a human-readable size string
func (f *FileInfo) GetSizeString() string {
	if f.IsDirectory {
		return "<DIR>"
	}

	size := f.Size
	units := []string{"B", "KB", "MB", "GB", "TB"}
	unitIndex := 0

	for size >= 1024 && unitIndex < len(units)-1 {
		size /= 1024
		unitIndex++
	}

	if unitIndex == 0 {
		return fmt.Sprintf("%d %s", size, units[unitIndex])
	}
	return fmt.Sprintf("%.1f %s", float64(size), units[unitIndex])
}

// IsHidden returns true if the file/directory is hidden (starts with .)
func (f *FileInfo) IsHidden() bool {
	return len(f.Name) > 0 && f.Name[0] == '.'
}
