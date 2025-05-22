package a

import (
	"database/sql"
	"log"
)

func MultipleQueriesMixed(db *sql.DB) {
	// Correctly closed
	rowsOK, errOK := db.Query("SELECT * FROM table_ok")
	if errOK != nil {
		log.Fatal(errOK)
	}
	defer rowsOK.Close()
	for rowsOK.Next() {
		// ...
	}

	// Missing close
	rowsMissing, errMissing := db.Query("SELECT * FROM table_missing") // want "immediately defer close the rows returned here to avoid a memory leak"
	if errMissing != nil {
		log.Fatal(errMissing)
	}
	// defer rowsMissing.Close() is missing
	for rowsMissing.Next() {
		// ...
	}

	// Closed, but not deferred
	rowsNotDeferred, errNotDeferred := db.Query("SELECT * FROM table_not_deferred") // want "immediately defer close the rows returned here to avoid a memory leak"
	if errNotDeferred != nil {
		log.Fatal(errNotDeferred)
	}
	for rowsNotDeferred.Next() {
		// ...
	}
	rowsNotDeferred.Close() // Not deferred
}

func AnotherFunctionWithMixedCases(db *sql.DB) {
	// Missing close in another function
	data, err := db.Query("SELECT data FROM logs") // want "immediately defer close the rows returned here to avoid a memory leak"
	if err != nil {
		log.Fatal(err)
	}
	for data.Next() {
		// process
	}

	// Correctly closed in the same function
	userRows, userErr := db.Query("SELECT name, email FROM users WHERE active = true")
	if userErr != nil {
		log.Fatal(userErr)
	}
	defer userRows.Close()
	for userRows.Next() {
		// process
	}
}
