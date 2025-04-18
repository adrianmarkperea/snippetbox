package models

import (
	"database/sql"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	sql := `INSERT INTO snippets (title, content, created, expires)
		      VALUES ($1, $2, CURRENT_TIMESTAMP AT TIME ZONE 'UTC', (CURRENT_TIMESTAMP AT TIME ZONE 'UTC') + make_interval(days => $3))
		      RETURNING id`

	var id int
	err := m.DB.QueryRow(sql, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *SnippetModel) Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
