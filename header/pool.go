package header

import "sync"

var (
	RequestPool  sync.Pool
	ResponsePool sync.Pool
)

func init() {
	// boo
}

func (r *RequestHeader) ResetHeader() error {
	return nil
}

func (r *ResponseHeader) ResetHeader() error {
	return nil
}
