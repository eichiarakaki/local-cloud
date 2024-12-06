package src

import (
	"errors"
	"strings"
)

func URLFilter(data string) (string, error) {
	url := strings.TrimSpace(data)
	if strings.HasPrefix(url, "https://www.youtube.com/") || strings.HasPrefix(url, "https://youtu.be/") {
		return url, nil
	}
	return "", errors.New("URL invalid.")
}
