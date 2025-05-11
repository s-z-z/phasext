package util

import (
	"path/filepath"
)

func GetAbsolutePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	path, _ = filepath.Abs(path)
	return path
}
