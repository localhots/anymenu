package fonts

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// Find looks up font file location by the name.
func Find(name string) (string, error) {
	// Normalize font file name
	fileName := strings.Replace(name, " ", "-", -1) + ".ttf"

	for _, dir := range fontDirs() {
		dir, err := expandHome(dir)
		if err != nil {
			return "", err
		}

		path := filepath.Join(dir, fileName)
		if fileExist(path) {
			return path, nil
		}
	}

	return "", nil
}

func expandHome(path string) (string, error) {
	if !strings.HasPrefix(path, "~/") {
		return path, nil
	}

	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(u.HomeDir, path[2:]), nil
}

func fileExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
