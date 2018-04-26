package ndsim

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

/*ForwardIndex define forwardindex node struct
包含 docid、hashlist磁盘偏移、hashlist长度
*/
type ForwardIndex struct {
	Docid          DocID
	HashListOffset DiskOffset
	HashListLen    int32
}

//FILength define ForwardIndex length
const FILength = 8 + 8 + 4

//AddDocHashList add hashsign list to a docid, add hs list to forward content and change forwardindex
func AddDocHashList(docid DocID, hsList []HashSign) error {
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
	fi := ForwardIndex{Docid: docid, HashListOffset: GIndexControl.LastForwardContentOffset, HashListLen: int32(len(hsList))}
	GIndexControl.LastForwardContentOffset += DiskOffset(len(hsList) * HashLength)
	return freshForwardIndexToDisk(fi)
}

//GetForwardHashList get hashsign list belong doc
func GetForwardHashList(docid DocID) []HashSign {
	indexOffset := FILength * int64(docid-gConfig.DocIDFrom)
	if _, err := GIndexControl.ForwardIndexFile.Seek(indexOffset, 0); err != nil {
		gLog.Warning(err)
		return nil
	}
	buf := make([]byte, FILength)
	gLog.Debug(indexOffset)
	if n, err := GIndexControl.ForwardIndexFile.Read(buf); nil != err {
		gLog.Warning(err)
	} else {
		if n != FILength {
			gLog.Warning(fmt.Sprintf("get doc[%d] forward index expect[%d]bytes but get [%d]bytes ", docid, FILength, n))
		} else {
			if fi, err := UnserializeForwardIndex(buf); err == nil {
				gLog.Debug(*fi)
				if fi.Docid != docid {
					gLog.Warning(fmt.Sprintf("docid[%d] index content[%v] with docid dont match", docid, fi))
					return nil
				}
				hsByte := make([]byte, fi.HashListLen*HashLength)
				if _, err := GIndexControl.ForwardContentFile.Seek(int64(fi.HashListOffset), 0); err != nil {
					gLog.Warning(err)
				} else {
					if n, err := GIndexControl.ForwardContentFile.Read(hsByte); err != nil || int32(n) != fi.HashListLen*HashLength {
						gLog.Warning(err, n)
					} else {
						hsList := make([]HashSign, fi.HashListLen)
						for i := 0; i < int(fi.HashListLen); i++ {
							copy(hsList[i][:], hsByte[i*HashLength:(i+1)*HashLength])
						}
						return hsList
					}
				}
			} else {
				gLog.Warning(err)
			}
		}
	}
	return nil
}

func freshForwardIndexToDisk(fi ForwardIndex) error {
	buf := SerializeForwardIndex(fi)
	indexOffset := FILength * int64(fi.Docid-gConfig.DocIDFrom)
	if _, err := GIndexControl.ForwardIndexFile.Seek(indexOffset, 0); err != nil {
		gLog.Warning(err)
		return err
	}
	n, err := GIndexControl.ForwardIndexFile.Write(buf)
	if err != nil {
		gLog.Warning(err)
		return err
	}
	if n != FILength {
		errMsg := fmt.Sprintf("fresh docid [%d] ForwardIndexcontent [%v]to disk with unexpect length[%d] ", fi.Docid, fi, n)
		gLog.Warning(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

//SerializeForwardIndex trans ForwardIndex to byte slice
func SerializeForwardIndex(fi ForwardIndex) []byte {

	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, fi.Docid)
	binary.Write(buf, binary.LittleEndian, fi.HashListOffset)
	binary.Write(buf, binary.LittleEndian, fi.HashListLen)

	return buf.Bytes()
}

//UnserializeForwardIndex trans byte slice to ReverseIndex
func UnserializeForwardIndex(raw []byte) (*ForwardIndex, error) {
	fi := ForwardIndex{}
	rawLen := FILength
	if len(raw) < rawLen {
		return nil, errors.New("raw has no enough length")
	}

	buf := bytes.NewReader(raw)
	binary.Read(buf, binary.LittleEndian, &fi.Docid)
	binary.Read(buf, binary.LittleEndian, &fi.HashListOffset)
	binary.Read(buf, binary.LittleEndian, &fi.HashListLen)
	return &fi, nil
}
