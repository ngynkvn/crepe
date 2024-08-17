package main

import (
	"github.com/ngynkvn/crepe/crepe/index"
	"go.uber.org/zap"
)

var atom = zap.NewAtomicLevel()

func main() {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = atom
	zap.ReplaceGlobals(zap.Must(cfg.Build()))
	atom.SetLevel(zap.DebugLevel)

	ix, err := index.New()
	if err != nil {
		panic(err)
	}
	ix.AddRepo(".")
	ix.Serve()
}
