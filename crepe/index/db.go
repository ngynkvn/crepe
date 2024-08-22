package index

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/ngynkvn/crepe/sql/gen/cindex"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
	cfg, err := pgx.ParseConfig(DefaultDSN())
	if err != nil {
		return nil, err
	}
	cfg.Tracer = newTracer()

	conn, err := pgx.ConnectConfig(context.TODO(), cfg)
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

type tracer struct {
	numQueries prometheus.Counter
	queryTime  prometheus.Summary
}

func newTracer() *tracer {
	return &tracer{
		numQueries: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "cindex",
			Subsystem: "query",
			Name:      "count",
			Help:      "",
		}),
		queryTime: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: "cindex",
			Subsystem: "query",
			Name:      "time",
			Help:      "",
		}),
	}
}

type startQueryTimeKey struct{}

func (t *tracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	t.numQueries.Inc()
	ctx = context.WithValue(ctx, startQueryTimeKey{}, time.Now())
	return ctx
}

func (t *tracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	if val, ok := ctx.Value(startQueryTimeKey{}).(time.Time); ok {
		t.queryTime.Observe(time.Since(val).Seconds())
	}
}
