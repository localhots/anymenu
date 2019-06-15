// +build !linux,!darwin

package fonts

func fontDirs() []string {
	return []string{
		".", // Current directory
	}
}