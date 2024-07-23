package view

import (
	"fmt"
	"io"

	"github.com/flosch/pongo2/v6"
)

type PasteView struct {
	tpls map[string]*pongo2.Template
}

const IndexPage = "index.html"
const PastePage = "paste.html"
const NotFoundPage = "404.html"

func CreatePasteView(root string) PasteView {
	pv := PasteView{
		tpls: make(map[string]*pongo2.Template),
	}

	pv.tpls[IndexPage] = pongo2.Must(pongo2.FromFile(fmt.Sprintf("%s/%s", root, IndexPage)))
	pv.tpls[PastePage] = pongo2.Must(pongo2.FromFile(fmt.Sprintf("%s/%s", root, PastePage)))
	pv.tpls[NotFoundPage] = pongo2.Must(pongo2.FromFile(fmt.Sprintf("%s/%s", root, NotFoundPage)))

	return pv
}

func (v PasteView) Render(w io.Writer, filename string, data map[string]any) error {
	tpl, ok := v.tpls[filename]
	if !ok {
		return fmt.Errorf("template %s not found", filename)
	}
	return tpl.ExecuteWriter(data, w)
}
