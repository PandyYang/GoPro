package gee

import (
	"strings"
	"testing"
)

func TestParsePattern2(t *testing.T) {
	vs := strings.Split("/a/b/c/*d/e", "/")
	parts := make([]string, 0)

	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	for _, v := range parts {
		t.Log(v)
	}
}
