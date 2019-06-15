// +build linux

package fonts

func fontDirs() []string {
	return []string{
		".",                             // Current directory
		"~/.fonts",                      // User
		"~/.fonts/truetype",             // User
		"~/.local/share/fonts",          // User
		"~/.local/share/fonts/truetype", // User
		"/usr/share/fonts",              // System
		"/usr/share/fonts/truetype",     // System
	}
}
