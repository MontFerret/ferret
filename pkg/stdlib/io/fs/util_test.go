package fs_test

import (
	"io/ioutil"
	"os"

	. "github.com/smartystreets/goconvey/convey"
)

func tempFile() (*os.File, func()) {
	file, err := ioutil.TempFile("", "fstest")
	So(err, ShouldBeNil)

	fn := func() {
		file.Close()
		os.Remove(file.Name())
	}

	return file, fn
}
