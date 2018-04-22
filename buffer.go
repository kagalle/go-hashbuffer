package rollhash

import (
	"io"
	"os"
	"testing"
)

// HashBuffer defines method to retrieve one or multiple bytes from a buffered stream of data.
type HashBuffer interface {
	Get(numberOfBytes int) ([]byte, int, error)
	GetNext() (byte, bool, error)
	Close() error
	SetTesting(t *testing.T)
}

// FileHashBuffer is a file based HashBuffer.
type FileHashBuffer struct {
	reader     *os.File
	bufferSize int
	pointer    int
	fillLevel  int
	buffer     []byte
	isOpen     bool
	t          *testing.T // may be nil
}

/*
// the buffer must be at least as big as the window size
if fhb.bufferSize < windowSize {
	bufferSize = windowSize
}
*/

// NewHashBuffer creates a FileHashBuffer against the specified filespec, with the specified buffersize.
func NewHashBuffer(filespec string, bufferSize int) (HashBuffer, error) {
	fhb := new(FileHashBuffer)
	fhb.bufferSize = bufferSize
	fhb.fillLevel = 0
	fhb.pointer = 0
	fhb.buffer = make([]byte, bufferSize)

	f, err := os.Open(filespec) // f : *os.File which implements io.Reader
	if err != nil {
		return nil, err
	}
	fhb.reader = f
	fhb.isOpen = true
	return fhb, nil
}

// SetTesting allows for logging to be sent when testing HashBuffer.
func (fhb *FileHashBuffer) SetTesting(t *testing.T) {
	fhb.t = t
}

// Close clses the file stream if it is not already closed.
func (fhb *FileHashBuffer) Close() error {
	if fhb.isOpen {
		err := fhb.reader.Close()
		if err != nil {
			return err
		}
		fhb.isOpen = false
	}
	return nil
}

// Get returns up to numberOfBytes of data as byte[], along with the number of bytes returned; if no bytes are available, return nil and 0.
func (fhb *FileHashBuffer) Get(numberOfBytes int) ([]byte, int, error) {
	if fhb.bufferEmpty() || fhb.bytesAvailable() < numberOfBytes {
		err := fhb.fillBuffer()
		if err != nil {
			return nil, 0, err
		}
	}
	numberToUse := numberOfBytes
	if fhb.bytesAvailable() < numberOfBytes {
		numberToUse = fhb.bytesAvailable()
	}
	fhb.logf("number to use %d", numberToUse)
	if numberToUse > 0 {
		start := fhb.pointer
		end := fhb.pointer + numberToUse
		fhb.pointer += numberToUse
		fhb.logf("start %d  end %d  val %#x", start, end, fhb.buffer[start:end])
		return fhb.buffer[start:end], numberToUse, nil
	}
	return nil, 0, nil
}

// GetNext returns the next available byte of data if available and true; if not available return nil and false.
func (fhb *FileHashBuffer) GetNext() (byte, bool, error) {
	if fhb.bufferEmpty() {
		err := fhb.fillBuffer()
		if err != nil {
			fhb.logf("Error %v", err)
			return 0, false, err
		}
	}
	var retval byte
	sent := false
	if !fhb.bufferEmpty() {
		retval = fhb.buffer[fhb.pointer]
		sent = true
		fhb.pointer++
	}
	return retval, sent, nil
}

func (fhb *FileHashBuffer) fillBuffer() error {
	if fhb.isOpen {
		fhb.log("Filling buffer")
		// if we've read all of the buffer, then reset the pointer back to zero
		if fhb.bufferEmpty() {
			fhb.pointer = 0
			fhb.fillLevel = 0
		}
		// beginning at the pointer, begin reading to fill as much of the buffer as we can
		bytesread, err := fhb.reader.Read(fhb.buffer[fhb.fillLevel:]) // reads up to len(buffer) bytes
		if err != nil {
			if err != io.EOF {
				fhb.logf("Error %v", err)
				return err
			} else {
				fhb.log("End of file, closing.")
			}
			fhb.Close()
		} else {
			// add the amount read to the fillLevel
			fhb.fillLevel += bytesread

			// log amount read and the fillLevel
			fhb.logf("current fillLevel after read: %d  bytes read: %d\n",
				fhb.fillLevel, bytesread)
		}
	} else {
		fhb.log("File is not open.")
	}
	return nil
}

func (fhb *FileHashBuffer) bufferEmpty() (isEmpty bool) {
	return fhb.fillLevel == fhb.pointer
}

func (fhb *FileHashBuffer) bytesAvailable() (amt int) {
	return fhb.fillLevel - fhb.pointer
}

func (fhb *FileHashBuffer) log(message string) {
	if fhb.t != nil {
		fhb.t.Log(message)
	}
}
func (fhb *FileHashBuffer) logf(format string, args ...interface{}) {
	if fhb.t != nil {
		fhb.t.Logf(format, args...)
	}
}
