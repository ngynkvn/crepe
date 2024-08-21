package index

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/ngynkvn/crepe/sql/gen/cindex"
)

// TODO: make this a proper interface
type DB struct {
	conn   *pgx.Conn
	cindex *cindex.Queries
}

func DefaultDSN() string {
	return "postgres://postgres:postgres@localhost/postgres?sslmode=disable"
}

func NewDB() (*DB, error) {
	conn, err := pgx.Connect(context.TODO(), DefaultDSN())
	if err != nil {
		return nil, err
	}
	cindex := cindex.New(conn)
	return &DB{
		conn:   conn,
		cindex: cindex,
	}, nil
}

// TODO: This should ideally accept a different struct as the second parameter which we can then map to the DB Insert
// AddFile(ctx, FileInfo{...})
func (db *DB) AddFile(ctx context.Context, params cindex.AddFileParams) error {
	_, err := db.cindex.AddFile(ctx, params)
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
	// db.Select(&tokens, "SELECT * FROM tokens")

	writeJSON(w, tokens)
}

func (db *DB) Mount(srv *http.ServeMux) *http.ServeMux {
	srv.HandleFunc("GET /tokens", db.tokens)
	return srv
}

// TODO: make this configurable on what indexer selects to parse from tree-sitter grammar.
// Scrape node types from https://github.com/tree-sitter/tree-sitter-go/blob/master/src/node-types.json
