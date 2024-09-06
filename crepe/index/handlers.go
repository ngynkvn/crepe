package index

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/ngynkvn/crepe/sql/gen/cindex"
	"github.com/ngynkvn/crepe/util/collections"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (ix Indexer) Serve() error {
	log := slog.Default()

	logger := httplog.NewLogger("cindex", httplog.Options{
		// JSON:             true,
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		// TimeFieldFormat: time.RFC850,
		Tags: map[string]string{
			"version": "v1.0-81aa4244d9fc8076a",
			"env":     "dev",
		},
		QuietDownRoutes: []string{
			"/",
			"/metrics",
		},
		QuietDownPeriod: time.Minute,
		// SourceFieldName: "source",
	})

	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(logger))

	r.Route("/repositories", func(r chi.Router) {
		r.Get("/", http.HandlerFunc(ix.HandleGetRepositories))
		r.Post("/", http.HandlerFunc(ix.HandlePostRepositories))
		r.Get("/{repoId}", http.HandlerFunc(ix.HandleGetRepositoriesRepoId))
		r.Delete("/{repoId}", http.HandlerFunc(ix.HandleDeleteRepositoriesRepoId))
	})

	r.Get("/search", http.HandlerFunc(ix.HandleSearch))
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	log.Info("starting server")
	return http.ListenAndServe("0.0.0.0:8080", r)
}

func (ix Indexer) HandleGetRepositories(response http.ResponseWriter, req *http.Request) {
	type RepositoryListItem struct {
		cindex.CindexCodeRepository
		FileCount int `json:"file_count"`
	}
	ctx := req.Context()
	repos, err := ix.db.cindex.GetRepositories(ctx)
	if err != nil {
		err := encode(response, req, http.StatusInternalServerError, err.Error())
		if err != nil {
			slog.Default().Error(err.Error())
		}
	}
	encode(response, req, http.StatusOK, collections.Map(repos, func(t cindex.GetRepositoriesRow) RepositoryListItem {
		return RepositoryListItem{
			t.CindexCodeRepository,
			int(t.NumFiles),
		}
	}))
}
func (ix Indexer) HandlePostRepositories(response http.ResponseWriter, req *http.Request) {}

func (ix Indexer) HandleGetRepositoriesRepoId(response http.ResponseWriter, req *http.Request)    {}
func (ix Indexer) HandleDeleteRepositoriesRepoId(response http.ResponseWriter, req *http.Request) {}

func (ix Indexer) HandleSearch(response http.ResponseWriter, req *http.Request) {}
