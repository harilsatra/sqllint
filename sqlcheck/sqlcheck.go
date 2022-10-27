package sqlcheck

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "sqllint",
	Doc:  "reports whether SQL rows are closed or not",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {

		selectQueryExists := false
		var selectQueryPos token.Pos
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
			}

			if selectQueryExists {
				count_statement++
			}

			if count_statement > 1 && fun.Sel.Name != "Close" {
				pass.Reportf(selectQueryPos, "defer close the rows returned here, immediately to avoid a memory leak ")
			}

			return true
		})
	}
	return nil, nil
}
