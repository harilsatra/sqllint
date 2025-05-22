package a

import (
	"context" // Added import
	"database/sql"
	"log"
)

func MissingClose(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM users") // want "immediately defer close the rows returned here to avoid a memory leak"
	if err != nil {
		log.Fatal(err)
	}
	// defer rows.Close() is missing

	for rows.Next() {
		// ...
	}
}

func MissingCloseQueryContext(db *sql.DB) {
	results, err := db.QueryContext(context.Background(), "SELECT * FROM products") // want "immediately defer close the rows returned here to avoid a memory leak"
	if err != nil {
		log.Fatal(err)
	}
	// defer results.Close() is missing

	for results.Next() {
		// ...
	}
}
