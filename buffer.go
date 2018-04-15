package rollhash

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"
)

// HashBuffer defines method to retrieve one or multiple bytes from a buffered stream of data.
type HashBuffer interface {
	Get(numberOfBytes int) ([]byte, int)
	GetNext() (byte, bool)
	Close()
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
func NewHashBuffer(filespec string, bufferSize int) HashBuffer {
	fhb := new(FileHashBuffer)
	fhb.bufferSize = bufferSize
	fhb.fillLevel = 0
	fhb.pointer = 0
	fhb.buffer = make([]byte, bufferSize)

	// s, err := ioutil.ReadFile("/home/ken/2002-06-19.pdf")
	f, err := os.Open(filespec) // f : *os.File which implements io.Reader
	fhb.reader = f
	fhb.check(err)
	fhb.isOpen = true
	return fhb
}

// SetTesting allows for logging to be sent when testing HashBuffer.
func (fhb *FileHashBuffer) SetTesting(t *testing.T) {
	fhb.t = t
}

// Close clses the file stream if it is not already closed.
func (fhb *FileHashBuffer) Close() {
	if fhb.isOpen {
		fhb.log("Closing")
		err := fhb.reader.Close()
		fhb.check(err)
		fhb.isOpen = false
	}
}

// Get returns up to numberOfBytes of data as byte[], along with the number of bytes returned; if no bytes are available, return nil and 0.
func (fhb *FileHashBuffer) Get(numberOfBytes int) ([]byte, int) {
	if fhb.bufferEmpty() || fhb.bytesAvailable() < numberOfBytes {
		fhb.fillBuffer()
	}
	numberToUse := numberOfBytes
	if fhb.bytesAvailable() < numberOfBytes {
		numberToUse = fhb.bytesAvailable()
	}
	fhb.log(fmt.Sprintf("number to use %d", numberToUse))
	if numberToUse > 0 {
		start := fhb.pointer
		end := fhb.pointer + numberToUse
		fhb.pointer += numberToUse
		return fhb.buffer[start:end], numberToUse
	} else {
		return nil, 0
	}
}

// GetNext returns the next available byte of data if available and true; if not available return nil and false.
func (fhb *FileHashBuffer) GetNext() (byte, bool) {
	if fhb.bufferEmpty() {
		fhb.fillBuffer()
	}
	var retval byte
	sent := false
	if !fhb.bufferEmpty() {
		retval = fhb.buffer[fhb.pointer]
		sent = true
	}
	return retval, sent
}

func (fhb *FileHashBuffer) fillBuffer() {
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
				fhb.check(err)
			} else {
				fhb.log("End of file, closing.")
			}
			fhb.Close()
		} else {
			// add the amount read to the fillLevel
			fhb.fillLevel += bytesread - 1

			// log amount read and the fillLevel
			fhb.log(fmt.Sprintf("current fillLevel after read: %d  bytes read: %d\n",
				fhb.fillLevel, bytesread))
		}
	} else {
		fhb.log("File is not open.")
	}
}

func (fhb *FileHashBuffer) bufferEmpty() (isEmpty bool) {
	return fhb.fillLevel == fhb.pointer
}

func (fhb *FileHashBuffer) bytesAvailable() (amt int) {
	return fhb.fillLevel - fhb.pointer
}

func (fhb *FileHashBuffer) check(e error) {
	if e != nil {
		log.Printf("Error %v", e)
		fhb.log(fmt.Sprintf("Error %v", e))
		panic(e)
	}
}

func (fhb *FileHashBuffer) log(message string) {
	if fhb.t != nil {
		fhb.t.Log(message)
	}
}
