package codec

import (
	"bufio"
	"io"
	"net/rpc"
	"sync"

	"github.com/mizumoto-cn/TRPcG/compressor"
	"github.com/mizumoto-cn/TRPcG/header"
	"github.com/mizumoto-cn/TRPcG/serializer"
)

// To implement /net/rpc required ClientCodec interface
// src/net/rpc/server.go

// type ClientCodec interface {
// 	WriteRequest(*Request, any) error
// 	ReadResponseHeader(*Response) error
// 	ReadResponseBody(any) error

// 	Close() error
// }


// Request 
// type Request struct {
// 	ServiceMethod string 
// 	Seq           uint64 
// 	next          *Request 
// }

// // Response 
// type Response struct {
// 	ServiceMethod string 
// 	Seq           uint64 
// 	Error         string 
// 	next          *Response 
// }

type clientCodec struct {
	r io.Reader
	w io.Writer
	c io.Closer

	compressor compressor.CompressType
	serializer serializer.Serializer
	response   header.ResponseHeader
	mutex      sync.Mutex // protect pending map
	pending    map[uint64]string
}

func (this *clientCodec) Close() error {
	return this.c.Close()
}

func (this *clientCodec) WriteRequest(r *rpc.Request, param any) error{
	this.mutex.Lock()
	this.pending[r.Seq] /*sequence number chosen by client*/ = r.ServiceMethod // format service.method
	this.mutex.Unlock()

	// to check whether a map contains a key
	if _,ok := compressor.Compressors[this.compressor]; !ok {
		return ErrCompressorNotFound
	}
	
}

// use bufio 
func NewClientCodec(conn io.ReadWriteCloser, compressType compressor.CompressType,
	serializer serializer.Serializer) rpc.ClientCodec {

	return nil
	}
}
