package a

import (
	"database/sql"
	"log"
)

func WrongVariableClosed(db *sql.DB) {
	rows1, err1 := db.Query("SELECT * FROM table1") // want "immediately defer close the rows returned here to avoid a memory leak"
	if err1 != nil {
		log.Fatal(err1)
	}
	// Missing: defer rows1.Close()

	rows2, err2 := db.Query("SELECT * FROM table2")
	if err2 != nil {
		log.Fatal(err2)
	}
	defer rows2.Close() // Correct for rows2, but rows1 is still open

	for rows1.Next() {
		// process rows1
	}
	for rows2.Next() {
		// process rows2
	}
}

func DifferentScopes(db *sql.DB) {
	outerRows, err := db.Query("SELECT id FROM records") // want "immediately defer close the rows returned here to avoid a memory leak"
	if err != nil {
		log.Fatal(err)
	}
	// defer outerRows.Close() is missing

	if true {
		innerRows, err_inner := db.Query("SELECT name FROM details")
		if err_inner != nil {
			log.Fatal(err_inner)
		}
		defer innerRows.Close() // Correct for innerRows
		for innerRows.Next() {
			// ...
		}
	}

	for outerRows.Next() {
		// ...
	}
}
