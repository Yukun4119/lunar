package service

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"lunar_uml/models"
	"lunar_uml/util"
	"strings"
)

type LunarUML struct {
	Config       models.YamlConfig
	PlantUML     []string
	Participants []string
	FSet         *token.FileSet
	Node         *ast.File
}

// TODO: Do I need interface?
//type Analysis interface {
//	AnalyseCallExpr(expr *ast.CallExpr)
//}

func (l *LunarUML) Inspect() {
	ast.Inspect(l.Node, func(n ast.Node) bool {
		switch fn := n.(type) {
		case *ast.FuncDecl:
			if fn.Name.Name == l.Config.LunarConfig.TargetInf {
				log.Println("Found the target")
				l.TranverseFunc(fn)
			}
		}
		return true
	})
	log.Println("Finish Inspecting")
}

func (l *LunarUML) TranverseFunc(fn *ast.FuncDecl) {
	for _, node := range fn.Body.List {
		switch nodeType := node.(type) {
		case *ast.AssignStmt:
			l.AnalyseAssignStmt(nodeType)
		case *ast.SwitchStmt:
			l.AnalyseSwitchStmt(nodeType)
		case *ast.IfStmt:
			l.AnalyseIfStmt(nodeType)
		case *ast.RangeStmt:
			l.AnalyseRangeStmt(nodeType)
		case *ast.ExprStmt:
			// TODO: do something
			l.AnalyseCallExpr(nodeType.X)
		}
	}
}

func (l *LunarUML) Tranverse(fn ast.Node) {
	ast.Inspect(fn, func(n ast.Node) bool {
		switch fn := n.(type) {
		case *ast.SwitchStmt:
			l.AnalyseSwitchStmt(fn)
			return false
		case *ast.IfStmt:
			l.AnalyseIfStmt(fn)
			return false
		case *ast.RangeStmt:
			l.AnalyseRangeStmt(fn)
			return false
		case *ast.CallExpr:
			l.AnalyseCallExpr(fn)
			return false
		}
		return true
	})
}

func (l *LunarUML) AnalyseCallExpr(expr *ast.CallExpr) {
	if rpc, ok := util.IsRPC(expr); ok {
		l.Participants = append(l.Participants, rpc)
		l.PlantUML = append(l.PlantUML, l.Config.LunarConfig.CurService+" -> "+rpc+": call "+util.GetCallExprName(expr))
	} else {
		l.PlantUML = append(l.PlantUML, l.Config.LunarConfig.CurService+" -> "+l.Config.LunarConfig.CurService+": call "+util.GetCallExprName(expr))
	}
}

func (l *LunarUML) AnalyseRangeStmt(r *ast.RangeStmt) {
	moduleName := r.Key.(*ast.Ident).Obj.Decl.(*ast.AssignStmt).Rhs[0].(*ast.UnaryExpr).X.(*ast.SelectorExpr).Sel.Name
	l.PlantUML = append(l.PlantUML, "loop "+moduleName)
	l.Tranverse(r.Body)
	l.PlantUML = append(l.PlantUML, "end")
}

func (l *LunarUML) AnalyseIfStmt(ifStmt *ast.IfStmt) {
	Cond := l.AnalyseExprIfStmt(ifStmt.Cond)
	l.PlantUML = append(l.PlantUML, "alt if "+Cond)
	l.Tranverse(ifStmt.Body)
	l.PlantUML = append(l.PlantUML, "end")
}

func (l *LunarUML) AnalyseSwitchStmt(SwitchStmt *ast.SwitchStmt) {
	l.PlantUML = append(l.PlantUML, "alt switch stat")
	l.Tranverse(SwitchStmt.Body)
	l.PlantUML = append(l.PlantUML, "end")

}

func (l *LunarUML) AnalyseAssignStmt(AssignStmt *ast.AssignStmt) {
	l.Tranverse(AssignStmt)
}

func (l *LunarUML) AnalyseExprIfStmt(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.BinaryExpr:
		return fmt.Sprintf("%s %s %s", l.AnalyseExprIfStmt(e.X), e.Op, l.AnalyseExprIfStmt(e.Y))
	case *ast.CallExpr:
		fun := l.AnalyseExprIfStmt(e.Fun)
		var args []string
		for _, a := range e.Args {
			args = append(args, l.AnalyseExprIfStmt(a))
		}
		return fmt.Sprintf("%s(%s)", fun, strings.Join(args, ", "))
	case *ast.UnaryExpr:
		return fmt.Sprintf("%s%s", e.Op, l.AnalyseExprIfStmt(e.X))
	default:
		return fmt.Sprintf("%T", e)
	}
}
