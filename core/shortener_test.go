package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShortenURL(t *testing.T) {
	tests := []struct {
		name    string
		longURL string
	}{
		{
			name:    "shorten chatgpt url",
			longURL: "https://chatgpt.com/c/6a5baca9-ed74-83e8-9b49-881a49bd8d7c",
		},
	}

	var url string
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url = ShortenURL(test.longURL)

			require.NotEqual(t, url, test.longURL)
			// https://short.ly/lZyD
			require.LessOrEqual(t, 17+4, len(url))
			require.GreaterOrEqual(t, 17+14, len(url))
		})
	}
}

func TestURLIsCollisionFree(t *testing.T) {
	test_long_url := "https://chatgpt.com/c/6a5baca9-ed74-83e8-9b49-881a49bd8d7c"

	results := []string{}
	for range 1000 {
		results = append(results, ShortenURL(test_long_url))
	}

	isCollisionFree := true
	for i, r1 := range results {
		for j, r2 := range results {
			if i != j && r1 == r2 {
				isCollisionFree = false
				break
			}
		}
	}

	require.True(t, isCollisionFree)
}
