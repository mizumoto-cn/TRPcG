package TRPcG

import (
	"encoding/json"
	"log"
	"net"
	"testing"

	"github.com/mizumoto-cn/TRPcG/compressor"
	jsonp "github.com/mizumoto-cn/TRPcG/testing/json"
	message "github.com/mizumoto-cn/TRPcG/testing/message"
	"github.com/stretchr/testify/assert"
)

// Json
type Json struct{}

// Json::Marshal
func (j *Json) Marshal(message any) ([]byte, error) {
	return json.Marshal(message)
}

// Json::Unmarshal
func (j *Json) Unmarshal(data []byte, message any) error {
	return json.Unmarshal(data, message)
}

func init() {
	listen, err := net.Listen("tcp", ":8008")
	if err != nil {
		log.Fatal("listen error:", err)
	}

	server := NewServer()
	err = server.Register(new(message.ArithService))
	if err != nil {
		log.Fatal("register error:", err)
	}
	go server.Serve(listen)

	listen, err = net.Listen("tcp", ":8009")
	if err != nil {
		log.Fatal("listen error:", err)
	}

	server = NewServer(WithSerializer(&Json{}))
	err = server.Register(new(jsonp.TestService))
	if err != nil {
		log.Fatal("register error:", err)
	}
	go server.Serve(listen)
}

// Test_Client_Call tests the synchronous call of the client.
func Test_Client_Call(t *testing.T) {
	compressType := compressor.Gzip
	conn, err := net.Dial("tcp", ":8008")
	if err != nil {
		log.Fatal("dial error:", err)
		// t.Fatal("dial error:", err)
	}
	defer conn.Close()

	client := NewClient(conn, WithCompress(compressType))
	defer client.Close()

	type expect struct {
		reply *message.ArithResponse
		err   error
	}

	cases := []struct {
		client        *Client
		name          string
		serviceMethod string
		args          *message.ArithRequest
		expect        expect
	}{
		{
			client,
			"Arith-1",
			"ArithService.Add",
			&message.ArithRequest{A: 1, B: 2},
			expect{
				&message.ArithResponse{C: 3},
				nil,
			},
		},
		{
			client,
			"Arith-2",
			"ArithService.Sub",
			&message.ArithRequest{A: 1, B: 2},
			expect{
				&message.ArithResponse{C: -1},
				nil,
			},
		},
		{
			client,
			"Arith-3",
			"ArithService.Mul",
			&message.ArithRequest{A: 1, B: 2},
			expect{
				&message.ArithResponse{C: 2},
				nil,
			},
		},
		{
			client,
			"Arith-4",
			"ArithService.Div",
			&message.ArithRequest{A: 6, B: 2},
			expect{
				&message.ArithResponse{C: 3},
				nil,
			},
		},
		// {
		// 	client,
		// 	"Arith-5",
		// 	"ArithService.Div",
		// 	&message.ArithRequest{A: 1, B: 0},
		// 	expect{
		// 		nil,
		// 		rpc.ServerError("divided by zero"),
		// 	},
		// },
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reply := &message.ArithResponse{}
			err := c.client.Call(c.serviceMethod, c.args, reply)
			assert.Equal(t, c.expect.reply.C, reply.C)
			assert.Equal(t, c.expect.err, err)
		})
	}
}
