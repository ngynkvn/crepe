package index

import (
	"context"
	"errors"
	"path/filepath"
	"slices"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"go.uber.org/zap"
)

type Indexer struct {
	db *DB
}

func Start() (*Indexer, error) {
	db, err := NewDB()
	if err != nil {
		return nil, err
	}
	return &Indexer{
		db,
	}, nil
}

func getLanguage(ext string) (*sitter.Language, error) {
	switch ext {
	case ".go":
		return golang.GetLanguage(), nil
	default:
		return nil, errors.New("unknown file extension")
	}
}

var allowedGoTypes = []string{
	"func_literal",
	"identifier",
	"interpreted_string_literal",
	"import_spec",
	"package_clause",
	"type_identifier",
}

func (ix Indexer) AddFile(f *object.File) error {
	log := zap.S().With("pkg", "crepe/index")
	ext := filepath.Ext(f.Name)

	language, err := getLanguage(ext)
	if err != nil {
		return err
	}
	parser := sitter.NewParser()
	parser.SetLanguage(language)
	contents, err := f.Contents()
	if err != nil {
		return err
	}

	tree, err := parser.ParseCtx(context.TODO(), nil, []byte(contents))
	if err != nil {
		return err
	}
	log.Debug(contents)
	log.Debug(tree.RootNode())
	// Walk the tree and add all nodes that are of a type that we want to index
	walk(tree.RootNode(), (func(n *sitter.Node) {
		if slices.Contains(allowedGoTypes, n.Type()) {
			ix.db.MustExec(`
			INSERT INTO tokens (id, filename, type, language, contents) 
			VALUES (nextval('tokens_id_seq'),?, ?, ?, ?)`,
				f.Name,
				n.Type(),
				ext,
				n.Content([]byte(contents)),
			)
		}
	}))
	return nil
}

func walk(n *sitter.Node, f func(n *sitter.Node)) {
	f(n)
	for i := uint32(0); i < n.NamedChildCount(); i++ {
		walk(n.NamedChild(int(i)), f)
	}
}

func (ix Indexer) AddRepo(path string) error {
	log := zap.S().With("pkg", "crepe/index")
	log.With("path", path).Info("checking path")
	// Check that path is valid git repository
	// It is either a url to a git repo, or it is a local path.
	repo, err := git.PlainOpen(path)
	if err != nil {
		log.Error(err)
		return err
	}

	treeObjects, err := repo.TreeObjects()
	if err != nil {
		log.Error(err)
		return err
	}
	// Iterate over all files in the repository
	// and add them to the slice
	var files []*object.File
	err = treeObjects.ForEach(func(t *object.Tree) error {
		return t.Files().ForEach(func(f *object.File) error {
			log.With("file", f).Info("")
			files = append(files, f)
			return nil
		})
	})
	if err != nil {
		log.Error(err)
		return err
	}

	// Iterate over all files and add them to the index
	for _, f := range files {
		log.With("path", path, "file", f.Name).Info("")
		ix.AddFile(f)
	}
	return nil
}
