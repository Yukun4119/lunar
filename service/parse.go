package service

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func (l *LunarUML) ParseCodeWithAST() {
	l.FSet = token.NewFileSet()
	l.Node, _ = parser.ParseFile(l.FSet, l.Config.LunarConfig.FilePath, nil, parser.ParseComments)

	// TODO: handle error
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func (l *LunarUML) PrintNodeForDebug() {
	if l.Config.LunarConfig.IsDebug {
		err := ast.Print(l.FSet, l.Node)
		if err != nil {
			log.Fatal(err)
		}
	}
}
