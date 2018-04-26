package ndsim

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"unsafe"
)

/*ForwardIndex define forwardindex node struct
包含 docid、hashlist磁盘偏移、hashlist长度
*/
type ForwardIndex struct {
	DocID          int64
	HashListOffset int32
	HashListLen    int32
}

//FILength define ForwardIndex length
const FILength = 8 + 4 + 4

//AddForwardIndex add hashsign list to a docid, add hs list to forward content and change forwardindex
func AddForwardIndex(docid DocID, hsList []HashSign) error {
	if docid > gConfig.DocIDEnd || docid < gConfig.DocIDFrom {
		errMsg := ""
		gLog.Warning("AddForwardIndex to docid", docid, "out of range[", gConfig.DocIDFrom, "~", gConfig.DocIDEnd, "]")
		return errors.New(errMsg)
	}
	buf := make([]byte, len(hsList)*HashLength)
	for i := 0; i < len(hsList); i++ {
		copy(buf[i*HashLength:(i+1)*HashLength], hsList[i][:])
	}
	if _, err := GIndexControl.ForwardContentFile.Seek(0, os.SEEK_END); err != nil {
		gLog.Warning(err)
		return err
	}
	if n, err := GIndexControl.ForwardContentFile.Write(buf); err != nil {
		gLog.Warning(err)
		return err
	} else if n != len(hsList)*HashLength {
		errMsg := fmt.Sprintf("AddForwardIndex hslist[%v] to docid[%d] only write [%d]bytes where expect [%d]bytes", hsList, docid, n, len(hsList)*HashLength)
		gLog.Warning(errMsg)
		return errors.New(errMsg)
	}
	GIndexControl.LastForwardContentOffset += DiskOffset(len(hsList) * HashLength)

	return nil
}

//SerializeForwardIndex trans ForwardIndex to byte slice
func SerializeForwardIndex(fi ForwardIndex) []byte {

	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, fi.DocID)
	binary.Write(buf, binary.LittleEndian, fi.HashListOffset)
	binary.Write(buf, binary.LittleEndian, fi.HashListLen)

	return buf.Bytes()
}

//UnserializeForwardIndex trans byte slice to ReverseIndex
func UnserializeForwardIndex(raw []byte) (*ForwardIndex, error) {
	fi := ForwardIndex{}
	rawLen := (int)(unsafe.Sizeof(fi))
	if len(raw) < rawLen {
		return nil, errors.New("raw has no enough length")
	}

	buf := bytes.NewReader(raw)
	binary.Read(buf, binary.LittleEndian, &fi.DocID)
	binary.Read(buf, binary.LittleEndian, &fi.HashListOffset)
	binary.Read(buf, binary.LittleEndian, &fi.HashListLen)
	return &fi, nil
}
