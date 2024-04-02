package common

import "path/filepath"

type FileMetadata struct {
	Filename    string
	Size        int64
	ContentHash []byte
}

func (m *FileMetadata) Read(srcFile string, hasher Hasher) {
	m.Filename = filepath.Base(srcFile)
	if hasher != nil {
		m.ContentHash = hasher(srcFile)
	}
	m.Size, _ = FileSize(srcFile)
}
