package codec

import (
	"bufio"
	"hash/crc32"
	"io"
	"net/rpc"
	"sync"

	"github.com/mizumoto-cn/TRPcG/compressor"
	"github.com/mizumoto-cn/TRPcG/header"
	"github.com/mizumoto-cn/TRPcG/serializer"
)

// type ServerCodec interface{
// 	ReadRequestHeader(*Request)error
// 	ReadRequestBody(any)error
// 	WriteResponse(*Response,any)error
// 	Close()error
// }

// reqContext is a context for a request.
type reqContext struct {
	id             uint64
	compressorType compressor.CompressType
}

type serverCodec struct {
	r io.Reader
	w io.Writer
	c io.Closer

	request    header.RequestHeader
	serializer serializer.Serializer
	mutex      sync.Mutex
	seq        uint64
	pending    map[uint64]*reqContext
}

// ServerCodec::ReadRequestHeader()
func (server *serverCodec) ReadRequestHeader(r *rpc.Request) error {
	server.request.ResetHeader()
	data, err := receiveFrame(server.r)
	if err != nil {
		return err
	}
	err = server.request.UnMarshal(data)
	if err != nil {
		return err
	}
	server.mutex.Lock()
	server.seq++ // add one to seqID
	server.pending[server.seq] = &reqContext{server.request.ID, server.request.GetCompressType()}
	r.ServiceMethod = server.request.Method
	r.Seq = server.seq
	server.mutex.Unlock()
	return nil
}

// ServerCodec::ReadRequestBody()
func (server *serverCodec) ReadRequestBody(param any) error {
	if param == nil {
		// throw unused bytes
		if server.request.RequestLen != 0 {
			err := read(server.r, make([]byte, server.request.RequestLen))
			if err != nil {
				return err
			}
		}
		return nil
	}

	reqBody := make([]byte, server.request.RequestLen)
	// read bytes of sizeof request body
	err := read(server.r, reqBody)
	if err != nil {
		return err
	}
	// check
	if server.request.Checksum != 0 {
		if crc32.ChecksumIEEE(reqBody) != server.request.Checksum {
			return ErrUnexpectedChecksum
		}
	}
	// check compressor
	_, ok := compressor.Compressors[server.request.GetCompressType()]
	if !ok {
		return ErrCompressorNotFound
	}
	// Unzip
	req, err := compressor.Compressors[server.request.GetCompressType()].Unzip(reqBody)
	if err != nil {
		return err
	}
	// Unmarshal
	return server.serializer.Unmarshal(req, param)
}

// ServerCodec::WriteResponse()
func (server *serverCodec) WriteResponse(r *rpc.Response, param any) error {
	server.mutex.Lock()
	reqContext, ok := server.pending[r.Seq]
	if !ok {
		server.mutex.Unlock()
		return ErrInvalidSeqID
	}
	delete(server.pending, r.Seq)
	server.mutex.Unlock()

	// if it's not a adequate rpc-call, set param to nil
	if r.Error != "" {
		param = nil
	}
	// check compressor
	if _, ok := compressor.Compressors[reqContext.compressorType]; !ok {
		return ErrCompressorNotFound
	}

	var (
		resBody []byte
		err     error
	)
	if param != nil {
		resBody, err = server.serializer.Marshal(param)
		if err != nil {
			return err
		}
	}

	// Zip resBody
	compressedResBody, err := compressor.
		Compressors[reqContext.compressorType].Zip(resBody)
	if err != nil {
		return err
	}

	// Get a new res header
	h := header.ResponsePool.Get().(*header.ResponseHeader)
	defer func() {
		h.ResetHeader()
		header.ResponsePool.Put(h)
	}()
	h.ID = reqContext.id
	h.Error = r.Error
	h.ResponseLen = uint32(len(compressedResBody))
	h.CheckSum = crc32.ChecksumIEEE(compressedResBody)
	h.CompressType = reqContext.compressorType

	err = sendFrame(server.w, h.Marshal())
	if err != nil {
		return err
	}
	err = write(server.w, compressedResBody)
	if err != nil {
		return err
	}
	server.w.(*bufio.Writer).Flush()
	return nil
}

func (server *serverCodec) Close() error {
	return server.c.Close()
}

func NewServerCodec(conn io.ReadWriteCloser, serializer serializer.Serializer) rpc.ServerCodec {
	return &serverCodec{
		r:          bufio.NewReader(conn),
		w:          bufio.NewWriter(conn),
		c:          conn,
		serializer: serializer,
		pending:    make(map[uint64]*reqContext),
	}
}
