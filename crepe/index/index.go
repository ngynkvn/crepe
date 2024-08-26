package index

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ngynkvn/crepe/sql/gen/cindex"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"go.uber.org/zap"
)

type Indexer struct {
	db *DB
}

func New() (*Indexer, error) {
	db, err := NewDB()
	if err != nil {
		return nil, err
	}
	return &Indexer{
		db,
	}, nil
}

// func (ix Indexer) AddFile(f *object.File) error {
func (ix Indexer) AddFile(repo string, fp string) error {
	ctx := context.TODO()
	log := zap.S().With("pkg", "crepe/index")
	ext := filepath.Ext(fp)

	parser, err := getTreesitterParser(ext)
	if err != nil {
		return err
	}
	f, err := os.Open(fp)
	if err != nil {
		return err
	}

	contents, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	tree, err := parser.ParseCtx(ctx, nil, contents)
	if err != nil {
		return err
	}
	_, err = ix.db.cindex.AddFile(
		ctx,
		cindex.AddFileParams{
			Repo:                repo,
			FilePath:            fp,
			FileName:            fp,
			ProgrammingLanguage: ext,
			Contents:            string(contents),
		})
	if err != nil {
		return err
	}

	// Walk the tree and add all nodes that are of a type that we want to index
	walk(tree.RootNode(), (func(n *sitter.Node) {
		// TODO: slice should be determined by extension
		if !slices.Contains(allowedGoNodeTypes, n.Type()) {
			return
		}
		start := n.StartPoint()
		end := n.EndPoint()
		_, err := ix.db.cindex.AddCodeElement(ctx, cindex.AddCodeElementParams{
			FileName:    fp,
			ElementType: n.Type(),
			Contents:    n.Content(contents),
			StartLine:   int32(start.Row),
			EndLine:     int32(end.Row),
		})
		if err != nil {
			log.Error(err)
		}
	}))
	return nil
}

func (ix Indexer) AddRepo(repoPath string) error {
	ctx := context.TODO()
	log := zap.S()
	// TODO: add support for link to github repos
	// TODO: and find a more appropriate place for this.
	ix.db.cindex.AddRepo(ctx, cindex.AddRepoParams{
		Repo:     repoPath,
		RepoType: "git",
	})

	// Iterate over all files and add them to the index
	err := filepath.WalkDir(repoPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the .git directory
		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}

		// Process only regular files
		if !d.IsDir() {
			relPath, err := filepath.Rel(repoPath, path)
			if err != nil {
				log.Errorf("failed to get relative path: %w", err)
				return err
			}

			log.Infof("Indexing file: %s", relPath)
			ix.AddFile(repoPath, relPath)
		}
		return nil
	})
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (ix Indexer) Serve() error {
	log := zap.S()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

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

var allowedGoNodeTypes = []string{
	"func_literal",
	"identifier",
	"interpreted_string_literal",
	"import_spec",
	"package_clause",
	"type_identifier",
}

func getTreesitterParser(ext string) (*sitter.Parser, error) {
	language, err := getLanguage(ext)
	if err != nil {
		return nil, err
	}
	parser := sitter.NewParser()
	parser.SetLanguage(language)
	return parser, nil
}

func getLanguage(ext string) (*sitter.Language, error) {
	switch ext {
	case ".go":
		return golang.GetLanguage(), nil
	default:
		return nil, errors.New("unknown file extension")
	}
}

func walk(n *sitter.Node, f func(n *sitter.Node)) {
	f(n)
	for i := uint32(0); i < n.NamedChildCount(); i++ {
		walk(n.NamedChild(int(i)), f)
	}
}
