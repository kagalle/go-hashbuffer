package hashbuffer

import (
	"testing"
)

/*
 * Implementations:
 *   fileHashBuffer.go : NewFileHashBuffer(filespec string, bufferSize int) (HashBuffer, error)
 *
 */

// HashBuffer defines method to retrieve one or multiple bytes from a buffered stream of data.
type HashBuffer interface {
	// Get one window of data; each call moves data forward by one byte.
	// Param []byte: buffer of window
	// Param error: non-nil if an error occurred trying to read (something other than EOF).
	GetWindow() (window []byte, err error)
	// Get next available byte of data; push this byte into the window.
	// This is equivelant to calling GetWindow() and using the right-most byte returned.
	// This is meant for rolling-hash algorithms that take an initial buffer of data and
	// then additional bytes are added in.
	GetNext() (nextByte byte, byteAvailable bool, err error)
	// Skip over the next `count` bytes in the input stream.  This is equivelant to calling
	// `GetNext()` `count` times, and discarding the results.
	// Returns the number actually skipped (less than `count` if EOF is reached).
	Skip(count int) (numberSkipped int, err error)
	// Close the file handle
	Close() (err error)
	// Send testing object in to which HashBuffer will write information on its progress
	SetTesting(t *testing.T)
}
