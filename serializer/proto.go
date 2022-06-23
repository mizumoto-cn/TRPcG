package serializer

import (
	"errors"

	"google.golang.org/protobuf/proto"
)

var ErrProtoMsgNotImplemented = errors.New("param does not implement proto.Message")

// implements Serializer interface
type ProtoSerializer struct {
}

var Proto = ProtoSerializer{}

// Marshal
func (_ ProtoSerializer) Marshal(message any) ([]byte, error) {
	if message == nil {
		return []byte{}, nil
	}
	//Message is the top-level interface that all messages must implement.
	// It provides access to a reflective view of a message. Any implementation of this interface
	// may be used with all functions in the protobuf module that accept a Message, except where otherwise specified.
	body, ok := message.(proto.Message)
	if !ok {
		return nil, ErrProtoMsgNotImplemented
	}
	return proto.Marshal(body)
}

// Unmarshal
func (_ ProtoSerializer) Unmarshal(data []byte, message any) error {
	if message == nil {
		return nil
	}

	body, ok := message.(proto.Message)
	if !ok {
		return ErrProtoMsgNotImplemented
	}
	// func Unmarshal(b []byte, m Message) error
	// parses the wire-format message in b and places the result in m.
	return proto.Unmarshal(data, body)
}
