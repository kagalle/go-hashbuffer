package rollhash

import "testing"

// empty file
// 1 byte file
// 1 byte less than full buffer
// 1 full buffer
// 1 full buffer + 1
// long

// {"empty", 0},
// {"onebyte", 1},
// {"onelessthanone", 1023},
// {"onebuffer", 1024},
// {"onebufferplusone", 1025},
// {"long", 35538},

const bufferSize = 1024

func TestBufferEmptyFileWithGet(t *testing.T) {
	t.Log("start TestBufferEmptyFileWithGet")
	hb := NewHashBuffer("./testdata/empty", bufferSize)
	defer hb.Close()
	{
		t.Log("starting first")
		buf, count := hb.Get(0)
		if count != 0 {
			t.Errorf("Error TestBufferVariousLengths: got count=%d, want 0", count)
		}
		len := len(buf)
		if len != 0 {
			t.Errorf("Error TestBufferVariousLengths: got length=%d, want 0", len)
		}
	}
	{
		t.Log("starting second")
		buf, count := hb.Get(0)
		if count != 0 {
			t.Errorf("Error TestBufferVariousLengths(2): got count=%d, want 0", count)
		}
		len := len(buf)
		if len != 0 {
			t.Errorf("Error TestBufferVariousLengths(2): got length=%d, want 0", len)
		}

	}
}
func TestBufferEmptyFileWithGetNext(t *testing.T) {
	t.Log("start TestBufferEmptyFileWithGetNext")
	hb := NewHashBuffer("./testdata/empty", bufferSize)
	defer hb.Close()
	t.Log("starting")
	outByte, ok := hb.GetNext()
	if ok {
		t.Error("Error TestBufferVariousLengths: got ok=true, want ok=false")
	}
	if outByte != 0 {
		t.Errorf("Error TestBufferVariousLengths: got outByte=%d, want 0", outByte)
	}
}
