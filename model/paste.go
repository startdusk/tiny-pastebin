package model

import "context"

type Paste struct {
	ID        int    `json:"-"`
	Code      string `json:"code"`
	Body      string `json:"body"`
	Language  string `json:"language"`
	Hash      string `json:"-"`
	CreatedAt int64  `json:"created_at"`
}

type CreatePaste struct {
	Code     string
	Body     string
	Language string
	Hash     string
}

const createPasteQuery = `
	INSERT INTO paste (code, body, language) VALUES (?,?,?) RETURNING *
`

func (p *Postgres) CreatePaste(ctx context.Context, cp CreatePaste) (Paste, error) {
	var paste Paste
	err := p.db.QueryRowContext(ctx, createPasteQuery, cp.Code, cp.Body, cp.Language).Scan(
		&paste.ID,
		&paste.Code,
		&paste.Body,
		&paste.Language,
		&paste.CreatedAt,
	)
	return paste, err
}

const getPasteByCodeQuery = `
	SELECT id, code, body, language, created_at FROM paste WHERE code = ?
`

func (p *Postgres) GetPasteByCode(ctx context.Context, code string) (Paste, error) {
	var paste Paste
	err := p.db.QueryRowContext(ctx, getPasteByCodeQuery, code).Scan(
		&paste.ID,
		&paste.Code,
		&paste.Body,
		&paste.Language,
		&paste.CreatedAt,
	)
	return paste, err
}

const latestPasteByCodeQuery = `
	SELECT id, code, language, created_at FROM paste ORDER BY created_at DESC LIMIT $1
`

func (p *Postgres) LatestPaste(ctx context.Context, limit uint) ([]Paste, error) {
	rows, err := p.db.QueryContext(ctx, latestPasteByCodeQuery, limit)
	if err != nil {
		return nil, err
	}
	var pastes []Paste
	for rows.Next() {
		var paste Paste
		if err := rows.Scan(
			&paste.ID,
			&paste.Code,
			&paste.Language,
			&paste.CreatedAt,
		); err != nil {
			return nil, err
		}
	}
	return pastes, rows.Err()
}
