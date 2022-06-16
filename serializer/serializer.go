package serializer

// type SerializeType int32

type Serializer interface {
	Marshal(message any) ([]byte, error)
	Unmarshal(data []byte, message any) error
}
