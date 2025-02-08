package utils

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

func ShortenURL(ctx context.Context, url string, size int) (string, error) {
	hash := sha256.Sum256([]byte(url))
	hashHex := hex.EncodeToString(hash[:])

	for len(hashHex) < size {
		hashHex += hashHex
	}

	if len(hashHex) > size {
		hashHex = hashHex[:size]
	}

	if strings.HasPrefix(url, "https://") {
		return fmt.Sprintf("https://%s", hashHex), nil
	}

	if strings.HasPrefix(url, "http://") {
		return fmt.Sprintf("http://%s", hashHex), nil
	}

	return hashHex, nil
}
