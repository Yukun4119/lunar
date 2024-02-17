package main

import (
	"gopkg.in/yaml.v3"
	"log"
	"lunar_uml/consts"
	"lunar_uml/models"
	"lunar_uml/service"
	"os"
)

func main() {
	initLogConfig()
	lunar := &service.LunarUML{
		Config: loadConfig(),
	}
	lunar.InitUML()
	lunar.ParseCodeWithAST()
	lunar.PrintNodeForDebug()
	lunar.Inspect()
	lunar.OutputUML()
}

func initLogConfig() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	log.Println("Finish init log")
}

func loadConfig() models.YamlConfig {
	file, err := os.ReadFile(consts.RelativeYamlConfigFilePath)
	if err != nil {
		log.Fatal("read file error")
	}
	config := models.YamlConfig{}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatal("unmarshal error")
	}
	return config
}
