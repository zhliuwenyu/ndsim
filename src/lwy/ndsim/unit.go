package ndsim

//HashLength define simhash ret length in byte array
const HashLength = 64

//DocIDLength docid length is  sizeof(DocID)/8
const DocIDLength = 8

//HashSign define type name of sign
type HashSign [HashLength]byte

//DocID name docId type
type DocID uint64

//DiskOffset name diskoffset
type DiskOffset uint64
