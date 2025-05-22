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

// rowsInfo stores information about a sql.Rows variable.
// The funcScope field was removed as scope is handled by iterating FuncDecls.
type rowsInfo struct {
	pos    token.Pos // position of the Query/QueryContext call
	closed bool      // whether a 'defer rows.Close()' was found
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			// Only inspect function declarations.
			fd, ok := n.(*ast.FuncDecl)
			if !ok {
				return true // Continue to find other function declarations
			}

			// For each function, create a map for its rows variables.
			// Key is the variable name (e.g., "rows").
			funcRows := make(map[string]rowsInfo)

			// Inspect the body of the function.
			ast.Inspect(fd.Body, func(node ast.Node) bool {
				// Find "rows := db.Query(...)" or "rows, err := db.Query(...)"
				if assignStmt, ok := node.(*ast.AssignStmt); ok {
					if len(assignStmt.Lhs) > 0 && len(assignStmt.Rhs) > 0 {
						// Check if the RHS is a call expression like db.Query() or db.QueryContext()
						if callExpr, ok := assignStmt.Rhs[0].(*ast.CallExpr); ok {
							if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
								if selectorExpr.Sel.Name == "Query" || selectorExpr.Sel.Name == "QueryContext" {
									// Get the variable name assigned to sql.Rows.
									// We assume the first LHS operand is the rows variable.
									if ident, ok := assignStmt.Lhs[0].(*ast.Ident); ok {
										if ident.Name != "_" { // Ignore blank identifier
											// Store information about this rows variable.
											// Mark as not closed initially.
											funcRows[ident.Name] = rowsInfo{pos: selectorExpr.Pos(), closed: false}
										}
									}
								}
							}
						}
					}
				}

				// Find "defer rows.Close()"
				if deferStmt, ok := node.(*ast.DeferStmt); ok {
					// deferStmt.Call is already an *ast.CallExpr, no need for type assertion.
					if callExpr := deferStmt.Call; callExpr != nil {
						if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok { // Check if it's a selector like rows.Close
							if selectorExpr.Sel.Name == "Close" {
								// Get the variable name from "variable.Close()"
								if ident, ok := selectorExpr.X.(*ast.Ident); ok {
									// If this variable is one we are tracking, mark it as closed.
									if info, exists := funcRows[ident.Name]; exists {
										info.closed = true
										funcRows[ident.Name] = info
									}
								}
							}
						}
					}
				}
				return true // Continue inspecting other nodes in the function body
			})

			// After inspecting the entire function body, report any unclosed rows.
			for _, info := range funcRows { // varName is not needed here, only info
				if !info.closed {
					pass.Reportf(info.pos, "immediately defer close the rows returned here to avoid a memory leak")
				}
			}
			return true // Continue to find other function declarations in the file
		})
	}
	return nil, nil
}
