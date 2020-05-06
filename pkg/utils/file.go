package utils

import "os"

// FileMode returns the file mode of path.
func FileMode(path string) (os.FileMode, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	return stat.Mode(), nil
}
