package hashbuffer

import (
	"io"
	"os"
	"testing"
)

// fileHashBuffer is a file based HashBuffer.
type fileHashBuffer struct {
	reader     *os.File
	bufferSize int
	pointer    int
	fillLevel  int
	buffer     []byte
	isOpen     bool
	// window holds the current window content
	windowSize int
	//window      []byte
	windowReady bool // true after first read
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
	fhb.fillLevel = -1
	fhb.pointer = 0
	fhb.buffer = make([]byte, bufferSize)
	fhb.windowSize = windowSize
	//fhb.window = make([]byte, windowSize)
	fhb.windowReady = false

	f, err := os.Open(filespec) // f : *os.File which implements io.Reader
	if err != nil {
		return nil, err
	}
	fhb.reader = f
	fhb.isOpen = true
	hashBuffer = fhb
	return hashBuffer, nil
}

// Get returns up to numberOfBytes of data as byte[], along with the number of bytes returned; if no bytes are available, return nil and 0.
func (fhb *fileHashBuffer) GetWindow() (window []byte, windowSize int, err error) {
	// If we need the first read or if the buffer is empty, attempt to read in more data.
	fhb.logf("windowReady %v  bufferEmpty %v", fhb.windowReady, fhb.bufferEmpty())
	if (!fhb.windowReady) || fhb.bufferEmpty() {
		err = fhb.fillBuffer()
		if err != nil {
			return nil, 0, err
		}
	}
	if (!fhb.windowReady) || fhb.bufferEmpty() {
		fhb.log("out of data after an attempt to load")
		return nil, 0, nil
	}
	// else...
	fhb.logf("available window size: %d", fhb.windowSize)
	start := fhb.pointer
	end := fhb.pointer + fhb.windowSize
	// After getting start and end, advance the pointer.
	fhb.pointer++
	fhb.logf("start %d  end %d  len %d", start, end, len(fhb.buffer))
	fhb.logf("val %#x", fhb.buffer[start:end])
	return fhb.buffer[start:end], fhb.windowSize, nil
}

// GetNext returns the next available byte of data if available and true; if not available return nil and false.
func (fhb *fileHashBuffer) GetNext() (nextByte byte, byteAvailable bool, err error) {
	var buffer []byte
	var bytesReceived int
	buffer, bytesReceived, err = fhb.GetWindow()
	if bytesReceived > 0 {
		return buffer[fhb.windowSize-1], true, nil
	} else {
		return 0, false, err
	}
}

// Close the file stream if it is not already closed.
func (fhb *fileHashBuffer) Close() error {
	if fhb.isOpen {
		err := fhb.reader.Close()
		if err != nil {
			return err
		}
		fhb.isOpen = false
	}
	return nil
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
		if fhb.windowReady {
			// move the window at the end, less the first character, to the beginning of the buffer
			// read in as much as we can after that
			fhb.log("Preparing buffer to be refilled")
			// move buffer at pointer + 1, windowsize - 1 to 0
			// pointer is 0
			copy(fhb.buffer[0:], fhb.buffer[fhb.pointer+1:fhb.pointer+fhb.windowSize-1])
			fhb.pointer = 0
			fhb.fillLevel = fhb.windowSize - 1 // -1 zero based, -1 drop one char of window ???
		}
		fhb.log("Filling buffer")
		// beginning at the pointer, begin reading to fill as much of the buffer as we can
		var bytesread int
		bytesread, err = fhb.reader.Read(fhb.buffer[fhb.fillLevel:]) // reads up to len(buffer) bytes
		if err != nil {
			if err != io.EOF {
				fhb.logf("Error %v", err)
				return
			}
			// else
			err = nil
			fhb.log("End of file, closing:")
			fhb.Close()
		} else {
			// add the amount read to the fillLevel
			fhb.fillLevel += bytesread
			// set windowReady
			if !fhb.windowReady {
				fhb.windowReady = true
				// if the whole file has already been read and it is less than the window size, adjust the windowsize
				if fhb.fillLevel < (fhb.windowSize - 1) {
					fhb.windowSize = fhb.fillLevel + 1
				}
				// w w w w w w
				//         f
				// 0 1 2 3 4 5 6 7 8
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

//         f
// p
// w w w w w w
// 0 1 2 3 4 5 6 7 8 9
func (fhb *fileHashBuffer) bufferEmpty() (isEmpty bool) {
	fhb.logf("fillLevel %d  pointer %d  windowSize %d  RHS %d  bufferEmpty %v",
		fhb.fillLevel, fhb.pointer, fhb.windowSize,
		(fhb.pointer + fhb.windowSize),
		(fhb.fillLevel < (fhb.pointer + fhb.windowSize)))
	return fhb.fillLevel < (fhb.pointer + fhb.windowSize)
}

//  0 < (0 + w - 1)

// func (fhb *fileHashBuffer) bytesAvailable() (amt int) {
// 	return fhb.fillLevel - fhb.pointer
// }

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
