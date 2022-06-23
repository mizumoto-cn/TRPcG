package TRPcG

import (
	"log"
	"net"
	"net/rpc"

	"github.com/mizumoto-cn/TRPcG/codec"
	"github.com/mizumoto-cn/TRPcG/serializer"
)

// wrap /net/rpc :: Server

// Server is a RPC server based on /net/rpc.Server
type Server struct {
	*rpc.Server
	serializer.Serializer
}

// Serve accepts incoming connections on the listener l, creating a new
// ServerCodec to handle each connection. The Server will close the listener
// when it receives a signal on the Done channel.
func (server *Server) Serve(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print("trpcg.Serve: accept:", err.Error())
			return
		}
		go server.Server.ServeCodec(codec.NewServerCodec(conn, server.Serializer))
	}
}

// NewServer returns a new Server.
func NewServer(opts ...Option) *Server {
	options := options{
		serializer: serializer.Proto,
	}
	for _, opt := range opts {
		opt(&options)
	}
	return &Server{
		&rpc.Server{},
		options.serializer,
	}
}

// Register registers a rpc service with a given receiver.
func (server *Server) Register(rcvr interface{}) error {
	return server.Server.Register(rcvr)
}

// RegisterName registers a rpc service with a given receiver and a given name.
func (server *Server) RegisterName(name string, rcvr interface{}) error {
	return server.Server.RegisterName(name, rcvr)
}
