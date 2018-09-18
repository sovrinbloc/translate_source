package directory

import "testing"

func TestCopy(t *testing.T) {
	Copy("/Users/joealai/go/src/github.com/shen100/golang123/", "../translated_sources/golang123/")
}
