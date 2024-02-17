package service

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

func (l *LunarUML) ParseCodeWithAST() {
	l.FSet = token.NewFileSet()
	l.Node, _ = parser.ParseFile(l.FSet, l.Config.LunarConfig.FilePath, nil, parser.ParseComments)

	log.Println("Finish parsing code")

	// TODO: handle error
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func (l *LunarUML) PrintNodeForDebug() {
	if l.Config.LunarConfig.IsDebug {
		log.Println("Debug mode on")
		// Create a new file
		f, err := os.Create("output/ast.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		// Redirect the output of ast.Print to the file
		err = ast.Fprint(f, l.FSet, l.Node, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("Debug mode off")
	}
}
