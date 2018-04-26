package ndsim

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

//IndexControl define index global info
type IndexControl struct {
	Path      string
	DocIDFrom int64
	DocIDEnd  int64
}

/*ReverseIndex define reverseindex node struct
包含  hashkey， hashkey所在磁盘偏移，  doclist所在磁盘偏移，doclist容量，doclist 长度
*/
type ReverseIndex struct {
	HashKey       HashSign
	KeyOffset     int64
	DocListOffset int64
	DocListLen    int32
	DocListCap    int32
}

//SerializeReverseIndex trans ReverseIndex to byte slice
func SerializeReverseIndex(ri ReverseIndex) []byte {

	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, ri.HashKey)
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
			fmt.Println(err)
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

//UnserializeReverseIndex trans byte slice to ReverseIndex
func UnserializeReverseIndex(raw []byte) (*ReverseIndex, error) {
	ri := ReverseIndex{}
	rawLen := (int)(unsafe.Sizeof(ri))
	if len(raw) < rawLen {
		return nil, errors.New("raw has no enough length")
	}

	for i := 0; i < HashLength; i++ {
		ri.HashKey[i] = raw[i]
	}
	buf := bytes.NewReader(raw[HashLength:])
	binary.Read(buf, binary.LittleEndian, &ri.KeyOffset)
	binary.Read(buf, binary.LittleEndian, &ri.DocListOffset)
	binary.Read(buf, binary.LittleEndian, &ri.DocListLen)
	binary.Read(buf, binary.LittleEndian, &ri.DocListCap)
	return &ri, nil
}
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
