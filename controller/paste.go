package controller

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/jmoiron/sqlx"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/startdusk/tiny-pastebin/model"
	"github.com/startdusk/tiny-pastebin/view"
)

type PasteHandler struct {
	postgres *model.Postgres
	view     view.PasteView
}

func CreatePasteHandler(conn *sqlx.DB, view view.PasteView) (PasteHandler, error) {
	postgres, err := model.CreateDatabase(conn)
	return PasteHandler{
		postgres: postgres,
		view:     view,
	}, err
}

func (h *PasteHandler) Index(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	ctx := r.Context()
	pastes, err := h.postgres.LatestPaste(ctx, 5)
	if err != nil {
		http.Error(w, err.Error(),
			http.StatusBadRequest)
		return
	}

	h.view.Render(w, view.IndexPage, map[string]any{
		"pastes":    pastes,
		"languages": "",
	})
}

func (h *PasteHandler) CreatePaste(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	body := r.FormValue("body")
	if body != "" {
		h.Index(w, r)
		return
	}
	code, err := gonanoid.New()
	if err != nil {
		http.Error(w, err.Error(),
			http.StatusUnprocessableEntity)
		return
	}

	language := lexers.Analyse(body).Config().Name
	ctx := r.Context()
	paste, err := h.postgres.CreatePaste(ctx, model.CreatePaste{
		Code:     code,
		Body:     body,
		Language: language,
		Hash:     "123456",
	})
	if err != nil {
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/%s", paste.Code), http.StatusFound)
}

func (h *PasteHandler) GetPaste(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimSpace(r.PathValue("code"))
	if code == "" {
		http.Error(w, "Invalid url",
			http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	paste, err := h.postgres.GetPasteByCode(ctx, code)
	if err != nil {
		http.Error(w, err.Error(),
			http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	bw := bufio.NewWriter(&buf)
	if err = quick.Highlight(bw, paste.Code, paste.Language, "html", "monokai"); err != nil {
		http.Error(w, err.Error(),
			http.StatusUnprocessableEntity)
		return
	}
	bw.Flush()
	paste.Body = buf.String()
	err = h.view.Render(w, view.PastePage, map[string]any{
		"paste": paste,
	})
	if err != nil {
		http.Error(w, err.Error(),
			http.StatusUnprocessableEntity)
		return
	}
}
