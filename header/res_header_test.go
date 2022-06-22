package header

import (
	"reflect"
	"testing"

	"github.com/mizumoto-cn/TRPcG/compressor"
	"github.com/stretchr/testify/assert"
)

// TestResponseHeader_Marshal tests ResponseHeader::Marshal
func TestResponseHeader_Marshal(t *testing.T) {
	header := &ResponseHeader{
		CompressType: CompressType(compressor.Raw),
		Error:        "error",
		ID:           12345,
		ResponseLen:  123,
		CheckSum:     12345,
	}
	assert.Equal(t, []byte{0x0, 0x0, 0xb9, 0x60, 0x5, 0x65, 0x72, 0x72, 0x6f,
		0x72, 0x7b, 0x39, 0x30, 0x0, 0x0}, header.Marshal())
}

// TestResponseHeader_Unmarshal tests ResponseHeader::Unmarshal
func TestResponseHeader_Unmarshal(t *testing.T) {
	type expect struct {
		header *ResponseHeader
		err    error
	}
	cases := []struct {
		name   string
		data   []byte
		expect expect
	}{
		{
			"test-1",
			[]byte{0x0, 0x0, 0xb9, 0x60, 0x5, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x7b, 0x39, 0x30, 0x0, 0x0},
			expect{
				&ResponseHeader{
					CompressType: CompressType(compressor.Raw),
					Error:        "error",
					ID:           12345,
					ResponseLen:  123,
					CheckSum:     12345,
				},
				nil,
			},
		},
		{
			"test-2",
			[]byte{0x0},
			expect{&ResponseHeader{}, ErrUnmarshalFail},
		},
		{
			"test-3",
			nil,
			expect{&ResponseHeader{}, ErrUnmarshalFail},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			header := &ResponseHeader{}
			err := header.Unmarshal(c.data)
			assert.Equal(t, c.expect.err, err)
			assert.Equal(t, true, reflect.DeepEqual(c.expect.header, header))
		})
	}
}

// TestResponseHeader_ResetHeader tests ResponseHeader::ResetHeader
func TestResponseHeader_ResetHeader(t *testing.T) {
	header := &ResponseHeader{
		CompressType: CompressType(compressor.Raw),
		Error:        "error",
		ID:           12345,
		ResponseLen:  123,
		CheckSum:     12345,
	}
	header.ResetHeader()
	assert.Equal(t, true, reflect.DeepEqual(&ResponseHeader{}, header))
}

// TestResponseHeader_GetCompressType tests ResponseHeader::GetCompressType
func TestResponseHeader_GetCompressType(t *testing.T) {
	header := &ResponseHeader{
		CompressType: CompressType(compressor.Raw),
		Error:        "error",
		ID:           12345,
		ResponseLen:  123,
		CheckSum:     12345,
	}
	assert.Equal(t, true, reflect.DeepEqual(compressor.Raw, header.GetCompressType()))
}
