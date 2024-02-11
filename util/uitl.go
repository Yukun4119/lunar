package util

import "go/ast"

func IsRPC(expr *ast.CallExpr) (rpc string, isRpc bool) {
	rpc = ""
	isRpc = false

	// expr is the expression you parsed
	ast.Inspect(expr, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			if selExpr, ok := x.Fun.(*ast.SelectorExpr); ok {
				if callExpr, ok := selExpr.X.(*ast.CallExpr); ok {
					if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
						if ident, ok := selExpr.X.(*ast.Ident); ok {
							if ident.Name == "rpc" {
								rpc = selExpr.Sel.Name
								isRpc = true
								return false // return false if you want to stop traversing after found
							}
						}
					}
				}
			}
		}
		return true // return true to keep traversing the tree
	})
	return
}

// Helper function to get the name from a CallExpr
func GetCallExprName(expr *ast.CallExpr) string {
	switch expr := expr.Fun.(type) {
	case *ast.SelectorExpr:
		return expr.Sel.Name
	case *ast.Ident:
		return expr.Name
	default:
		return "unknown"
	}
}
