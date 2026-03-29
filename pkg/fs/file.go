package fs

import (
	"io"
	"io/fs"
)

type (
	ReadableFile interface {
		fs.File
	}

	WritableFile interface {
		ReadableFile
		io.Writer
		io.Seeker
	}
)
