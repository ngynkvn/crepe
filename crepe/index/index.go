package index

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"github.com/ngynkvn/crepe/sql/gen/cindex"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
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

func (ix Indexer) AddFile(repoUrl string, fp string) error {
	ctx := context.TODO()
	log := slog.Default().With("pkg", "crepe/index")
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
			Url:                 repoUrl,
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
			log.Error(err.Error())
		}
	}))
	return nil
}

func (ix Indexer) Clear() {
	_, err := ix.db.conn.Exec(context.Background(), "TRUNCATE cindex.code_repositories CASCADE;")
	if err != nil {
		panic(err)
	}
}

func (ix Indexer) AddRepo(name, repoUrl string) error {
	ctx := context.TODO()
	log := slog.Default()
	// TODO: add support for link to github repos
	// TODO: and find a more appropriate place for this.
	ix.db.cindex.AddRepo(ctx, cindex.AddRepoParams{
		RepoName: name,
		Url:      repoUrl,
		RepoType: "git",
	})

	// Iterate over all files and add them to the index
	err := filepath.WalkDir(repoUrl, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the .git directory
		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}

		// Process only regular files
		if !d.IsDir() {
			relPath, err := filepath.Rel(repoUrl, path)
			if err != nil {
				log.Error("failed to get relative path", "error", err)
				return err
			}

			log.Info("Indexing file", "file", relPath)
			err = ix.AddFile(repoUrl, relPath)
			if err != nil && !errors.Is(err, errUnknownFileExtension) {
				log.Error("Failed to index", "error", err)
			}
		}
		return nil
	})
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
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

var errUnknownFileExtension = errors.New("unknown file extension")

func getLanguage(ext string) (*sitter.Language, error) {
	switch ext {
	case ".go":
		return golang.GetLanguage(), nil
	default:
		return nil, errUnknownFileExtension
	}
}

func walk(n *sitter.Node, f func(n *sitter.Node)) {
	f(n)
	for i := uint32(0); i < n.NamedChildCount(); i++ {
		walk(n.NamedChild(int(i)), f)
	}
}
