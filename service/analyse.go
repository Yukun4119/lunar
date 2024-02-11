package service

import (
	"fmt"
	"go/ast"
	"lunar_uml/util"
	"strings"
)

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
	if rpc, ok := util.IsRPC(expr); ok {
		participants = append(participants, rpc)
		plantUML = append(plantUML, Config.LunarConfig.CurService+" -> "+rpc+": call "+util.GetCallExprName(expr))
	} else {
		plantUML = append(plantUML, Config.LunarConfig.CurService+" -> "+Config.LunarConfig.CurService+": call "+util.GetCallExprName(expr))
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
