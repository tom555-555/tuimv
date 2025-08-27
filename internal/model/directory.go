package model

// DirectoryInfo represents a directory for fuzzy search results
type DirectoryInfo struct {
	Name    string // Directory name (basename)
	Path    string // Relative or user-friendly path
	AbsPath string // Absolute path
	Score   int    // Fuzzy search match score (higher is better)
}

// GetDisplayPath returns a formatted path for display
func (d *DirectoryInfo) GetDisplayPath() string {
	if len(d.Path) > 0 {
		return d.Path
	}
	return d.AbsPath
}

// GetDisplayName returns the directory name with a trailing slash
func (d *DirectoryInfo) GetDisplayName() string {
	return d.Name + "/"
}
