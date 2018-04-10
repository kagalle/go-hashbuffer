package rollhash

type HashBuffer interface {
	Open(filespec string) (err error)
	Get(numberOfBytes int) (output []byte, err error)
	GetNext() (output byte, err error)
}

type FileHashBuffer struct {
	buffersize int
	buffer []byte
}

/*
// the buffer must be at least as big as the window size
if fhb.bufferSize < windowSize {
	bufferSize = windowSize
}
*/

// bufferSize := 1024
func NewFileHashBuffer(buffersize int) (fhb *FileHashBuffer) {
	fhb := new(FileHashBuffer)
	fhb.buffersize = buffersize
}

func (fhb *FileHashBuffer)	Open(filespec string) (err error) {

		fhb.buffer = make([]byte, bufferSize)

		s, err := ioutil.ReadFile("/home/ken/2002-06-19.pdf")
		f, err := os.Open("/tmp/dat") // f : *os.File which implements io.Reader
		check(err)
		defer f.Close()

		for {
			bytesread, err := file.Read(buffer)

			if err != nil {
				if err != io.EOF {
					fmt.Println(err)
				}
				break
			}
			fmt.Printf("%d bytes read\n", bytesread)
	}
	Get(numberOfBytes int) (output []byte, err error) {

	}
	GetNext() (output byte, err error) {

	}

}
