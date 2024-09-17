package main

import (
	"fmt"
	"net/url"
	"strings"
)

func Normalize_url(URL string) (string, error) {
	res, err := url.Parse(URL)
	if err != nil {
		return "", fmt.Errorf("couldn't parse Url %w", err)
	}
	out := strings.ToLower(res.Host + res.Path)

	out = strings.TrimSuffix(out, "/")
	return out, nil
}
