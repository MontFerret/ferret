package file

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRead(t *testing.T) {
	Convey("Read function", t, func() {
		Convey("Should read existing file successfully", func() {
			// Create a temporary file
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "test.txt")
			testContent := "hello world\nline 2"
			
			err := os.WriteFile(testFile, []byte(testContent), 0644)
			So(err, ShouldBeNil)
			
			// Test reading the file
			source, err := Read(testFile)
			
			So(err, ShouldBeNil)
			So(source, ShouldNotBeNil)
			So(source.Name(), ShouldEqual, testFile)
			So(source.Content(), ShouldEqual, testContent)
			So(source.Empty(), ShouldBeFalse)
		})
		
		Convey("Should handle non-existent file", func() {
			nonExistentFile := "/path/that/does/not/exist/file.txt"
			
			source, err := Read(nonExistentFile)
			
			So(err, ShouldNotBeNil)
			So(source, ShouldBeNil)
		})
		
		Convey("Should handle empty file", func() {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "empty.txt")
			
			err := os.WriteFile(testFile, []byte(""), 0644)
			So(err, ShouldBeNil)
			
			source, err := Read(testFile)
			
			So(err, ShouldBeNil)
			So(source, ShouldNotBeNil)
			So(source.Name(), ShouldEqual, testFile)
			So(source.Content(), ShouldEqual, "")
			So(source.Empty(), ShouldBeTrue)
		})
		
		Convey("Should handle file with special characters", func() {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "special.txt")
			testContent := "Hello 世界\n\tTab\r\nWindows line ending"
			
			err := os.WriteFile(testFile, []byte(testContent), 0644)
			So(err, ShouldBeNil)
			
			source, err := Read(testFile)
			
			So(err, ShouldBeNil)
			So(source, ShouldNotBeNil)
			So(source.Content(), ShouldEqual, testContent)
		})
		
		Convey("Should handle directory instead of file", func() {
			tmpDir := t.TempDir()
			
			source, err := Read(tmpDir)
			
			So(err, ShouldNotBeNil)
			So(source, ShouldBeNil)
		})
	})
}