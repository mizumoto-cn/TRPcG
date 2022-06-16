package tinyrpcgo

import (
	"io"
	"net/rpc"

	"github.com/mizumoto-cn/TRPcG/codec"
	"github.com/mizumoto-cn/TRPcG/compressor"
	"github.com/mizumoto-cn/TRPcG/serializer"
)

type Client struct {
	*rpc.Client
}

// Functional Options Pattern
type Option func(o *options)

type options struct {
	compressType compressor.CompressType
	serializer   serializer.Serializer
}

// set compression type
func WithCompress(c compressor.CompressType) Option {
	return func(o *options) {
		o.compressType = c
	}
}

// set client serializer
func WithSerializer(serializer serializer.Serializer) Option {
	return func(o *options) {
		o.serializer = serializer
	}
}

// Create New rpc client object
func NewClient(conn io.ReadWriteCloser, args ...Option) *Client {
	options := options{
		compressType: compressor.Raw,
		serializer:   serializer.Proto,
	}
	for _,option := range args{
		option(&options)
	}
	return &Client{
		rpc.NewClientWithCodec(codec.NewClientCodec(conn, options.compressType, options.serializer))
	}
}

// synchronous call
func (c*Client)Call (serviceMethod string, args any, reply any ) error {
	return c.Client.Call(serviceMethod, args, reply)
}

// Async call  asynchronously calls the rpc function and returns a channel of *rpc.Call
func (c *Client)AsyncCall(serviceMethod string, args any, reply any) chan *rpc.Call {
	return c.Go(serviceMethod, args, reply, nil).Done	
}