package main

import (
	"gopkg.in/yaml.v3"
	"log"
	"lunar_uml/consts"
	"lunar_uml/service"
	"os"
)

func main() {
	initLogConfig()
	// TODO: what is the difference whether there is & or not
	lunar := &service.LunarUML{
		Config: loadConfig(),
	}
	lunar.InitUML()
	lunar.ParseCodeWithAST()
	lunar.PrintNodeForDebug()
	// start traversing
	lunar.Inspect()
	lunar.PrintUML()
}

func initLogConfig() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
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
