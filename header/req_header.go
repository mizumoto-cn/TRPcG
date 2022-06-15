package header

import (
	"encoding/binary"
	"errors"
	"sync"

	"github.com/mizumoto-cn/TRPcG/compressor"
)

const (
	MaxHeaderSize = 36
	Uint32Size    = 4
	Uint16Size    = 2
)

var ErrUnmarshalFail = errors.New("an error occurred in Unmarshal")

type CompressType uint16

type RequestHeader struct {
	// A RWMutex is a reader/writer mutual exclusion lock.
	// The lock can be held by an arbitrary number of readers or a single writer.
	// The zero value for a RWMutex is an unlocked mutex.
	sync.RWMutex
	CompressType CompressType // uint16, used to indicate the compress type. TRPcG supports Raw/Gzip/Snappy/Zlib
	Method       string       //
	ID           uint64       // ID of the request
	RequestLen   uint32       // Lenghth of the request body
	Checksum     uint32       // CRC32 hashed value for checksum
}

// Maarshal is somewhat a encoder
func (r *RequestHeader) Marshal() []byte {
	// lock and Unlock Readlock at the end
	r.RLock()
	defer r.RUnlock()
	itor := 0
	// 2 + 10 * 3 + 4 + string length
	header := make([]byte, MaxHeaderSize+len(r.Method))

	// | CompressType |      Method    |    ID    | RequestLen | Checksum |
	// |    uint16    | uvarint+string |  uvarint |   uvarint  |  uint32  |
	// write uint16 compressType
	// LittleEndian PutType functions encode Type into buf and returns the number of bytes written
	// Here it writes uint16 type info into header
	binary.LittleEndian.PutUint16(header[itor:], uint16(r.CompressType))
	itor += Uint16Size

	// Then id:
	itor += writeString(header[itor:], r.Method)
	itor += binary.PutUvarint(header[itor:], r.ID)
	itor += binary.PutUvarint(header[itor:], uint64(r.RequestLen))

	// At last CheckSum
	binary.LittleEndian.PutUint32(header[itor:], r.Checksum)
	itor += Uint32Size

	return header[:itor]
}

func (r *RequestHeader) UnMarshal(data []byte) (err error) {
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
	// reads out bytes of a uint16 from []byte
	r.CompressType = CompressType(binary.LittleEndian.Uint16(data[itor:]))
	itor += Uint16Size

	r.Method, size = readString(data[itor:])
	itor += size

	r.ID, size = binary.Uvarint(data[itor:])
	itor += size

	length, size := binary.Uvarint(data[itor:])
	r.RequestLen = uint32(length)
	itor += size

	r.Checksum = binary.LittleEndian.Uint32(data[itor:])

	return
}

func writeString(data []byte, str string) int {
	// return 0
	itor := 0
	itor += binary.PutUvarint(data, uint64(len(str)))
	copy(data[itor:], str)
	itor += len(str)
	return itor
}

func readString(data []byte) (string, int) {
	// return "", 0
	itor := 0
	length, size := binary.Uvarint(data)
	itor += size
	str := string(data[itor : itor+int(length)])
	itor += len(str)

	return str, itor
}

func (r *ResponseHeader) GetCompressType() compressor.CompressType {
	r.RLock()
	defer r.RUnlock()
	return compressor.CompressType(r.CompressType)

}
