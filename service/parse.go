package service

import (
	"github.com/ryqdev/golang_utils/log"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func (l *LunarUML) ParseCodeWithAST() *LunarUML {
	if l.Err != nil {
		return l
	}
	l.FSet = token.NewFileSet()
	l.Node, _ = parser.ParseFile(l.FSet, l.Config.LunarConfig.FilePath, nil, parser.ParseComments)

	log.Info("Finish parsing code")
	return l
}

func (l *LunarUML) PrintNodeForDebug() *LunarUML {
	if l.Err != nil {
		return l
	}
	if l.Config.LunarConfig.IsDebug {
		log.Info("Debug mode on")
		// Create a new file
		f, err := os.Create("output/ast.txt")
		if err != nil {
			log.Error("error: %+v", err)
			os.Exit(1)
		}
		defer f.Close()

		// Redirect the output of ast.Print to the file
		err = ast.Fprint(f, l.FSet, l.Node, nil)
		if err != nil {
			log.Error("error: %+v", err)
			os.Exit(1)
		}
	} else {
		log.Info("Debug mode off")
	}
	return l
}
