package index

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/marcboeker/go-duckdb"
	_ "github.com/marcboeker/go-duckdb"
)

// TODO: make this a proper interface
type DB struct {
	*sqlx.DB
}

func NewDB() (*DB, error) {
	connector, err := duckdb.NewConnector("crepe.db", nil)
	sdb := sql.OpenDB(connector)
	db := sqlx.NewDb(sdb, "duckdb")
	if err != nil {
		return nil, err
	}
	initDB(db)
	return &DB{
		db,
	}, nil
}

func initDB(conn *sqlx.DB) {
	// TODO: make this a proper migration system
	conn.MustExec(`
		DROP TABLE IF EXISTS tokens;
		DROP SEQUENCE IF EXISTS tokens_id_seq;
		CREATE TABLE tokens (
			id INTEGER PRIMARY KEY,
			filename TEXT,
			type TEXT,
			language TEXT,
			contents TEXT 
		);

		CREATE SEQUENCE tokens_id_seq;
	`)
}

// TODO: make this configurable on what indexer selects to parse from tree-sitter grammar.
// Scrape node types from https://github.com/tree-sitter/tree-sitter-go/blob/master/src/node-types.json
