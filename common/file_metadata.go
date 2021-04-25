package common

import "path"

type FileMetadata struct {
	Filename    string
	Size        int64
	ContentHash []byte
}

func (m *FileMetadata) Read(srcFile string) {
	m.Filename = path.Base(srcFile)
	m.ContentHash = Sha256File(srcFile)
	m.Size, _ = FileSize(srcFile)
}
