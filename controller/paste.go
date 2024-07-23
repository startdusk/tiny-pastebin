package controller

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
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
	languages := lexers.Names(true)
	h.view.Render(w, view.IndexPage, map[string]any{
		"pastes":    pastes,
		"languages": languages,
	})
}

func (h *PasteHandler) CreatePaste(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	body := r.FormValue("body")
	if strings.TrimSpace(body) == "" {
		h.Index(w, r)
		return
	}

	language := strings.TrimSpace(r.FormValue("language"))
	if language == "" {
		h.Index(w, r)
		return
	}
	code, err := gonanoid.New()
	if err != nil {
		log.Println("generate nanoid error:", err)
		http.Error(w, err.Error(),
			http.StatusUnprocessableEntity)
		return
	}

	ctx := r.Context()
	paste, err := h.postgres.CreatePaste(ctx, model.CreatePaste{
		Code:     code,
		Body:     body,
		Language: language,
		Hash:     GetMD5Hash(body),
	})
	if err != nil {
		log.Println("create paste error:", err)
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
	if err = quick.Highlight(bw, paste.Body, paste.Language, "html", "vs"); err != nil {
		http.Error(w, err.Error(),
			http.StatusUnprocessableEntity)
		return
	}
	bw.Flush()
	err = h.view.Render(w, view.PastePage, map[string]any{
		"paste":   paste,
		"content": buf.String(),
	})
	if err != nil {
		http.Error(w, err.Error(),
			http.StatusUnprocessableEntity)
		return
	}
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
