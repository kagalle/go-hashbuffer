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
	hb.SetTesting(t)
	defer hb.Close()
	testGet(t, hb, "TestBufferEmptyFileWithGet first pass", 0, 0)
	testGet(t, hb, "TestBufferEmptyFileWithGet second pass", 0, 0)
}

func TestBufferEmptyFileWithGetNext(t *testing.T) {
	t.Log("start TestBufferEmptyFileWithGetNext")
	hb := NewHashBuffer("./testdata/empty", bufferSize)
	hb.SetTesting(t)
	defer hb.Close()
	testGetNext(t, hb, "TestBufferEmptyFileWithGetNext first pass", false)
	testGetNext(t, hb, "TestBufferEmptyFileWithGetNext second pass", false)
}

func testGet(t *testing.T, hb HashBuffer, title string, amountToGet int, expectedLength int) {
	t.Logf("starting %s", title)
	buf, count := hb.Get(amountToGet)
	if count != expectedLength {
		t.Errorf("Error %s :got count=%d, want %d", title, count, expectedLength)
	}
	if amountToGet == 0 {
		if buf != nil {
			t.Errorf("Error %s: got buf[] with length=%d, want buf == nil", title, len(buf))
		}
	} else {
		len := len(buf)
		if len != expectedLength {
			t.Errorf("Error %s: got length=%d, want %d", title, len, expectedLength)
		}
	}
}

func testGetNext(t *testing.T, hb HashBuffer, title string, byteExpectedInReturn bool) {
	t.Logf("starting %s", title)
	outByte, ok := hb.GetNext()
	if byteExpectedInReturn {
		if !ok {
			t.Errorf("Error %s: got ok=false, want ok=true", title)
		}
	} else {
		if ok {
			t.Errorf("Error %s: got ok=true, want ok=false", title)
		}
		if outByte != 0 {
			t.Errorf("Error %s: got outByte=%d, want 0", title, outByte)
		}
	}
}
