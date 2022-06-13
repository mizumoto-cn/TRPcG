package codec

import (
	"encoding/binary"
	"io"
	"net"
)

func sendFrame(w io.Writer, data []byte) (err error) {
	var buf [binary.MaxVarintLen64]byte

	if data == nil || len(data) == 0 {
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
		if _, flag := err.(net.Error); !flag {
			return err
		}
		i += n
	}
	return nil
}
