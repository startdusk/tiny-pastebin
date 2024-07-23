package model

import (
	"context"
	"time"
)

type Paste struct {
	ID        int       `json:"-" db:"id"`
	Code      string    `json:"code" db:"code"`
	Body      string    `json:"body" db:"body"`
	Language  string    `json:"language" db:"language"`
	Hash      string    `json:"-" db:"hash"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreatePaste struct {
	Code     string
	Body     string
	Language string
	Hash     string
}

const createPasteQuery = `
	INSERT INTO paste (code, body, language, hash) VALUES ($1,$2,$3,$4) 
	ON CONFLICT(hash) DO UPDATE SET updated_at=NOW() RETURNING *
`

func (p *Postgres) CreatePaste(ctx context.Context, cp CreatePaste) (Paste, error) {
	var paste Paste
	err := p.db.QueryRowxContext(ctx, createPasteQuery, cp.Code, cp.Body, cp.Language, cp.Hash).StructScan(&paste)
	return paste, err
}

const getPasteByCodeQuery = `
	SELECT id, code, body, language, created_at ,updated_at FROM paste WHERE code = $1
`

func (p *Postgres) GetPasteByCode(ctx context.Context, code string) (Paste, error) {
	var paste Paste
	err := p.db.GetContext(ctx, &paste, getPasteByCodeQuery, code)
	return paste, err
}

const latestPasteByCodeQuery = `
	SELECT id, code, language, created_at, updated_at FROM paste ORDER BY created_at DESC LIMIT $1
`

func (p *Postgres) LatestPaste(ctx context.Context, limit uint) ([]Paste, error) {
	var pastes []Paste
	err := p.db.SelectContext(ctx, &pastes, latestPasteByCodeQuery, limit)
	return pastes, err
}
