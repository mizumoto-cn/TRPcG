package codec

import (
	"bufio"
	"hash/crc32"
	"io"
	"log"
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

func (client *clientCodec) Close() error {
	return client.c.Close()
}

// WriteRequest Write the rpc request header and body to the io stream
func (client *clientCodec) WriteRequest(r *rpc.Request, param any) error {
	client.mutex.Lock()
	client.pending[r.Seq] /*sequence number chosen by client*/ = r.ServiceMethod // format service.method
	client.mutex.Unlock()
	// to check whether a map contains a key
	// that is whether there is a compressor
	if _, ok := compressor.Compressors[client.compressor]; !ok {
		return ErrCompressorNotFound
	}
	reqBody, err := client.serializer.Marshal(param)
	if err != nil {
		return err
	}
	//compress
	compressorMethod := compressor.Compressors[client.compressor]
	c_reqBody, err := compressorMethod.Zip(reqBody)
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
	h.CompressType = compressor.CompressType(client.compressor)
	h.Checksum = crc32.ChecksumIEEE(c_reqBody)
	// send Req Header
	err = sendFrame(client.w, h.Marshal())
	if err != nil {
		return err
	}
	// send req body
	if err = write(client.w, c_reqBody); err != nil {
		return err
	}

	client.w.(*bufio.Writer).Flush()
	return nil
}

// ClientCodec::ReadResponseHeader() implement
func (client *clientCodec) ReadResponseHeader(r *rpc.Response) error {
	//reset req header
	client.response.ResetHeader()
	//receive header
	data, err := receiveFrame(client.r)
	if err != nil {
		return nil
	}
	err = client.response.Unmarshal(data)
	if err != nil {
		return err
	}
	client.mutex.Lock()

	r.Seq = client.response.ID
	r.Error = client.response.Error
	// infer service method from seqID
	r.ServiceMethod = client.pending[r.Seq]
	// delete seqID
	delete(client.pending, r.Seq)
	client.mutex.Unlock()
	return nil
}

// ClientCodec::ReadResponseBody implementation
func (client *clientCodec) ReadResponseBody(param any) error {
	if param == nil {
		if client.response.ResponseLen != 0 {
			err := read(client.r, make([]byte, client.response.ResponseLen))
			if err != nil {
				return err
			}
		}
		return nil
	}
	// read  ResLen size bytes
	resBody := make([]byte, client.response.ResponseLen)
	err := read(client.r, resBody)
	if err != nil {
		return err
	}
	// Check
	if client.response.CheckSum != 0 {
		if crc32.ChecksumIEEE(resBody) != client.response.CheckSum {
			return ErrUnexpectedChecksum
		}
	}
	// check compressor
	boo := client.response.GetCompressType()
	_, ok := compressor.Compressors[boo]
	if !ok {
		return ErrCompressorNotFound
	}
	if boo != client.compressor {
		log.Fatalf("compressor type mismatch: %d != %d", client.response.GetCompressType(), client.compressor)
		return ErrCompressorTypeMismatch
	}
	// unzip
	res, err := compressor.Compressors[client.response.GetCompressType()].Unzip(resBody)
	if err != nil {
		return err
	}
	// Unmarshal
	return client.serializer.Unmarshal(res, param)
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
