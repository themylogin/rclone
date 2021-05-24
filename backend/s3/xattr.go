// xattr

package s3

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
)

func xattrDigest(xattr []byte) string {
	h := sha1.New()
	h.Write(xattr)
	return hex.EncodeToString(h.Sum(nil))
}

func xattrEncode(xattr []byte) (string, error) {
	if xattr == nil {
		return "", nil
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	_, err := gz.Write(xattr)
	if err != nil {
		return "", err
	}

	if err = gz.Flush(); err != nil {
		return "", err
	}

	if err = gz.Close(); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func canBeMetadataXattr(encoded string) bool {
	// S3 max metadata size is 2048 bytes (keys and values)
	// 2048 - len("mtime") - 20 - len("xattr") = 2018
	return len(encoded) <= 2018
}

func decodeMetadataXattr(data string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	gr, err := gzip.NewReader(bytes.NewBuffer(decoded))
	defer gr.Close()

	decompressed, err := ioutil.ReadAll(gr)
	if err != nil {
		return nil, err
	}

	return decompressed, nil
}
