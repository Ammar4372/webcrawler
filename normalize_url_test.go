package main

import (
	"strings"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      string
		errorContains string
	}{
		{
			name:          "Remove Scheme",
			input:         "https://blog.boot.dev/path",
			expected:      "blog.boot.dev/path",
			errorContains: "",
		},
		{
			name:          "Remove Upper Case",
			input:         "https://Blog.boot.dev/path",
			expected:      "blog.boot.dev/path",
			errorContains: "",
		},
		{
			name:          "Remove Query String",
			input:         "https://Blog.boot.dev/path?test=1",
			expected:      "blog.boot.dev/path",
			errorContains: "",
		},
		{
			name:          "Remove Trailing Slash",
			input:         "https://Blog.boot.dev/path/",
			expected:      "blog.boot.dev/path",
			errorContains: "",
		},
		{
			name:          "Wrong Url",
			input:         "://Blog.boot.dev/path/",
			expected:      "",
			errorContains: "couldn't parse Url",
		},
	}
	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := Normalize_url(tc.input)
			if err != nil && !strings.Contains(err.Error(), tc.errorContains) {
				t.Errorf("Test %v - '%s' FAIL: unexcepted error %v", i, tc.name, err)
				return
			} else if err != nil && tc.errorContains == "" {
				t.Errorf("Test %v - '%s' FAIL: unexcepted error %v", i, tc.name, err)
				return
			} else if err == nil && tc.errorContains != "" {
				t.Errorf("Test %v - '%s' FAIL: expected error containing '%v', got none.", i, tc.name, tc.errorContains)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - '%s' FAIL: excepted: %s got: %s", i, tc.name, tc.expected, actual)
			}
		})
	}

}
