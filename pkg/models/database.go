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

func (db *Database) InsertSnippet(title, content, expires string) (int, error) {
	statement := `
		INSERT INTO snippets (title, content, created, expires)
		VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? SECOND))
	`

	// This returns a sql.Result object which contains some basic info
	// about what happened when the statement was executed
	result, err := db.Exec(statement, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Not all drivers support the LastInsertId method (e.g. Postgres)
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// ID returned by LastInsertId() is of type int64, so convert it to an int type
	return int(id), nil
}

// Use `transaction` to:
// 1. guarantee the same database connection is used
// 2. execute multiple SQL statements as a single atomic action
// So long as we use tx.Rollback() in the event of any error, tx ensures either:
// a. ALL statements are executed successfully; OR
// b. NO statements are executed and the db remains unchanged

// Example 1
// func ExampleTransaction(db *sql.DB) error {
// 	// Calling the Begin method on the connection pool creates a new sql.Tx object
// 	// which represents the in-progress database connection
// 	tx, err := db.Begin()
// 	if err != nil {
// 		return err
// 	}

// 	// Call Exec() on the transaction, passing in your statement and any
// 	// parameters. It's important to notice that tx.Exec() is called on the
// 	// transaction object we just created, NOT the connection pool. Although we're
// 	// using tx.Exec() here you can also use tx.Query() and tx.QueryRow() in
// 	// exactly the same way.
// 	_, err = tx.Exec("INSERT INTO ...")
// 	if err != nil {
// 		// If there is any error, we call the tx.Rollback() method on the
// 		// transaction. This will abort the transaction and no changes will be
// 		// made to the database.
// 		tx.Rollback()
// 		return err
// 	}
// 	// Carry out another transaction in exactly the same way.
// 	_, err = tx.Exec("UPDATE ...")
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}
// 	// If there are no errors, the statements in transaction can be commited
// 	// to the database with the tx.Commit() method. It's really important to ALWAYS
// 	// call either Rollback() or Commit() before your function returns. If you
// 	// don't the connection will stay open and not be returned to the connection
// 	// pool. This can lead to hitting you maximum connection limit/running out of
// 	// resources.
// 	err = tx.Commit()
// 	return err
// }
