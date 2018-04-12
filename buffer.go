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
}

type FileHashBuffer struct {
	reader     io.Reader
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
	defer f.Close()
	fhb.isOpen = true
	return fhb
}

// return byte[], number of bytes returned
func (fhb *FileHashBuffer) Get(numberOfBytes int) ([]byte, int) {
	if fhb.bytesAvailable() < numberOfBytes {
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
		// if we've read all of the buffer, then reset the pointer back to zero
		if fhb.bufferEmpty() {
			fhb.pointer = 0
			fhb.fillLevel = 0
		}
		// beginning at the pointer, begin reading to fill as much of the buffer as we can
		bytesread, err := fhb.reader.Read(fhb.buffer[fhb.fillLevel:]) // reads up to len(buffer) bytes
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			fhb.isOpen = false
			fmt.Println("End of file, closing.")
		} else {
			// add the amount read to the fillLevel
			fhb.fillLevel += bytesread - 1

			// log amount read and the fillLevel
			fmt.Printf("current fillLevel after read: %d  bytes read: %d\n",
				fhb.fillLevel, bytesread)
		}
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
		log.Fatal(e)
		fhb.isOpen = false
		// panic(e)
	}
}
