package controller

import (
	"bufio"
	"bytes"
	"log"
	"testing"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
)

func TestHighlightCode(t *testing.T) {
	code := `package main

func main() { }
`

	var buf bytes.Buffer
	bw := bufio.NewWriter(&buf)
	err := quick.Highlight(bw, code, "go", "html", "monokai")
	if err != nil {
		log.Fatal(err)
	}

	bw.Flush()
	log.Printf("%s", buf.String())

	languages := lexers.Names(true)
	log.Println(languages)
}
