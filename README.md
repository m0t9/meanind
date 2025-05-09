# meanind

[![build](https://github.com/m0t9/meanind/actions/workflows/go.yml/badge.svg)](https://github.com/m0t9/meanind/actions/workflows/go.yml)
![coverage](https://raw.githubusercontent.com/m0t9/meanind/badges/.badges/master/coverage.svg)
**meanind** (meaningful indexing) — Go linter detecting confusing variable names for slice/array indexes in for-range loops.

## Details

One of the misleading pieces of Go language is for-range loops without value.
Newcomers from other languages may assume that the code below will print array items.
However, the indexes will be printed out.

```go
for item := range []string{"pen", "pineapple", "apple", "pen"} {
	fmt.Println(item)
}
```

In order to achieve the desired behavior, one small patch is needed
```go
for _, item := range []string{"pen", "pineapple", "apple", "pen"} {
	fmt.Println(item)
}
```

Linter from this repository is supposed to detect such confusing cases.
It searches for-range loops over arrays / slices where only key value is used with
iteratee name seems **meaningful** (that is not usual for Go language) and reports them.

The iteratee name is considered to be **meaningful** when **all** the case-insensitive rules **failing**
- iteratee identifier name is `_`, `i`, `j` or `j`
- iteratee is `ind` with some prefix / suffix
- iteratee looks like `idx`, `jdx` or `kdx` with some possible prefix / suffix 
- iteratee indeed is used in index expression of iterable (checked by AST traversal)

## Usage

### Installation

`go install github.com/m0t9/meanind/cmd/...@latest`

### Running

Linter is compatible with `go vet`

`go vet -vettool=$(which meanind) ./...` 

### Customization

Optionally one can rewrite regular expression for defining index identifiers

`go vet -vettool=$(which meanind) -index-regexp='custom-regex' ./...`
