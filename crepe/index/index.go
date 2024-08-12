package index

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"go.uber.org/zap"
)

type Indexer struct{
}


func Start() Indexer {
	return Indexer{}
}

func getLanguage(ext string) (*sitter.Language, error) {
	switch ext {
	case ".go":
		return golang.GetLanguage(), nil
	default:
		return nil, errors.New("unknown file extension")
	}
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
	return nil
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
		err = ix.AddFile(f)
		if err != nil {
			// log.Error(err)
		}
	}
	return nil
}
