package hashbuffer

import (
	"crypto/sha256"
	"testing"

	"github.com/chmduquesne/rollinghash"
	_rabinkarp64 "github.com/chmduquesne/rollinghash/rabinkarp64"
)

// Test and provide basic examples of useage.

// Test hashBuffer using a rolling hash algorithm.
func TestRollingHash(t *testing.T) {
	// size of the file read buffer that HashBuffer should allocate
	const bufferSize = 1024
	// size of the word to be hashed
	const windowSize = 16
	// file to hash
	const filespec = "testdata/data_1025"

	// Get data
	hb, err := NewHashBuffer(filespec, bufferSize, windowSize)
	if err != nil {
		panic(err)
	}
	defer hb.Close()

	// create the Hash class
	var rolling rollinghash.Hash64
	rolling = _rabinkarp64.New()

	// read the first word from the file
	var window []byte
	window, err = hb.GetWindow()
	if err != nil {
		panic(err)
	}

	// This is an edge case that would normally be handled
	if len(window) == 0 {
		panic("Empty file")
	}

	// This is an edge case that would normally be handled
	if len(window) < windowSize {
		panic("File too small")
	}

	// initialize the rolling hash with the first word
	rolling.Write(window)

	// Continue reading though to the end of the file
	var nextByte byte
	for byteAvailable := true; byteAvailable; { // do...while
		// output the hash value for the current word
		t.Log(rolling.Sum64())
		nextByte, byteAvailable, err = hb.GetNext()
		if err != nil {
			panic(err)
		}
		// Roll the incoming byte in rolling
		if byteAvailable {
			rolling.Roll(nextByte)
		}
	}
}

// Test hashBuffer using a normal hash algorithm.
func TestStandardHash(t *testing.T) {
	// size of the file read buffer that HashBuffer should allocate
	const bufferSize = 1024
	// size of the word to be hashed
	const windowSize = 16
	// file to hash
	const filespec = "testdata/data_1025"

	// Get data
	hb, err := NewHashBuffer(filespec, bufferSize, windowSize)
	if err != nil {
		panic(err)
	}
	defer hb.Close()

	// Read though to the end of the file
	var window []byte
	for dataAvailable := true; dataAvailable; { // do...while
		// read the next available word from the file
		window, err = hb.GetWindow()
		if err != nil {
			panic(err)
		}

		if len(window) > 0 {
			// output the hash value for the current word
			t.Log(sha256.Sum256(window))
		} else {
			// we reached the end of the file
			dataAvailable = false
		}

	}
}
