package utils

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

func ShortenURL(ctx context.Context, url string, size int) (string, error) {
	// check is optional, potentially may be possible to implement Redirect, so that is why a valid url is required
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		return "", NewError("unknown url format: expect http or https", BadRequestErr)
	}

	hash := makeHashFromString(ctx, url, size)

	return hash, nil
}

func makeHashFromString(ctx context.Context, str string, size int) string {
	hash := sha256.Sum256([]byte(str))
	hashHex := hex.EncodeToString(hash[:])

	for len(hashHex) < size {
		hashHex += hashHex
	}

	if len(hashHex) > size {
		hashHex = hashHex[:size]
	}

	return hashHex
}
