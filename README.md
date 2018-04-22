# hashbuffer

A hashbuffer is a file buffer designed to feed data to a rolling hash, such as can be found at [chmduquesne/rollinghash](http://github.com/chmduquesne/rollinghash)

The `HashBuffer` interface defines the available operations and `FileHashBuffer` provides a file-based implementation.

A rolling hash is an algorithm that generates checksums of a sliding window over a long array of bytes. It does so more efficiently than if each block was checksumed individually, building the next checksum from the value of the previous one.

`NewHashBuffer()` creates a `FileHashBuffer` from a specified file name and the size of buffer to be used. The buffer can be any reasonable size larger than the window size.  This opens the file.  Your code should call, or defer a call, to `Close()`.

`Close()` closes the associated file and the hashbuffer.

`Get()` retrieves a slice of bytes of up to the specified length, which is the window length.  This is normally the first call after opening to "prime" the rolling checksum, after which `GetNext()` is called to continue adding additional bytes.  It returns the byte slice, the number of bytes actually returned; or an error.

`GetNext()` retrieves the next available byte.  It returns the byte and true to indicate success, or 0 and false if no byte is available; or an error.

`SetTesting()` allows for logging to be sent when testing HashBuffer.  The output is available if the test is run in verbose mode (`go test -test.v`).
