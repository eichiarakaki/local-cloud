package src

import (
	"errors"
	"strings"
)

func URLFilter(data string) (string, error) {
	URLPrefixes := []string{
		"https://www.youtube.com/",
		"https://youtu.be/",
	}
	url := strings.TrimSpace(data)

	for i := 0; i < len(URLPrefixes); i++ {
		if strings.HasPrefix(url, URLPrefixes[i]) {
			return url, nil
		}
	}

	return "", errors.New("URL invalid.")
}
