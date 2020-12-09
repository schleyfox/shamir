package shares

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

// Encode prepends a little-endian crc32 checksum to `byteShare` and converts
// to a base64 encoded string.
//
// Format of the result is base64(littleEndian(crc32(byteShare)) + byteShare)
func Encode(byteShare []byte) string {
	val := make([]byte, 4, 4+len(byteShare))
	crc := crc32.ChecksumIEEE(byteShare)
	binary.LittleEndian.PutUint32(val, crc)
	val = append(val, byteShare...)

	return base64.StdEncoding.EncodeToString(val)
}

// Decode takes a `share` string in the format produced by `Encode`, verifies
// the checksum, and converts into the binary `byteShare` without the checksum.
func Decode(share string) ([]byte, error) {
	checksumShare, err := base64.StdEncoding.DecodeString(share)
	if err != nil {
		return nil, fmt.Errorf("could not decode share: %w", err)
	}

	if len(checksumShare) < 4 {
		return nil, fmt.Errorf("decoded share too short (len=%d, expected >= 4)", len(checksumShare))
	}

	checksum := binary.LittleEndian.Uint32(checksumShare[:4])
	byteShare := checksumShare[4:]
	computedChecksum := crc32.ChecksumIEEE(byteShare)

	if checksum != computedChecksum {
		return nil, fmt.Errorf("checksums do not match (checksum=%d, computed=%d)", checksum, computedChecksum)
	}

	return byteShare, nil
}
