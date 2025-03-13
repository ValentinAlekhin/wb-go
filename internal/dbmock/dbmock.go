package dbmock

import (
	"database/sql"
	"github.com/ValentinAlekhin/wb-go/internal/db"
	_ "modernc.org/sqlite"
)

func NewDBMock() *sql.DB {
	instance, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		panic(err)
	}

	err = db.MigrateOnce(instance)
	if err != nil {
		panic(err)
	}

	return instance
}
