package compressor

// implements Compressor interface
type RawCompressor struct {
}

// Zip
func (_ RawCompressor) Zip(data []byte) ([]byte, error) {
	return data, nil
}

// Unzip
func (_ RawCompressor) Unzip(data []byte) ([]byte, error) {
	return data, nil
}
