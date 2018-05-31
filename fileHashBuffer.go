package hashbuffer

import (
	"io"
	"os"
	"testing"
)

// fileHashBuffer is a file based HashBuffer.
type fileHashBuffer struct {
	reader *os.File
	// total size of the buffer
	bufferSize int
	// current index into buffer (0 based)
	pointer int
	// total number of available bytes in buffer
	fillLevel int
	// buffer where file data is read into
	buffer []byte
	// remains true while file has data left to read in
	isOpen bool
	// current size of the window (may be reduced at the last read)
	windowSize int
	// optional testing object to send progress information to
	t *testing.T // user-supplied; may be left nil
}

// NewHashBuffer creates a FileHashBuffer against the specified filespec, with the specified buffersize.
func NewHashBuffer(filespec string, bufferSize int, windowSize int) (hashBuffer HashBuffer, err error) {
	fhb := new(fileHashBuffer)
	// The buffer needs to be at least as big as the window size.
	if bufferSize < windowSize {
		fhb.bufferSize = windowSize
	} else {
		fhb.bufferSize = bufferSize
	}
	fhb.fillLevel = 0
	fhb.pointer = 0
	fhb.buffer = make([]byte, bufferSize)
	fhb.windowSize = windowSize

	f, err := os.Open(filespec) // f : *os.File which implements io.Reader
	if err != nil {
		return
	}
	fhb.reader = f
	fhb.isOpen = true
	hashBuffer = fhb
	return
}

// Get returns up to numberOfBytes of data as byte[], along with the number of bytes returned; if no bytes are available, return nil and 0.
func (fhb *fileHashBuffer) GetWindow() (window []byte, err error) {
	// If we need the first read or if the buffer is empty, attempt to read in more data.
	fhb.logf("GetWindow() starting;  bufferEmpty %v", fhb.bufferEmpty())
	if fhb.bufferEmpty() {
		err = fhb.fillBuffer()
		if err != nil {
			fhb.logf("fillBuffer err %v", err)
			return
		}
		if fhb.bufferEmpty() {
			fhb.log("out of data after an attempt to load")
			return
		}
	}
	fhb.logf("available window size: %d", fhb.windowSize)
	start := fhb.pointer
	end := fhb.pointer + fhb.windowSize
	window = fhb.buffer[start:end]
	// After getting start and end, advance the pointer.
	fhb.pointer++
	fhb.logf("start %d  end %d  len %d", start, end, len(fhb.buffer))
	fhb.logf("val %#x", window)
	return
}

// GetNext returns the next available byte of data if available and true; if not available return nil and false.
func (fhb *fileHashBuffer) GetNext() (nextByte byte, byteAvailable bool, err error) {
	var window []byte
	window, err = fhb.GetWindow()
	bytesReceived := len(window)
	if bytesReceived > 0 {
		nextByte = window[bytesReceived-1]
		byteAvailable = true
		fhb.logf("GetNext returning from %d  len %d", fhb.pointer+fhb.windowSize-1, len(fhb.buffer))
	}
	return
}

// Close the file stream if it is not already closed.
func (fhb *fileHashBuffer) Close() (err error) {
	if fhb.isOpen {
		err := fhb.reader.Close()
		if err == nil {
			fhb.isOpen = false
		}
	}
	return
}

// SetTesting allows for logging to be sent when testing HashBuffer.
func (fhb *fileHashBuffer) SetTesting(t *testing.T) {
	fhb.t = t
}

// Since we know we are reading and using one byte at a time after the initial read,
// we can assume that when this gets called, we can clear and reload the entire buffer.
func (fhb *fileHashBuffer) fillBuffer() (err error) {
	if fhb.isOpen {
		// if we reloading the buffer, we need to save the current window and then continue loading
		if fhb.pointer != 0 {
			// move the window at the end, less the first character, to the beginning of the buffer
			// read in as much as we can after that
			from := fhb.pointer + 1
			to := fhb.pointer + fhb.windowSize - 1
			fhb.logf("Preparing buffer to be refilled  from %d:%d to 0", from, to)
			if to > from {
				copy(fhb.buffer[0:], fhb.buffer[fhb.pointer+1:fhb.pointer+fhb.windowSize-1])
				fhb.pointer = 0
				fhb.fillLevel = fhb.windowSize - 1 // drop one char of window
			}
		}
		fhb.log("Filling buffer")
		// beginning just past fillLevel, fill as much of the buffer as we can
		var bytesread int
		bytesread, err = fhb.reader.Read(fhb.buffer[fhb.fillLevel:])
		if err != nil {
			if err != io.EOF {
				fhb.logf("Error %v", err)
			} else {
				err = nil
				fhb.log("End of file, closing:")
				fhb.Close()
			}
		} else {
			// add the amount read to the fillLevel
			fhb.fillLevel += bytesread
			// if the whole file has already been read and it is less than the window size, adjust the windowsize
			if fhb.fillLevel < fhb.windowSize {
				fhb.windowSize = fhb.fillLevel
			}
			// log amount read and the fillLevel
			fhb.logf("current fillLevel after read: %d  bytes read: %d\n",
				fhb.fillLevel, bytesread)
		}
	} else {
		fhb.log("File is not open.")
	}
	return
}

func (fhb *fileHashBuffer) bufferEmpty() bool {
	fhb.logf("Calc bufferEmpty(): fillLevel %d  pointer %d  windowSize %d  LHS %d  RHS %d  bufferEmpty %v",
		fhb.fillLevel, fhb.pointer, fhb.windowSize,
		(fhb.pointer + fhb.windowSize),
		(fhb.fillLevel),
		(fhb.pointer+fhb.windowSize > fhb.fillLevel))
	return (fhb.pointer + fhb.windowSize) > fhb.fillLevel
}

func (fhb *fileHashBuffer) log(message string) {
	if fhb.t != nil {
		fhb.t.Log(message)
	}
}
func (fhb *fileHashBuffer) logf(format string, args ...interface{}) {
	if fhb.t != nil {
		fhb.t.Logf(format, args...)
	}
}
