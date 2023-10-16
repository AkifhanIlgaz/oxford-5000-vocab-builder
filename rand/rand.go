package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)

	nRead, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("rand bytes: %w", err)
	}

	if nRead < n {
		return nil, fmt.Errorf("bytes didn't read enough random bytes")
	}

	return b, nil
}

func String(bytesPerToken int) (string, error) {
	b, err := Bytes(bytesPerToken)
	if err != nil {
		return "", fmt.Errorf("rand string: %w", err)
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
