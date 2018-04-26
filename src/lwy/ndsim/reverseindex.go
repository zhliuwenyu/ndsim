package ndsim

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

func loadReverseIndexToMap() error {
	buf := make([]byte, RISerLength)
	GIndexControl.ReverseIndexFile.Seek(0, 0)
	n, err := GIndexControl.ReverseIndexFile.Read(buf)

	for i := 0; n == RISerLength && err == nil; n, err = GIndexControl.ReverseIndexFile.Read(buf) {
		hs, pri, err := UnserializeReverseIndex(buf)
		if err != nil {
			return err
		}
		GIndexControl.ReverseIndexMap[hs] = *pri
		if pri.DocListOffset+DiskOffset(pri.DocListCap) > GIndexControl.MaxReverseContentOffset {
			GIndexControl.MaxReverseContentOffset = pri.DocListOffset + DiskOffset(pri.DocListCap)
		}
		if pri.KeyOffset != DiskOffset(i*RISerLength) {
			errMsg := fmt.Sprintf("load reverseindex file fail with index format err ,block[%d] offset[%d] should be[%d] ", i, pri.KeyOffset, i*RISerLength)
			gLog.Fatal(errMsg)
			return errors.New(errMsg)
		}
		i++
	}
	if err != nil && err != io.EOF {
		errMsg := "load reverseindex file fail with io error"
		gLog.Fatal(errMsg)
		return err
	} //else is for err == nil   or  err == io.EOF
	if err == nil && n != RISerLength { //err ==nil && n== blockLength will stay in for loop
		errMsg := "load reverseindex file fail with unexpect file length"
		gLog.Fatal(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

const defaultDocListCap = 16

//AddReverseDoc add new docId to hashsign hs
func AddReverseDoc(hs HashSign, docid DocID) error {
	ri, ok := GIndexControl.ReverseIndexMap[hs]
	if docid > gConfig.DocIDEnd || docid < gConfig.DocIDFrom {
		errMsg := ""
		gLog.Warning("hs", hs, "add new docid", docid, "out of range[", gConfig.DocIDFrom, "~", gConfig.DocIDEnd, "]")
		return errors.New(errMsg)
	}

	//为了线程安全  保证都   先写数据文件  再更新map 最后更新索引文件
	if ok {
		if ri.DocListCap > ri.DocListLen {
			if _, err := GIndexControl.ReverseContentFile.Seek(int64(ri.DocListOffset)+int64(ri.DocListLen*DocIDLength), 0); err != nil {
				gLog.Warning(err)
				return err
			}
			if err := binary.Write(GIndexControl.ReverseContentFile, binary.LittleEndian, uint64(docid)); err != nil {
				gLog.Warning(err)
				return err
			}
			ri.DocListLen++
		} else {
			//扩容流程
			buf := make([]byte, ri.DocListCap*DocIDLength)
			if _, err := GIndexControl.ReverseContentFile.Seek(int64(ri.DocListOffset), 0); err != nil {
				gLog.Warning(err)
				return err
			}
			if _, err := GIndexControl.ReverseContentFile.Read(buf); err != nil {
				gLog.Warning(err)
				return err
			}
			if _, err := GIndexControl.ReverseContentFile.Seek(int64(GIndexControl.MaxReverseContentOffset), 0); err != nil {
				gLog.Warning(err)
				return err
			}
			if _, err := GIndexControl.ReverseContentFile.Write(buf); err != nil {
				gLog.Warning(err)
				return err
			}
			if err := binary.Write(GIndexControl.ReverseContentFile, binary.LittleEndian, uint64(docid)); err != nil {
				gLog.Warning(err)
				return err
			}
			ri = ReverseIndex{
				KeyOffset:     ri.KeyOffset,
				DocListOffset: GIndexControl.MaxReverseContentOffset,
				DocListLen:    ri.DocListLen + 1,
				DocListCap:    ri.DocListCap * 2,
			}
			GIndexControl.MaxReverseContentOffset += DiskOffset(ri.DocListCap * 8)
			gLog.Debug("hs", hs, "increase capacity to", ri.DocListCap)
		}
	} else {
		//新加hashkey流程
		ri = ReverseIndex{
			KeyOffset:     DiskOffset(RISerLength * len(GIndexControl.ReverseIndexMap)),
			DocListOffset: GIndexControl.MaxReverseContentOffset,
			DocListLen:    1,
			DocListCap:    defaultDocListCap,
		}
		if _, err := GIndexControl.ReverseContentFile.Seek(int64(GIndexControl.MaxReverseContentOffset), 0); err != nil {
			gLog.Warning(err)
			return err
		}
		if err := binary.Write(GIndexControl.ReverseContentFile, binary.LittleEndian, uint64(docid)); err != nil {
			gLog.Warning(err)
			return err
		}
		GIndexControl.MaxReverseContentOffset += DiskOffset(ri.DocListCap * 8)
		gLog.Debug("add new hs", hs, "to hashmap")
	}
	GIndexControl.ReverseIndexMap[hs] = ri
	if err := freshReverseIndexToDisk(hs, ri); err != nil {
		return err
	}
	gLog.Debug("success add docid", docid, "to hs", hs)
	return nil
}

//GetReverseDocList get all doc list from hs
func GetReverseDocList(hs HashSign) []DocID {
	ri, ok := GIndexControl.ReverseIndexMap[hs]
	if !ok {
		gLog.Debug("hs", hs, "doclist don't exist")
		return nil
	}
	buf := make([]byte, ri.DocListLen*DocIDLength)
	fmt.Println(ri)
	if _, err := GIndexControl.ReverseContentFile.Seek(int64(ri.DocListOffset), 0); err != nil {
		gLog.Warning(err)
		return nil
	}
	if n, err := GIndexControl.ReverseContentFile.Read(buf); err != nil {
		gLog.Warning(err)
	} else {
		if int32(n) != ri.DocListLen*DocIDLength {
			gLog.Warning(fmt.Sprintln("get hs", hs, " doc list length [", n, "] but expect with [", ri.DocListLen*DocIDLength, "]"))
		} else {
			reader := bytes.NewReader(buf)
			docList := make([]DocID, ri.DocListLen)
			for i := 0; i < int(ri.DocListLen); i++ {
				if err := binary.Read(reader, binary.LittleEndian, &docList[i]); err != nil {
					gLog.Warning("decode hs doclist index[", i, "] with error [", err, "]")
					return nil
				}
			}
			return docList
		}

	}
	return nil
}

func freshReverseIndexToDisk(hs HashSign, ri ReverseIndex) error {
	if _, err := GIndexControl.ReverseIndexFile.Seek(int64(ri.KeyOffset), 0); err != nil {
		gLog.Warning(err)
		return err
	}
	buf := SerializeReverseIndex(hs, ri)
	n, err := GIndexControl.ReverseIndexFile.Write(buf)
	if err != nil {
		gLog.Warning(err)
		return err
	}
	if n != RISerLength {
		errMsg := fmt.Sprintf("fresh ReverseIndex[%v] content [%v]to disk with unexpect length[%d] ", hs, ri, n)
		gLog.Fatal(errMsg)
		return errors.New(errMsg)
	}
	gLog.Debug("success freash hs ", hs, "content", ri, "to disk")
	return nil
}

/*ReverseIndex define reverseindex node struct
包含  hashkey， hashkey所在磁盘偏移，  doclist所在磁盘偏移，doclist容量，doclist 长度
*/
type ReverseIndex struct {
	//HashKey       HashSign
	KeyOffset     DiskOffset //索引记录所在偏移
	DocListOffset DiskOffset
	DocListLen    int32
	DocListCap    int32
}

//RISerLength is short for ReverseIndex serialization length
const RISerLength = HashLength + 8 + 8 + 4 + 4

//SerializeReverseIndex trans ReverseIndex to byte slice
func SerializeReverseIndex(hs HashSign, ri ReverseIndex) []byte {

	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, hs)
	binary.Write(buf, binary.LittleEndian, ri.KeyOffset)
	binary.Write(buf, binary.LittleEndian, ri.DocListOffset)
	binary.Write(buf, binary.LittleEndian, ri.DocListLen)
	binary.Write(buf, binary.LittleEndian, ri.DocListCap)

	return buf.Bytes()
}

//本来想利用反射，但是结果是还不如一行一行写
func serializeOneLevelStruct(v reflect.Value) ([]byte, error) {
	buf := new(bytes.Buffer)
	for i := 0; i < v.NumField(); i++ {
		var err error

		switch v.Field(i).Type().Kind() {
		case reflect.Slice:
			err = binary.Write(buf, binary.LittleEndian, v.Field(i).Bytes())
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
			err = binary.Write(buf, binary.LittleEndian, v.Field(i).Interface())
		}
		if err != nil {
			gLog.Warning(err)
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

//UnserializeReverseIndex trans byte slice to ReverseIndex
func UnserializeReverseIndex(raw []byte) (HashSign, *ReverseIndex, error) {
	ri := ReverseIndex{}
	var hs HashSign
	rawLen := (int)(unsafe.Sizeof(ri))
	if len(raw) != rawLen+HashLength {
		return hs, nil, errors.New("UnserializeReverseIndex input raw has diff length")
	}

	for i := 0; i < HashLength; i++ {
		hs[i] = raw[i]
	}
	buf := bytes.NewReader(raw[HashLength:])

	binary.Read(buf, binary.LittleEndian, &ri.KeyOffset) //可以不存储，此处用作格式校验
	binary.Read(buf, binary.LittleEndian, &ri.DocListOffset)
	binary.Read(buf, binary.LittleEndian, &ri.DocListLen)
	binary.Read(buf, binary.LittleEndian, &ri.DocListCap)
	return hs, &ri, nil
}
