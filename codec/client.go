package codec

import (
	"bufio"
	"hash/crc32"
	"io"
	"net/rpc"
	"sync"

	// cSpell:ignore mizumoto
	"github.com/mizumoto-cn/TRPcG/compressor"
	"github.com/mizumoto-cn/TRPcG/header"
	"github.com/mizumoto-cn/TRPcG/serializer"
)

// To implement /net/rpc required ClientCodec interface
// src/net/rpc/server.go
//
// type ClientCodec interface {
// 	WriteRequest(*Request, any) error
// 	ReadResponseHeader(*Response) error
// 	ReadResponseBody(any) error
//
// 	Close() error
// }
//
// Request
// type Request struct {
// 	ServiceMethod string
// 	Seq           uint64
// 	next          *Request
// }
//
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

func (this_c *clientCodec) Close() error {
	return this_c.c.Close()
}

func (this_c *clientCodec) WriteRequest(r *rpc.Request, param any) error {
	this_c.mutex.Lock()
	this_c.pending[r.Seq] /*sequence number chosen by client*/ = r.ServiceMethod // format service.method
	this_c.mutex.Unlock()
	// to check whether a map contains a key
	// that is whether there is a compressor
	if _, ok := compressor.Compressors[this_c.compressor]; !ok {
		return ErrCompressorNotFound
	}
	reqBody, err := this_c.serializer.Marshal(param)
	if err != nil {
		return err
	}
	//compress
	c_reqBody, err := compressor.Compressors[this_c.compressor].Zip(reqBody)
	if err != nil {
		return err
	}
	// take out a req head from pool
	h := header.RequestPool.Get().(*header.RequestHeader) // any item from sync.Pool to RH
	defer func() {
		h.ResetHeader()
		header.RequestPool.Put(h)
	}()
	h.ID = r.Seq
	h.Method = r.ServiceMethod
	h.RequestLen = uint32(len(c_reqBody))
	h.CompressType = header.CompressType(this_c.compressor)
	h.Checksum = crc32.ChecksumIEEE(c_reqBody)
	// send Req Header
	if err = sendFrame(this_c.w, h.Marshal()); err != nil {
		return err
	}
	// send req body
	if err = write(this_c.w, c_reqBody); err != nil {
		return err
	}

	this_c.w.(*bufio.Writer).Flush()
	return nil
}

// ClientCodec::ReadResponseHeader() implement
func (this_c *clientCodec) ReadResponseHeader(r *rpc.Response) error {
	//reset req header
	this_c.response.ResetHeader()
	//receive header
	data, err := receiveFrame(this_c.r)
	if err != nil {
		return nil
	}
	err = this_c.response.Unmarshal(data)
	if err != nil {
		return err
	}
	this_c.mutex.Lock()

	r.Seq = this_c.response.ID
	r.Error = this_c.response.Error
	// infer service method from seqID
	r.ServiceMethod = this_c.pending[r.Seq]
	// delete seqID
	delete(this_c.pending, r.Seq)
	this_c.mutex.Unlock()
	return nil
}

// ClientCodec::ReadResponseBody implementation
func (this_c *clientCodec) ReadResponseBody(param any) error {
	if param == nil {
		if this_c.response.ResponseLen != 0 {
			err := read(this_c.r, make([]byte, this_c.response.ResponseLen))
			if err != nil {
				return err
			}
		}
		return nil
	}
	// read  ResLen size bytes
	resBody := make([]byte, this_c.response.ResponseLen)
	err := read(this_c.r, resBody)
	if err != nil {
		return err
	}
	// Check
	if this_c.response.CheckSum != 0 {
		if crc32.ChecksumIEEE(resBody) != this_c.response.CheckSum {
			return ErrUnexpectedChecksum
		}
	}
	// check compressor
	if _, ok := compressor.Compressors[this_c.response.GetCompressType()]; !ok {
		return ErrCompressorNotFound
	}
	// unzip
	res, err := compressor.Compressors[this_c.response.GetCompressType()].Unzip(resBody)
	if err != nil {
		return err
	}
	// Unmarshal
	return this_c.serializer.Unmarshal(res, param)
}

// use bufio
func NewClientCodec(conn io.ReadWriteCloser, compressType compressor.CompressType,
	serializer serializer.Serializer) rpc.ClientCodec {

	return &clientCodec{
		r:          bufio.NewReader(conn),
		w:          bufio.NewWriter(conn),
		c:          conn,
		compressor: compressType,
		serializer: serializer,
		pending:    make(map[uint64]string),
	}

}
