// Package meanind provides analyzer for spotting possibly misused for-range loop indices.
package meanind

import (
	"flag"
	"fmt"
	"go/ast"
	"go/types"
	"regexp"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "meanind",
	Doc:      "reports possibly improperly used for-range loop indices",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// SetupFlags allows
func SetupFlags() {
	Analyzer.Flags.Init("meanind", flag.ExitOnError)
	Analyzer.Flags.StringVar(
		&indexRegexp,
		"index-regexp",
		indexRegexp,
		"regexp defining identifier names that can be used for indexes",
	)
}

// indexRegexp is (by default) a case-insensitive regular expression
// that matches a usual index identifier
// - common `i`, `j` or `k` name
// - looks like `idx`, `jdx`, `kdx` with some possible suffix or prefix
// - it is `ind` with some suffix or prefix
// - used as index of iterable (checked by AST traversal)
// (with -index-regexp=... it can be overridden)
var (
	indexRegexp = `(?i)^([ijk_]|[ijk]dx.*|.*[ijk]dx|.*ind|ind.*)$`
)

// sameObject is a helper function helping to compare two identifiers referring to the same variable.
func sameObject(pass *analysis.Pass, ident1, ident2 *ast.Ident) bool {
	obj1 := getObject(pass.TypesInfo, ident1)
	obj2 := getObject(pass.TypesInfo, ident2)
	return obj1 == obj2 && obj1 != nil
}

// getObject is a helper function to extract object from the pass information.
func getObject(info *types.Info, ident *ast.Ident) types.Object {
	if obj := info.Defs[ident]; obj != nil {
		return obj
	}
	return info.Uses[ident]
}

func run(pass *analysis.Pass) (any, error) {
	indexIdent := regexp.MustCompile(indexRegexp)

	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	filter := []ast.Node{(*ast.RangeStmt)(nil)}

	inspector.Preorder(filter, func(n ast.Node) {
		rng := n.(*ast.RangeStmt)

		// Check is not needed for cases with key & value both used like
		// `for idx, item := range iterable`.
		if rng.Value != nil {
			return
		}

		keyIdent, ok := rng.Key.(*ast.Ident)
		if rng.Key == nil || !ok {
			// For some reason, key is nil or not an identifier :|
			return
		}

		// The syntax may be ambiguous only with array or slice iterable type.
		it := pass.TypesInfo.TypeOf(rng.X).Underlying()
		_, isArray := it.(*types.Array)
		_, isSlice := it.(*types.Slice)
		if !(isSlice || isArray) {
			return
		}

		// From here it is a loop over slice / array in the form of `for id := range arr`.
		if indexIdent.Match([]byte(keyIdent.Name)) {
			return
		}

		// Check whether iteratee is used as index of iterable (only for identifiers case).
		usedAsIndex := false
		if iterableIdent, ok := rng.X.(*ast.Ident); ok {
			ast.Inspect(rng.Body, func(n ast.Node) bool {
				idxExpr, ok := n.(*ast.IndexExpr)
				if !ok {
					return true
				}
				accessedIdent, ok := idxExpr.X.(*ast.Ident)
				if !ok {
					return true
				}

				if !sameObject(pass, accessedIdent, iterableIdent) {
					return true
				}

				idxUsed := false
				ast.Inspect(idxExpr.Index, func(n ast.Node) bool {
					ident, ok := n.(*ast.Ident)
					if !ok {
						return true
					}
					if sameObject(pass, keyIdent, ident) {
						idxUsed = true
						return false
					}
					return true
				})

				if idxUsed {
					usedAsIndex = true
					return false
				}
				return true
			})
		}

		if usedAsIndex {
			return
		}

		pass.Report(analysis.Diagnostic{
			Message: fmt.Sprintf(
				`for-range loop with confusing iteratee name detected. Is %q name suitable for array/slice index?`,
				keyIdent.Name),
			Pos: rng.Pos(),
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: fmt.Sprintf(`consider rename %q to more suitable index name`, keyIdent.Name),
				},
			},
		})
	})
	return nil, nil
}
