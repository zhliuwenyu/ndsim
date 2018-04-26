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

		for i, v := range []byte("abcdefgHIJKLMN") {
			if i >= 64 {
				break
			}
		fi.HashKey[i] = v
		}
	*/
	t.Log(fi)
	var hs ndsim.HashSign
	for i, v := range []byte("abcdefgHIJKLMN") {
		if i >= 64 {
			break
		}
		hs[i] = v
	}

	outBytes := ndsim.SerializeReverseIndex(hs, fi)
	t.Log(outBytes)

	hs, ri, err := ndsim.UnserializeReverseIndex(outBytes)
	if err != nil {
		t.Error(err)
	}
	t.Log(hs, ri)

}

func TestForwardIndexSerialize(t *testing.T) {
	fi := ndsim.ForwardIndex{Docid: 1000000, HashListOffset: 655389, HashListLen: 200}
	t.Log(fi)
	stream := ndsim.SerializeForwardIndex(fi)

	t.Log(stream)
	decodeFi, err := ndsim.UnserializeForwardIndex(stream)
	if err != nil {
		t.Error(err)
	}
	t.Log(decodeFi)
}

func TestAddReverseDoc(t *testing.T) {
	ndsim.InitTest()

	var a, b, c [64]byte
	a[0] = 'd'
	b[2] = 'e'
	c[4] = 'f'

	for i := 0; i < 10000000; i++ {
		ndsim.AddReverseDoc(a, ndsim.DocID(i))
	}
	t.Log(ndsim.GIndexControl)
	for i := 0; i < 10; i++ {
		ndsim.AddReverseDoc(b, ndsim.DocID(i))
	}
	t.Log(ndsim.GIndexControl)
	for i := 0; i < 1000; i++ {
		ndsim.AddReverseDoc(c, ndsim.DocID(i))
	}
	t.Log(ndsim.GIndexControl)

}

func TestGetReverseDocList(t *testing.T) {
	ndsim.InitTest()

	var a [64]byte
	a[0] = 'd'
	//b[2] = 'e'
	//c[4] = 'f'

	t.Log(ndsim.GIndexControl)
	docList := ndsim.GetReverseDocList(a)
	t.Log(docList[5], docList[len(docList)-199])
	/*
		docList = ndsim.GetReverseDocList(b)
		t.Log(docList)
		docList = ndsim.GetReverseDocList(c)
		t.Log(docList)*/
}

func TestAddForwardHashList(t *testing.T) {
	ndsim.InitTest()
	var a, b, c [64]byte
	a[0] = 'H'
	b[2] = 'e'
	c[4] = 'f'
	t.Log(ndsim.GIndexControl.LastForwardContentOffset)
	doc1, doc2 := ndsim.DocID(100000000), ndsim.DocID(100005000)
	err := ndsim.AddDocHashList(doc1, []ndsim.HashSign{a})
	t.Log(err, ndsim.GIndexControl.LastForwardContentOffset)
	err = ndsim.AddDocHashList(doc2, []ndsim.HashSign{a, c})
	t.Log(err, ndsim.GIndexControl.LastForwardContentOffset)

}

func TestGetForwardHashList(t *testing.T) {
	ndsim.InitTest()
	doc1, doc2 := ndsim.DocID(100000000), ndsim.DocID(100005000)
	hsList := ndsim.GetForwardHashList(doc1)
	t.Log(hsList)
	hsList = ndsim.GetForwardHashList(doc2)
	t.Log(hsList)
}
