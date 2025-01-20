package types

import (
	"path/filepath"
	"strings"
)

func FilenameWithoutExt(filename string) string {
	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}
