package rollhash

import (
	"fmt"
	"io"
	"log"
	"os"
)

type HashBuffer interface {
	Get(numberOfBytes int) ([]byte, int)
	GetNext() (byte, bool)
	Close()
}

type FileHashBuffer struct {
	reader     *os.File
	bufferSize int
	pointer    int
	fillLevel  int
	buffer     []byte
	isOpen     bool
}

/*
// the buffer must be at least as big as the window size
if fhb.bufferSize < windowSize {
	bufferSize = windowSize
}
*/

// bufferSize := 1024
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

func (fhb *FileHashBuffer) Close() {
	if fhb.isOpen {
		fmt.Println("Closing")
		err := fhb.reader.Close()
		fhb.check(err)
		fhb.isOpen = false
	}
}

// return byte[], number of bytes returned
func (fhb *FileHashBuffer) Get(numberOfBytes int) ([]byte, int) {
	if fhb.bufferEmpty() || fhb.bytesAvailable() < numberOfBytes {
		fhb.fillBuffer()
	}
	numberRead := numberOfBytes
	if fhb.bytesAvailable() < numberOfBytes {
		numberRead = fhb.bytesAvailable()
	}
	start := fhb.pointer
	end := fhb.pointer + numberOfBytes
	fhb.pointer += numberRead
	return fhb.buffer[start:end], numberRead
}

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
		fmt.Println("Filling buffer")
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
				fmt.Println("End of file, closing.")
			}
			fhb.Close()
		} else {
			// add the amount read to the fillLevel
			fhb.fillLevel += bytesread - 1

			// log amount read and the fillLevel
			fmt.Printf("current fillLevel after read: %d  bytes read: %d\n",
				fhb.fillLevel, bytesread)
		}
	} else {
		fmt.Println("File is not open.")
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
		panic(e)
	}
}
