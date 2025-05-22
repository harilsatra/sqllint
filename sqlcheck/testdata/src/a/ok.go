package a

import (
	"context" // Added import
	"database/sql"
	"log"
)

func GoodQuery(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() // This is correct

	// Process rows
	for rows.Next() {
		// ...
	}
}

func GoodQueryContext(db *sql.DB) {
	rows, err := db.QueryContext(context.Background(), "SELECT * FROM products")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() // This is correct

	for rows.Next() {
		// ...
	}
}

func AnotherGoodQuery(db *sql.DB) {
	// Different variable name
	results, err := db.Query("SELECT * FROM orders")
	if err != nil {
		log.Fatal(err)
	}
	defer results.Close() // This is correct

	for results.Next() {
		// ...
	}
}

func QueryWithBlankIdentifier(db *sql.DB) {
	_, err := db.Query("SELECT name FROM items WHERE id = ?", 1) // Rows not assigned to a variable
	if err != nil {
		log.Fatal(err)
	}
	// No rows variable to close, so this is fine.
}

func QueryInShortAssignment(db *sql.DB) {
	// This test case is to ensure Query calls that are not part of an assignment
	// are not flagged, as there's no rows variable to track.
	// However, the current linter logic specifically looks for *ast.AssignStmt.
	// If Query() is called without assignment, it won't be picked up by the linter.
	// This is acceptable given the current linter implementation.
	db.Query("UPDATE users SET last_login = NOW() WHERE id = 1")
}
