## sqllint

Checks if rows returned by db.Query() or db.QueryContext() are not immediately defer closed.

```
go get -u github.com/harilsatra/sqllint
```

## Example

Given the file

```
package main

import (
	"database/sql"
	"log"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
)

func runQuery(db *sql.DB) error {
	q := sq.Select("id").From("some_table")

	rows, err := q.RunWith(db).Query()
	if err != nil {
		log.Fatal(err)
		return err
	}

	ids := []uint64{}
	for rows.Next() {
		var id uint64
		rows.Scan(&id)

		ids = append(ids, id)
	}

	return nil
}
```

running sqllint on it with default options

```
sqllint file.go
```

will produce

```
file.go:14:15: immediately defer close the rows returned here to avoid a memory leak
```