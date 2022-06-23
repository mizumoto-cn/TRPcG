package codec

import "errors"

var (
	ErrInvalidSeqID           = errors.New("invalid sequence number in response")
	ErrUnexpectedChecksum     = errors.New("unexpected checksum")
	ErrCompressorNotFound     = errors.New("not found compressor")
	ErrCompressorTypeMismatch = errors.New("compressor type mismatch")
)
