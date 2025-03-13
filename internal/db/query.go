package db

import (
	"database/sql"
	"sync"
)

var queriesOnce sync.Once
var queries *Queries

func NewQueries(db *sql.DB) *Queries {
	queriesOnce.Do(func() {
		queries = New(db)
	})

	return queries
}
