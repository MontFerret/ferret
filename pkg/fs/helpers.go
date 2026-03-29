package fs

import "os"

func isWriteOpen(flag int) bool {
	const writeMask = os.O_WRONLY | os.O_RDWR | os.O_APPEND | os.O_CREATE | os.O_TRUNC
	return flag&writeMask != 0
}
