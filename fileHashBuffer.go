package hashbuffer

import (
	"os"
)

// fileHashBuffer is a file based HashBuffer.
type fileHashBuffer struct {
	*abstractHashBuffer
	file *os.File
}

// NewFileHashBuffer creates a FileHashBuffer against the specified filespec, with the specified buffersize.
func NewFileHashBuffer(filespec string, bufferSize int, windowSize int) (hashBuffer HashBuffer, err error) {
	fhb := new(fileHashBuffer)
	hashBuffer = fhb
	fhb.abstractHashBuffer = new(abstractHashBuffer)

	f, err := os.Open(filespec) // f : *os.File which implements io.Reader
	if err != nil {
		return
	}
	fhb.abstractHashBuffer.isOpen = true
	fhb.abstractHashBuffer.init(f, f, bufferSize, windowSize)
	return
}
