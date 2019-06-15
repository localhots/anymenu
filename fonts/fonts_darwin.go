// +build darwin

package fonts

// Docs: https://support.apple.com/en-us/HT201722
func fontDirs() []string {
	return []string{
		".",                     // Current directory
		"~/Library/Fonts",       // User
		"/Library/Fonts",        // Local
		"/System/Library/Fonts", // System
	}
}
