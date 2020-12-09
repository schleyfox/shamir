package shares

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"
)

var sharedCases = []struct {
	name      string
	byteShare []byte
	share     string
}{
	{
		name:      "empty byteShare",
		byteShare: []byte{},
		share:     "AAAAAA==",
	},
	{
		name:      "low binary",
		byteShare: []byte{1, 2, 3},
		share:     "HYC8VQECAw==",
	},
	{
		name:      "high binary",
		byteShare: []byte{253, 254, 255, 1},
		share:     "WIIxDv3+/wE=",
	},
	{
		name:      "arbitrary string foobar",
		byteShare: []byte("foobar"),
		share:     "lR/2nmZvb2Jhcg==",
	},
}

func TestEncode(t *testing.T) {
	for _, testCase := range sharedCases {
		t.Run(testCase.name, func(t *testing.T) {
			encoded := Encode(testCase.byteShare)
			if encoded != testCase.share {
				t.Fatalf("expected: %v, got: %v", testCase.share, encoded)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	for _, testCase := range sharedCases {
		t.Run(testCase.name, func(t *testing.T) {
			decoded, err := Decode(testCase.share)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !bytes.Equal(decoded, testCase.byteShare) {
				t.Fatalf("expected: %v, got: %v", testCase.byteShare, decoded)
			}
		})
	}
}

func TestDecode_invalid(t *testing.T) {
	highBinary, err := base64.StdEncoding.DecodeString("WIIxDv3+/wE=")
	if err != nil {
		t.Fatalf("%v", err)
	}
	corruptChecksum := make([]byte, len(highBinary))
	corruptShare := make([]byte, len(highBinary))
	copy(highBinary, corruptChecksum)
	copy(highBinary, corruptShare)
	corruptChecksum[0]++
	corruptShare[len(corruptShare)-1]++

	testCases := []struct {
		name      string
		share     string
		errPrefix string
	}{
		{
			name:      "invalid base64",
			share:     "R/2nm==",
			errPrefix: "could not decode share",
		},
		{
			name:      "empty string",
			share:     "",
			errPrefix: "decoded share too short",
		},
		{
			name:      "byte string 12",
			share:     "MTI=",
			errPrefix: "decoded share too short",
		},
		{
			name:      "byte string 123",
			share:     "MTIz",
			errPrefix: "decoded share too short",
		},
		{
			name:      "byte string 1234",
			share:     "MTIzNA==",
			errPrefix: "checksums do not match",
		},
		{
			name:      "byte string 12345",
			share:     "MTIzNDU=",
			errPrefix: "checksums do not match",
		},
		{
			name:      "corrupt checksum",
			share:     base64.StdEncoding.EncodeToString(corruptChecksum),
			errPrefix: "checksums do not match",
		},
		{
			name:      "corrupt share",
			share:     base64.StdEncoding.EncodeToString(corruptShare),
			errPrefix: "checksums do not match",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := Decode(testCase.share)
			if err == nil || !strings.HasPrefix(err.Error(), testCase.errPrefix) {
				t.Fatalf("expected error with prefix: %v, got: %v", testCase.errPrefix, err)
			}
		})
	}
}
