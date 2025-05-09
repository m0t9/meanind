package main

import (
	"github.com/m0t9/meanind/pkg/meanind"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	meanind.SetupFlags()
	singlechecker.Main(meanind.Analyzer)
}
