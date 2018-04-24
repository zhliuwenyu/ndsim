package ndsim_test

import (
	"lwy/ndsim"
	"testing"
)

func TestReverseIndexSerialize(t *testing.T) {
	var fi = ndsim.ReverseIndex{
		KeyOffset:     65539,
		DocListOffset: 700000,
		DocListLen:    20,
		DocListCap:    50,
	}
	/*
		readIO := bytes.NewReader([]byte("abcdefgHIJKLMN"))
		writeIO := bytes.NewBuffer(fi.HashKey[0:64])
		io.CopyN(writeIO, readIO, 64)
	*/
	for i, v := range []byte("abcdefgHIJKLMN") {
		if i >= 64 {
			break
		}
		fi.HashKey[i] = v
	}

	t.Log(fi)

	outBytes := ndsim.SerializeReverseIndex(fi)
	t.Log(outBytes)

	ri, err := ndsim.UnserializeReverseIndex(outBytes)
	if err != nil {
		t.Error(err)
	}
	t.Log(ri)

}

func TestForwardIndexSerialize(t *testing.T) {
	fi := ndsim.ForwardIndex{DocID: 1000000, HashListOffset: 655389, HashListLen: 200}
	t.Log(fi)
	stream := ndsim.SerializeForwardIndex(fi)

	t.Log(stream)
	decodeFi, err := ndsim.UnserializeForwardIndex(stream)
	if err != nil {
		t.Error(err)
	}
	t.Log(decodeFi)
}
