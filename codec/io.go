package codec

import (
	"encoding/binary"
	"io"
	"net"
)

// build buf (uvarint) to indicate the buffer (data slice) size
// then write data slice into stream
// wstream << data length | data []bytes
func sendFrame(w io.Writer, data []byte) (err error) {
	var buf [binary.MaxVarintLen64]byte

	if len(data) == 0 {
		n := binary.PutUvarint(buf[:], uint64(0))
		err = write(w, buf[:n])
		if err != nil {
			return
		}
		return
	}

	n := binary.PutUvarint(buf[:], uint64(len(data)))
	if err = write(w, buf[:n]); err != nil {
		return
	}
	if err = write(w, data); err != nil {
		return
	}
	return
}

func write(w io.Writer, data []byte) error {
	for i := 0; i < len(data); {
		// returns the number of bytes written from data (0 <= n <= len(data))
		// and any error encountered that caused the write to stop early
		n, err := w.Write(data[i:])

		// type assertion  value, ok := interface_name.(TypeName)
		if _, ok := err.(net.Error); !ok {
			return err
		}
		i += n
	}
	return nil
}

// do not use reference unless you want to change something
func receiveFrame(r io.Reader) (data []byte, err error) {
	// func ReadUvarint(r io.ByteReader) (uint64, error)
	// ReadUvarint reads an encoded unsigned integer from r and returns it as a uint64.
	// https://golang.google.cn/ref/spec#Type_assertions
	buf_len, err := binary.ReadUvarint(r.(io.ByteReader))
	if err != nil {
		return nil, err
	}
	if buf_len != 0 {
		data = make([]byte, buf_len)
		if err = read(r, data); err != nil {
			return nil, err
		}
	}
	return data, nil
}

func read(r io.Reader, data []byte) error {
	// return nil
	for i := 0; i < len(data); {
		n, err := r.Read(data[i:])
		if err != nil {
			if _, ok := err.(net.Error); !ok {
				return err
			}
		}
		i += n
	}
	return nil
}

// cSpell:ignore wstream
