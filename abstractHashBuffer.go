package hashbuffer

import (
	"io"
	"testing"
)

/*
break of the file-ness of this into AbstracthashBuffer and FileHashBuffer

type io.Reader interface {
        Read(p []byte) (n int, err error)
}
type io.Closer interface {
        Close() error
}
*/

// abstractHashBuffer is a base class of HashBuffer.
type abstractHashBuffer struct {
	// reader *os.File
	// total size of the buffer
	bufferSize int
	// current index into buffer (0 based)
	pointer int
	// total number of available bytes in buffer
	fillLevel int
	// buffer where data is read into
	buffer []byte
	// remains true while stream has data left to read in
	isOpen bool
	// current size of the window (may be reduced at the last read)
	windowSize int

	reader io.Reader
	closer io.Closer

	// optional testing object to send progress information to
	t *testing.T // user-supplied; may be left nil
}

// Init initializes an abstractHashBuffer against the specified i/o, with the specified buffer size and window size.
func (ahb *abstractHashBuffer) init(reader io.Reader, closer io.Closer, bufferSize int, windowSize int) {
	ahb.reader = reader
	ahb.closer = closer
	// The buffer needs to be at least as big as the window size.
	if bufferSize < windowSize {
		ahb.bufferSize = windowSize
	} else {
		ahb.bufferSize = bufferSize
	}
	ahb.fillLevel = 0
	ahb.pointer = 0
	ahb.buffer = make([]byte, bufferSize)
	ahb.windowSize = windowSize
}

// GetWindow returns up to numberOfBytes of data as byte[], along with the number of bytes returned; if no bytes are available, return nil and 0.
func (ahb *abstractHashBuffer) GetWindow() (window []byte, err error) {
	// if ahb.isOpen {
	// If we need the first read or if the buffer is empty, attempt to read in more data.
	ahb.logf("GetWindow() starting;  bufferEmpty %v", ahb.bufferEmpty())
	if ahb.bufferEmpty() {
		err = ahb.fillBuffer()
		if err != nil {
			ahb.logf("GetWindow(): fillBuffer err %v", err)
			return
		}
		if ahb.bufferEmpty() {
			ahb.log("GetWindow(): out of data after an attempt to load")
			return
		}
	}
	ahb.logf("available window size: %d", ahb.windowSize)
	start := ahb.pointer
	end := ahb.pointer + ahb.windowSize
	window = ahb.buffer[start:end]
	// After getting start and end, advance the pointer.
	ahb.pointer++
	ahb.logf("start %d  end %d  len %d", start, end, len(ahb.buffer))
	ahb.logf("val %#x (%s)", window, string(window))
	// }
	return
}

// GetNext returns the next available byte of data if available and true; if not available return nil and false.
func (ahb *abstractHashBuffer) GetNext() (nextByte byte, byteAvailable bool, err error) {
	var window []byte
	window, err = ahb.GetWindow()
	bytesReceived := len(window)
	if bytesReceived > 0 {
		nextByte = window[bytesReceived-1]
		byteAvailable = true
		ahb.logf("GetNext returning from %d  len %d", ahb.pointer+ahb.windowSize-1, len(ahb.buffer))
	}
	return
}

// Skip skips over the next `count` bytes in the input stream.
func (ahb *abstractHashBuffer) Skip(count int) (numberSkipped int, err error) {
	// if ahb.isOpen {
	// determine if there is not enough in the buffer currently to skip over
	if (ahb.pointer + ahb.windowSize + count) > ahb.fillLevel {
		// attempt to fill buffer
		err = ahb.fillBuffer()
		if err != nil {
			ahb.logf("Skip(): initial fillBuffer err %v", err)
			return
		}
		if ahb.bufferEmpty() {
			ahb.log("Skip(): out of data after an attempt to load")
			return
		}
	}
	// calculate the amount available to skip
	amountAvailableToSkip := ahb.fillLevel - (ahb.pointer + ahb.windowSize)
	ahb.logf("amountAvailableToSkip=%d", amountAvailableToSkip)
	// reduce the amount to skip to the max amount we have available
	if amountAvailableToSkip <= 0 {
		numberSkipped = 0
		return
	}
	if amountAvailableToSkip < count {
		numberSkipped = amountAvailableToSkip
	} else {
		numberSkipped = count
	}
	ahb.logf("numberSkipped=%d", numberSkipped)
	// advance the pointer by that amount
	ahb.pointer = ahb.pointer + numberSkipped
	ahb.logf("recurse=%v", ((!ahb.bufferEmpty()) && (numberSkipped < count)))
	if (!ahb.bufferEmpty()) && (numberSkipped < count) {
		// attempt to fill buffer
		// err = ahb.fillBuffer()
		// if err != nil {
		// 	ahb.logf("Skip(): fillBuffer err %v", err)
		// 	return
		// }
		// make the call again until either we reach count,
		// or we run out of data
		childCount := count - numberSkipped
		var childSkipped int
		ahb.logf("childCount=%d", childCount)
		childSkipped, err = ahb.Skip(childCount)
		if err != nil {
			ahb.logf("Skip(): child Skip() err %v", err)
		}
		numberSkipped += childSkipped
	}
	// }
	return
}

// Close the stream if it is not already closed.
func (ahb *abstractHashBuffer) Close() (err error) {
	if ahb.isOpen {
		err := ahb.closer.Close
		if err == nil {
			ahb.isOpen = false
		}
	}
	return
}

// SetTesting allows for logging to be sent when testing HashBuffer.
func (ahb *abstractHashBuffer) SetTesting(t *testing.T) {
	ahb.t = t
}

func (ahb *abstractHashBuffer) fillBuffer() (err error) {
	if ahb.isOpen {
		// if we reloading the buffer, we need to save the current window and then continue loading
		if ahb.pointer != 0 {
			// move the window at the end, to the beginning of the buffer
			// read in as much as we can after that
			from := ahb.pointer
			to := ahb.fillLevel
			ahb.logf("Preparing buffer to be refilled  from %d (pointer):%d (fillLevel)  to 0  -  new fillLevel %d",
				from, to, (ahb.fillLevel - ahb.pointer))
			if to > from {
				copy(ahb.buffer[0:], ahb.buffer[from:to])
				ahb.fillLevel = ahb.fillLevel - ahb.pointer
				ahb.pointer = 0
				ahb.logf("new fillLevel %d", ahb.fillLevel)
			}
		}
		ahb.log("Filling buffer")
		// beginning just past fillLevel, fill as much of the buffer as we can
		var bytesread int
		bytesread, err = ahb.reader.Read(ahb.buffer[ahb.fillLevel:])
		if err != nil {
			if err != io.EOF {
				ahb.logf("Error %v, closing", err)
			} else {
				err = nil
				ahb.log("End of stream, closing")
			}
			ahb.Close()
		} else {
			// add the amount read to the fillLevel
			ahb.fillLevel += bytesread
			// if the whole stream has already been read and it is less than the window size, adjust the windowsize
			if ahb.fillLevel < ahb.windowSize {
				ahb.windowSize = ahb.fillLevel
			}
			// log amount read and the fillLevel
			ahb.logf("current fillLevel after read: %d  bytes read: %d\n",
				ahb.fillLevel, bytesread)
		}
	} else {
		ahb.log("File is not open.")
	}
	return
}

func (ahb *abstractHashBuffer) bufferEmpty() bool {
	ahb.logf("Calc bufferEmpty(): fillLevel %d  pointer %d  windowSize %d  LHS %d  RHS %d  bufferEmpty %v",
		ahb.fillLevel, ahb.pointer, ahb.windowSize,
		(ahb.pointer + ahb.windowSize),
		(ahb.fillLevel),
		(ahb.pointer+ahb.windowSize > ahb.fillLevel))
	return (ahb.pointer + ahb.windowSize) > ahb.fillLevel
}

func (ahb *abstractHashBuffer) log(message string) {
	if ahb.t != nil {
		ahb.t.Helper()
		ahb.t.Log(message)
	}
}
func (ahb *abstractHashBuffer) logf(format string, args ...interface{}) {
	if ahb.t != nil {
		ahb.t.Helper()
		ahb.t.Logf(format, args...)
		// fmt.Printf(format+"\n", args...)
	}
}
