package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
)

type YamlConfig struct {
	LunarConfig LunarConfig `yaml:"lunarConfig"`
}

type LunarConfig struct {
	CurService string `yaml:"curService"`
	TargetInf  string `yaml:"targetInf"`
	FilePath   string `yaml:"filePath"`
}

var (
	Config       YamlConfig
	plantUML     []string
	participants []string
)

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	err := loadConfig()
	if err != nil {
		log.Println("load config error")
		os.Exit(1)
	}

	// Step 2
	// parse code
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, Config.LunarConfig.FilePath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	//ast.Print(fset, node)

	// Step 3
	// Get all the participants
	participants = append(participants, Config.LunarConfig.TargetInf)

	// step 4
	// Start PlantUML sequence diagram
	plantUML = append(plantUML, "@startuml")
	plantUML = append(plantUML, "autonumber")

	ast.Inspect(node, func(n ast.Node) bool {
		switch fn := n.(type) {
		case *ast.FuncDecl:
			if fn.Name.Name == Config.LunarConfig.TargetInf {
				TranverseFunc(fn)
			}
		}
		return true
	})

	// End PlantUML sequence diagram
	plantUML = append(plantUML, "@enduml")

	// Print PlantUML
	for _, line := range plantUML {
		println(line)
	}
}

func loadConfig() error {
	file, err := os.ReadFile("./config.yml")
	if err != nil {
		log.Println("read file error")
		return err
	}
	return yaml.Unmarshal(file, &Config)
}

// Helper function to get the name from a CallExpr
func getCallExprName(expr *ast.CallExpr) string {
	switch expr := expr.Fun.(type) {
	case *ast.SelectorExpr:
		return expr.Sel.Name
	case *ast.Ident:
		return expr.Name
	default:
		return "unknown"
	}
}

func isRPC(expr *ast.CallExpr) (rpc string, isRpc bool) {
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

func TranverseFunc(fn *ast.FuncDecl) {
	for _, node := range fn.Body.List {
		switch nodeType := node.(type) {
		case *ast.AssignStmt:
			AnalyseAssignStmt(nodeType)
		case *ast.SwitchStmt:
			AnalyseSwitchStmt(nodeType)
		case *ast.IfStmt:
			AnalyseIfStmt(nodeType)
		case *ast.RangeStmt:
			AnalyseRangeStmt(nodeType)
		}
	}
}

func Tranverse(fn ast.Node) {
	ast.Inspect(fn, func(n ast.Node) bool {
		switch fn := n.(type) {
		case *ast.SwitchStmt:
			AnalyseSwitchStmt(fn)
			return false
		case *ast.IfStmt:
			AnalyseIfStmt(fn)
			return false
		case *ast.RangeStmt:
			AnalyseRangeStmt(fn)
			return false
		case *ast.CallExpr:
			AnalyseCallExpr(fn)
			return false
		}
		return true
	})
}

func AnalyseCallExpr(expr *ast.CallExpr) {
	if rpc, ok := isRPC(expr); ok {
		participants = append(participants, rpc)
		plantUML = append(plantUML, Config.LunarConfig.CurService+" -> "+rpc+": call "+getCallExprName(expr))
	} else {
		plantUML = append(plantUML, Config.LunarConfig.CurService+" -> "+Config.LunarConfig.CurService+": call "+getCallExprName(expr))
	}
}

func AnalyseRangeStmt(r *ast.RangeStmt) {
	moduleName := r.Key.(*ast.Ident).Obj.Decl.(*ast.AssignStmt).Rhs[0].(*ast.UnaryExpr).X.(*ast.SelectorExpr).Sel.Name
	plantUML = append(plantUML, "loop "+moduleName)
	Tranverse(r.Body)
	plantUML = append(plantUML, "end")
}

func AnalyseIfStmt(ifStmt *ast.IfStmt) {
	Cond := AnalyseExprIfStmt(ifStmt.Cond)
	plantUML = append(plantUML, "alt if "+Cond)
	Tranverse(ifStmt.Body)
	plantUML = append(plantUML, "end")
}

func AnalyseSwitchStmt(SwitchStmt *ast.SwitchStmt) {
	plantUML = append(plantUML, "alt switch stat")
	Tranverse(SwitchStmt.Body)
	plantUML = append(plantUML, "end")

}

func AnalyseAssignStmt(AssignStmt *ast.AssignStmt) {
	Tranverse(AssignStmt)
}

func AnalyseExprIfStmt(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.BinaryExpr:
		return fmt.Sprintf("%s %s %s", AnalyseExprIfStmt(e.X), e.Op, AnalyseExprIfStmt(e.Y))
	case *ast.CallExpr:
		fun := AnalyseExprIfStmt(e.Fun)
		var args []string
		for _, a := range e.Args {
			args = append(args, AnalyseExprIfStmt(a))
		}
		return fmt.Sprintf("%s(%s)", fun, strings.Join(args, ", "))
	case *ast.UnaryExpr:
		return fmt.Sprintf("%s%s", e.Op, AnalyseExprIfStmt(e.X))
	default:
		return fmt.Sprintf("%T", e)
	}
}
