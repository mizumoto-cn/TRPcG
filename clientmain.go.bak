package main

import (
	"log"
	"net"

	"github.com/mizumoto-cn/TRPcG/testing/message"

	"github.com/mizumoto-cn/TRPcG"
)

func main(){
    conn, err := net.Dial("tcp", ":8082")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    client := TRPcG.NewClient(conn)
    resq := message.ArithRequest{A: 20, B: 5}
    resp := message.ArithResponse{}
    err = client.Call("ArithService.Add", &resq, &resp)
    log.Printf("Arith.Add(%v, %v): %v ,Error: %v", resq.A, resq.B, resp.C, err)
}