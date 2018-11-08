package models

import (
	"database/sql"
)

type Database struct {
	*sql.DB
}

func (db *Database) GetSnippet(id int) (*Snippet, error) {
	statement := `
		SELECT id, title, content, created, expires FROM snippets
		WHERE expires > UTC_TIMESTAMP() and id = ?
	`

	// This returns a pointer to a sql.Row object which holds the result returned by the db
	row := db.QueryRow(statement, id)

	// Initialize a pointer to a new zeroed Snippet struct
	s := &Snippet{}

	// Copy values from each field in sql.Row to corresponding fields in Snippet struct
	// The args to row.Scan() must be pointers to the place you want to copy the data into,
	// and the # of args must be exactly the same as the # of columns returned by the statement
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return s, nil
}

func (db *Database) LatestSnippets() (Snippets, error) {
	statement := `
		SELECT id, title, content, created, expires FROM snippets WHERE expires >
		UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10
	`

	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}

	// IMPORTANT!
	defer rows.Close()

	// Initialize empty Snippets object (a slice of the type []*Snippet)
	snippets := Snippets{}

	// Use rows.Next to iterate through the rows.
	// This prepares the first (and then each subsequent) row to be acted on
	// by the rows.Scan() method.
	for rows.Next() {
		s := &Snippet{}

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
