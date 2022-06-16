package TRPcG

import (
	"log"
	"net"
	"net/rpc"

	"github.com/mizumoto-cn/TRPcG/codec"
	"github.com/mizumoto-cn/TRPcG/serializer"
)

// wrap /net/rpc :: Server

type Server struct {
	*rpc.Server
	serializer.Serializer
}

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
