package main

import (
	"os"

	"github.com/ngynkvn/crepe/crepe/index"
	"go.uber.org/zap"
)

var atom = zap.NewAtomicLevel()

func main() {
	println(os.Getwd())
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = atom
	zap.ReplaceGlobals(zap.Must(cfg.Build()))
	atom.SetLevel(zap.DebugLevel)

	ix, err := index.Start()
	if err != nil {
		panic(err)
	}
	ix.AddRepo(".")
}
