package codec

import (
	"io"
	"net/rpc"
	"sync"

	"github.com/mizumoto-cn/TRPcG/header"
	"github.com/mizumoto-cn/TRPcG/serializer"
)

// type ServerCodec interface{
// 	ReadRequestHeader(*Request)error
// 	ReadRequestBody(any)error
// 	WriteResponse(*Response,any)error
// 	Close()error
// }

type serverCodec struct {
	r io.Reader
	w io.Writer
	c io.Closer

	request    header.RequestHeader
	serializer serializer.Serializer
	mutex      sync.Mutex
	seq        uint64
	pending    map[uint64]uint64
}

// ServerCodec::ReadRequestHeader()
func (thi *serverCodec) ReadRequestHeader(r *rpc.Request) error {
	thi.request.ResetHeader()
	data, err := receiveFrame(thi.r)
	if err != nil {
		return err
	}
	err = thi.request.UnMarshal(data)
	if err != nil {
		return err
	}
	thi.mutex.Lock()
	thi.seq++ // add one to seqID
	thi.pending[thi.seq] = thi.request.ID
	r.ServiceMethod = thi.request.Method
	r.Seq = thi.seq
	thi.mutex.Unlock()
	return nil
}

// ServerCodec::ReadRequestBody()
func (thi *serverCodec) ReadRequestBody(param any) error {
	if param == nil {
		// throw unused bytes
		if thi.request.RequestLen != 0 {
			err := read(thi.r, make([]byte, thi.request.RequestLen))
			if err != nil {
				return err
			}
		}
		return nil
	}

	reqBody := make([]byte, thi.request.RequestLen)
	// read bytes of sizeof request body
	err := read(thi.r, reqBody)
	if err != nil {
		return err
	}
}
