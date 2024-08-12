package index

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/marcboeker/go-duckdb"
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
		CREATE TABLE tokens (
			filename TEXT,
			type TEXT,
			language TEXT,
			contents TEXT 
		);
	`)
}

func (db *DB) AddFile(filename string, nodetype string, ext string, contents string) error {
	_, err := db.Exec(`
			INSERT INTO tokens (filename, type, language, contents) 
			VALUES (?, ?, ?, ?)`,
		filename,
		nodetype,
		ext,
		contents,
	)
	return err
}

func writeJSON(w http.ResponseWriter, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
	return nil
}

func (db *DB) tokens(w http.ResponseWriter, r *http.Request) {
	type TokenJSON struct {
		Filename string `json:"filename"`
		Type     string `json:"type"`
		Language string `json:"language"`
		Contents string `json:"contents"`
	}
	var tokens []TokenJSON
	db.Select(&tokens, "SELECT * FROM tokens")

	writeJSON(w, tokens)
}

func (db *DB) Mount(srv *http.ServeMux) *http.ServeMux {
	srv.HandleFunc("GET /tokens", db.tokens)
	return srv
}

// TODO: make this configurable on what indexer selects to parse from tree-sitter grammar.
// Scrape node types from https://github.com/tree-sitter/tree-sitter-go/blob/master/src/node-types.json
