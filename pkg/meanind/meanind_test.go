package meanind_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/m0t9/meanind/pkg/meanind"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestMeanindAnalyzer(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	testdata := filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata")
	meanind.SetupFlags()
	analysistest.Run(t, testdata, meanind.Analyzer, "meanind")
}
