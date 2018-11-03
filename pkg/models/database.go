package models

import (
	"database/sql"
	"time"
)

type Database struct {
	*sql.DB
}

func (db *Database) GetSnippet(id int) (*Snippet, error) {
	if id == 123 {
		snippet := &Snippet{
			ID:      id,
			Title:   "Example title",
			Content: "Example content",
			Created: time.Now(),
			Expires: time.Now(),
		}
		return snippet, nil
	}

	return nil, nil
}
