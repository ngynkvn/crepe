package main

import (
	"os"

	"github.com/ngynkvn/crepe/crepe/index"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"go.uber.org/zap"
)

var atom = zap.NewAtomicLevel()

func main() {
	println(os.Getwd())
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = atom
	zap.ReplaceGlobals(zap.Must(cfg.Build()))
	atom.SetLevel(zap.DebugLevel)

	language := golang.GetLanguage()
	parser := sitter.NewParser()
	parser.SetLanguage(language)

	ix, err := index.Start()
	if err != nil {
		panic(err)
	}
	ix.AddRepo(".")
}
