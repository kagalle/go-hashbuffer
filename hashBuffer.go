package hashbuffer

import (
	"testing"
)

/*
 * Implementations:
 *   fileHashBuffer.go : NewHashBuffer(filespec string, bufferSize int) (HashBuffer, error)
 *
 */

// HashBuffer defines method to retrieve one or multiple bytes from a buffered stream of data.
type HashBuffer interface {
	// Get one window of data; each call moves data forward by one byte.
	// Param []byte: buffer of window
	// Param int: width of buffer returned; will be windowSize unless amount in file was less than windowSize
	// Param error: non-nil if an error occurred trying to read (something other than EOF).
	GetWindow() ([]byte, error)
	// Get next available byte of data; push this byte into the window.
	// This is equivelant to calling GetWindow() and using the right-most byte returned.
	// This is meant for rolling-hash algorithms that take an initial buffer of data and then additional bytes are added in.
	GetNext() (byte, bool, error)
	// Close the file handle
	Close() error
	// Send testing object in to which HashBuffer will write information on its progress
	SetTesting(t *testing.T)
}

/*
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
*/
