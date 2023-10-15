package json

import "errors"

// Request .
type Request struct {
	A float64 `json:"a,omitempty"`
	B float64 `json:"b,omitempty"`
}

// Response .
type Response struct {
	C float64 `json:"c,omitempty"`
}

// TestService Defining Computational Digital Services
type TestService struct{}

// Add addition
func (this *TestService) Add(args *Request, reply *Response) error {
	reply.C = args.A + args.B
	return nil
}

// Sub subtraction
func (this *TestService) Sub(args *Request, reply *Response) error {
	reply.C = args.A - args.B
	return nil
}

// Mul multiplication
func (this *TestService) Mul(args *Request, reply *Response) error {
	reply.C = args.A * args.B
	return nil
}

// Div division
func (this *TestService) Div(args *Request, reply *Response) error {
	if args.B == 0 {
		reply.C = 0
		return errors.New("divided by zero")
	}
	reply.C = args.A / args.B
	return nil
}
