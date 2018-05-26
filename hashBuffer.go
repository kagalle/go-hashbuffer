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
	Get(numberOfBytes int) ([]byte, int, error)
	GetNext() (byte, bool, error)
	Close() error
	SetTesting(t *testing.T)
}
