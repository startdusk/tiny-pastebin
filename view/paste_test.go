package view

import (
	"log"
	"testing"

	"github.com/kluctl/go-jinja2"
)

func TestPaste(t *testing.T) {
	j2, err := jinja2.NewJinja2("example", 1,
		jinja2.WithGlobal("test_var1", 1),
		jinja2.WithGlobal("test_var2", map[string]any{"test": 2}))
	if err != nil {
		panic(err)
	}
	defer j2.Close()

	template := "{{ test_var2 }}"

	s, err := j2.RenderString(template)
	if err != nil {
		panic(err)
	}

	log.Printf("template: %s\nresult: %s", template, s)
}
