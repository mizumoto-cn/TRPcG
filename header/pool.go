package header

import "sync"

var (
	RequestPool  sync.Pool
	ResponsePool sync.Pool
)

func init() {
	// boo
	RequestPool = sync.Pool{New: func() any { return &RequestHeader{} }}

	ResponsePool = sync.Pool{New: func() any { return &ResponseHeader{} }}
}

func (r *RequestHeader) ResetHeader() error {
	r.ID = 0
	r.Checksum = 0
	r.Method = ""
	r.CompressType = 0
	r.RequestLen = 0
	return nil
}

func (r *ResponseHeader) ResetHeader() error {
	// return nil
	r.Error = ""
	r.ID = 0
	r.CompressType = 0
	r.ResponseLen = 0
	r.CheckSum = 0
	return nil
}
