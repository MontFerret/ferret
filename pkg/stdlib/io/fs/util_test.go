package fs_test

import (
	"os"

	. "github.com/smartystreets/goconvey/convey"
)

func tempFile() (*os.File, func()) {
	file, err := os.CreateTemp("", "fstest")
	So(err, ShouldBeNil)

	fn := func() {
		file.Close()
		os.Remove(file.Name())
	}

	return file, fn
}
