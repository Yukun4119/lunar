package util

import (
	"github.com/ryqdev/golang_utils/log"
	"go/ast"
	"gopkg.in/yaml.v3"
	"lunar_uml/models"
	"os"
)

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

func LoadConfig(configFile string) models.YamlConfig {
	file, err := os.ReadFile(configFile)
	if err != nil {
		log.Error("Read file error")
		os.Exit(1)
	}
	config := models.YamlConfig{}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Error("Unmarshal error")
		os.Exit(1)
	}
	return config
}
