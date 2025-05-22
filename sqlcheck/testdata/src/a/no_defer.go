package a

import (
	"context" // Added import
	"database/sql"
	"log"
)

func NoDefer(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM sessions") // want "immediately defer close the rows returned here to avoid a memory leak"
	if err != nil {
		log.Fatal(err)
	}

	// Process rows
	for rows.Next() {
		// ...
	}
	rows.Close() // Close is called, but not deferred
}

func NoDeferQueryContext(db *sql.DB) {
	rs, err := db.QueryContext(context.Background(), "SELECT * FROM events") // want "immediately defer close the rows returned here to avoid a memory leak"
	if err != nil {
		log.Fatal(err)
	}

	for rs.Next() {
		// ...
	}
	rs.Close() // Close is called, but not deferred
}
