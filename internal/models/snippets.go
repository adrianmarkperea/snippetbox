package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (Snippet, error)
	Latest() ([]Snippet, error)
}

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *pgxpool.Pool
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
		      VALUES ($1, $2, CURRENT_TIMESTAMP AT TIME ZONE 'UTC', (CURRENT_TIMESTAMP AT TIME ZONE 'UTC') + make_interval(days => $3))
		      RETURNING id`

	var id int
	err := m.DB.QueryRow(context.Background(), stmt, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
           WHERE expires > (CURRENT_TIMESTAMP AT TIME ZONE 'UTC') AND id = $1`

	var s Snippet

	if err := m.DB.QueryRow(context.Background(), stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	         WHERE expires > (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
	         ORDER BY id DESC LIMIT 10`
	rs, err := m.DB.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	var ss []Snippet
	for rs.Next() {
		var s Snippet
		if err := rs.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires); err != nil {
			return nil, err
		}
		ss = append(ss, s)
	}

	// don't assume everything went well
	if err := rs.Err(); err != nil {
		return nil, err
	}

	return ss, nil
}
