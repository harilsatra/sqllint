package sqlcheck

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "sqllint",
	Doc:  "reports whether SQL rows are closed or not",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	fmt.Println(len(pass.Files))
	for _, file := range pass.Files {

		selectQueryExists := false
		var selectQueryPos token.Pos
		var selectQueryExpr *ast.CallExpr
		count_statement := 0
		ast.Inspect(file, func(n ast.Node) bool {

			ce, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			fun, ok := ce.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			if fun.Sel.Name == "Query" || fun.Sel.Name == "QueryContext" {
				selectQueryExists = true
				selectQueryPos = fun.Pos()
				selectQueryExpr = ce

			}

			if selectQueryExists {
				count_statement++
			}

			if count_statement > 1 && fun.Sel.Name != "Close" {
				pass.Reportf(selectQueryPos, "defer close the rows returned here, immediately to avoid a memory leak ",
					render(pass.Fset, selectQueryExpr))
			}

			return true
		})
	}
	return nil, nil
}

// render returns the pretty-print of the given node
func render(fset *token.FileSet, x interface{}) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, x); err != nil {
		panic(err)
	}
	return buf.String()
}
