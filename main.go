package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"gopkg.in/yaml.v3"
	"log"
	"lunar_uml/consts"
	"lunar_uml/service"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	lunar := service.LunarUML{
		Config: loadConfig(),
	}

	// parse code
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, lunar.Config.LunarConfig.FilePath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	err = ast.Print(fset, node)
	if err != nil {
		log.Fatal(err)
	}

	// Get all the participants
	lunar.Participants = append(lunar.Participants, lunar.Config.LunarConfig.TargetInf)

	// Start PlantUML sequence diagram
	lunar.PlantUML = append(lunar.PlantUML, consts.UmlStartuml)
	lunar.PlantUML = append(lunar.PlantUML, consts.UmlSetAutonumber)

	// TODO: make the code neater
	ast.Inspect(node, func(n ast.Node) bool {
		switch fn := n.(type) {
		case *ast.FuncDecl:
			if fn.Name.Name == lunar.Config.LunarConfig.TargetInf {
				lunar.TranverseFunc(fn)
			}
		}
		return true
	})

	// End PlantUML sequence diagram
	lunar.PlantUML = append(lunar.PlantUML, consts.UmlEnduml)

	// Print PlantUML
	for _, line := range lunar.PlantUML {
		// TODO: output to a file
		println(line)
	}
}

func loadConfig() service.YamlConfig {
	file, err := os.ReadFile(consts.RelativeYamlConfigFilePath)
	if err != nil {
		log.Fatal("read file error")
	}
	config := service.YamlConfig{}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatal("unmarshal error")
	}
	return config
}
