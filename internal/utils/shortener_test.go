package utils

import (
	"context"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestShortener(t *testing.T) {
	tests := [][2]string{
		{"https://www.youtube.com/", "dbae2d0204"},
		{"https://github.com/", "09a8b930c8"},
		{"https://gitlab.com/", "9931966bed"},
		{"http://my_addr.com/", "c46ce31003"},
		{"http://example.com/" + strings.Repeat("a", 200), "4261ee63a9"},
		{"https://example.com/" + strings.Repeat("a", 200), "fbe873ba66"},
		{"https://example.com/" + strings.Repeat("a", 2000), "0a3614a41e"},
		{"http://a.com/", "e649b13c1c"},
		{"http://a.com/", "e649b13c1c"},
		{"invalid_url", ""},
		{"https://secure-site.com/path/to/resource", "f8bf80a547"},
		{"htp://invalid-url", ""},
		{"http://example.com/query?param1=test&param2=sample", "6a9ec101c8"},
	}

	for _, tc := range tests {
		shortened, err := ShortenURL(context.Background(), tc[0], 10)
		if shortened == "" {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}

		require.Equal(t, tc[1], shortened, tc[0])
	}
}
