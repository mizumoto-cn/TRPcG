package codec

import (
	"io"

	"github.com/mizumoto-cn/TRPcG/compressor"
	"github.com/mizumoto-cn/TRPcG/serializer"
)

// To implement /net/rpc required ClientCodec interface

type clientCodec struct {
	r io.Reader
	w io.Writer
	c io.Closer

	compressor compressor.CompressType
	serializer serializer.Serializer
}
