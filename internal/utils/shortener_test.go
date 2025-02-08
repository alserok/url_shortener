package utils

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestShortener(t *testing.T) {
	tests := [][2]string{
		{"https://www.youtube.com/", "https://dbae2d0204"},
		{"https://github.com/", "https://09a8b930c8"},
		{"https://gitlab.com/", "https://9931966bed"},
		{"http://my_addr.com/", "http://c46ce31003"},
	}

	for _, tc := range tests {
		shortened, err := ShortenURL(context.Background(), tc[0], 10)
		require.NoError(t, err)
		require.Equal(t, tc[1], shortened)
	}
}
