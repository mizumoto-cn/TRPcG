package header

import (
	"encoding/binary"
	"sync"
)

// const (
// 	MaxHeaderSize = 2 + 10 + 10 + 10 + 4
// 	(10 refer to binary.MaxVarintLen64)
// 	MaxHeaderSize = 36
// 	Uint32Size    = 4
// 	Uint16Size    = 2
// )

// var ErrUnmarshalFail = errors.New("an error occurred in Unmarshal")

// type CompressType uint16

// | CompressType |    ID   |      Error     | ResponseLen | CheckSum |
// |    uint16    | uvarint | uvarint+string |    uvarint  |  uint32  |
type ResponseHeader struct {
	sync.RWMutex
	CompressType CompressType // uint16
	ID           uint64       // response id
	Error        string       // error info
	ResponseLen  uint32       // Length of the response body
	CheckSum     uint32       // for check
}

// Marshal() encode response header into byte slice
func (r *ResponseHeader) Marshal() []byte {
	r.RLock()
	defer r.RUnlock()
	itor := 0
	// 36 + errstr length
	header := make([]byte, MaxHeaderSize+len(r.Error))
	// putin cType
	binary.LittleEndian.PutUint16(header[itor:], uint16(r.CompressType))
	itor += Uint16Size
	// putin ID errstr resBodyLength
	itor += binary.PutUvarint(header[itor:], r.ID)
	itor += writeString(header[itor:], r.Error)
	itor += binary.PutUvarint(header[itor:], uint64(r.ResponseLen))
	// putin checksum
	binary.LittleEndian.PutUint32(header[itor:], r.CheckSum)
	itor += Uint32Size
	return header[:itor]
}

// Unmarshal will decode request header into a byte slice
func (r *ResponseHeader) Unmarshal(data []byte) (err error) {
	r.Lock()
	defer r.Unlock()
	if len(data) == 0 {
		return ErrUnmarshalFail
	}
	defer func() {
		if r := recover(); r != nil {
			err = ErrUnmarshalFail
		}
	}()

	itor, size := 0, 0
	r.CompressType = CompressType(binary.LittleEndian.Uint16(data[itor:]))
	itor += Uint16Size

	r.ID, size = binary.Uvarint(data[itor:])
	itor += size

	r.Error, size = readString(data[itor:])
	itor += size

	length, size := binary.Uvarint(data[itor:])
	r.ResponseLen = uint32(length)
	itor += size

	r.CheckSum = binary.LittleEndian.Uint32(data[itor:])
	return
}

// func writeString(data []byte, str string) int {
//
// 	itor := 0
// 	itor += binary.PutUvarint(data, uint64(len(str)))
// 	copy(data[itor:], str)
// 	itor += len(str)
// 	return itor
// }

// func readString(data []byte) (string, int) {
//
// 	itor := 0
// 	length, size := binary.Uvarint(data)
// 	itor += size
// 	str := string(data[itor : itor+int(length)])
// 	itor += len(str)
//
// 	return str, itor
// }
