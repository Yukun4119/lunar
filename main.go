package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"gopkg.in/yaml.v3"
	"log"
	"lunar_uml/service"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	lunar := service.LunarUML{
		Config: loadConfig(),
	}

	// Step 2
	// parse code
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, lunar.Config.LunarConfig.FilePath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	ast.Print(fset, node)

	// Step 3
	// Get all the participants
	lunar.Participants = append(lunar.Participants, lunar.Config.LunarConfig.TargetInf)

	// step 4
	// Start PlantUML sequence diagram
	lunar.PlantUML = append(lunar.PlantUML, "@startuml")
	lunar.PlantUML = append(lunar.PlantUML, "autonumber")

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
	lunar.PlantUML = append(lunar.PlantUML, "@enduml")

	// Print PlantUML
	for _, line := range lunar.PlantUML {
		println(line)
	}
}

func loadConfig() service.YamlConfig {
	file, err := os.ReadFile("./config/config.yml")
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
