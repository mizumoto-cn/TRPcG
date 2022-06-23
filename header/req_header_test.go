package header

import (
	"reflect"
	"testing"

	"github.com/mizumoto-cn/TRPcG/compressor"
	"github.com/stretchr/testify/assert"
)

// TestRequestHeader_Marshal tests RequestHeader::Marshal
func TestRequestHeader_Marshal(t *testing.T) {
	header := &RequestHeader{
		CompressType: compressor.Raw,
		Method:       "Add",
		ID:           12345,
		RequestLen:   123,
		Checksum:     12345,
	}
	assert.Equal(t, []byte{0x0, 0x0, 0x3, 0x41, 0x64, 0x64, 0xb9, 0x60, 0x7b, 0x39, 0x30, 0x0, 0x0}, header.Marshal())
}

// TestRequestHeader_Unmarshal tests RequestHeader::Unmarshal
func TestRequestHeader_Unmarshal(t *testing.T) {
	type expect struct {
		header *RequestHeader
		err    error
	}
	cases := []struct {
		name   string
		data   []byte
		expect expect
	}{
		{
			"test-1",
			[]byte{0x0, 0x0, 0x3, 0x41, 0x64, 0x64, 0xb9, 0x60, 0x7b, 0x39, 0x30, 0x0, 0x0},
			expect{
				&RequestHeader{
					CompressType: compressor.Raw,
					Method:       "Add",
					ID:           12345,
					RequestLen:   123,
					Checksum:     12345,
				},
				nil,
			},
		},

		{
			"test-2",
			[]byte{0x0},
			expect{&RequestHeader{}, ErrUnmarshalFail},
		},
		{
			"test-3",
			nil,
			expect{&RequestHeader{}, ErrUnmarshalFail},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := &RequestHeader{}
			err := h.UnMarshal(c.data)
			assert.Equal(t, true, reflect.DeepEqual(c.expect.header, h))
			assert.Equal(t, err, c.expect.err)
		})
	}
}

// TestRequestHeader_ResetHeader tests RequestHeader::ResetHeader
func TestRequestHeader_ResetHeader(t *testing.T) {
	header := &RequestHeader{
		CompressType: compressor.Raw,
		Method:       "Add",
		ID:           12345,
		RequestLen:   123,
		Checksum:     12345,
	}
	header.ResetHeader()
	assert.Equal(t, &RequestHeader{}, header)
}

// TestRequestHeader_GetCompressType tests RequestHeader::GetCompressType
func TestRequestHeader_GetCompressType(t *testing.T) {
	header := &RequestHeader{
		CompressType: compressor.Raw,
		Method:       "Add",
		ID:           12345,
		RequestLen:   123,
		Checksum:     12345,
	}
	assert.Equal(t, compressor.CompressType(compressor.Raw), header.GetCompressType())
}
