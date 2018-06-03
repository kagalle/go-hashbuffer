# hashbuffer

Hashbuffer is a buffer designed to feed data (normally a file) to a hashing algorithm in a way that the input is broken up into fixed sized windows (also referred to as words), each of which is later hashed.  The output window starts at the beginning and moves forward one byte at a time.  If the windowSize is 16, then the first window returned would be bytes 0-15, the second, bytes 1-16, etc.

This can be used with a rolling hash, such as can be found at [chmduquesne/rollinghash](http://github.com/chmduquesne/rollinghash).  It can also be used with a standard hash, such as MD-5, etc.  Note that more complex hash algorithms generate wider hashes; there is little point in creating a hash of something when the checksum created is larger than the input.

A rolling hash is an algorithm that generates checksums of a sliding window over a long array of bytes. It does so more efficiently than if each block was checksumed individually, building the next checksum from the value of the previous one.  Typically, an initial word of data is used to create the first hash, and then additional single bytes of data for each additional hash generated.  For example:

```go
hb, err := NewHashBuffer(filespec, bufferSize, windowSize)
window, err = hb.GetWindow()
// initialize hash algorithm using `window`
// do...
    // use generated checksum
    nextByte, byteAvailable, err = hb.GetNext()
    // add `nextByte` to the hash algorithm
// ...while byteAvailable
```

A standard hash algorithm takes a word of data and creates a checksum.  For example:

```go
hb, err := NewHashBuffer(filespec, bufferSize, windowSize)
// do...
    // read the next available word from the file
    window, err = hb.GetWindow()
    // call the hash algorithm using window and use generated checksum
// ...while `window` is non-empty
}
```

The `HashBuffer` interface defines the available operations and `FileHashBuffer` provides a file-based implementation.

`NewHashBuffer()` creates a `FileHashBuffer` from a specified file name and the size of buffer to be used. The buffer can be any reasonable size larger than the window size.  This opens the file.  Your code should call, or defer a call, to `Close()`.

`Close()` closes the associated file and the hashbuffer.

`GetWindow()` retrieves a slice of bytes of up to the specified length, which is the window length.  If called repeatedly, it returns the next slice, one byte further in the stream, as described above.

`GetNext()` retrieves the next available byte.  It returns the byte and true to indicate success, or 0 and false if no byte is available; or an error.

`SetTesting()` allows for logging to be sent when testing HashBuffer.  The output is available if the test is run in verbose mode (`go test -v`).
