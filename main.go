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
		log.Fatal("load config error")
	}

	// Step 2
	// parse code
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, Config.LunarConfig.FilePath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	ast.Print(fset, node)

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
				service.TranverseFunc(fn)
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
	file, err := os.ReadFile("./config/config.yml")
	if err != nil {
		log.Println("read file error")
		return err
	}
	return yaml.Unmarshal(file, &Config)
}
