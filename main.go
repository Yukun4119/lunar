package main

import (
	"github.com/Yukun4119/golang_utils/log"
	"gopkg.in/yaml.v3"
	"lunar_uml/consts"
	"lunar_uml/models"
	"lunar_uml/service"
	"os"
)

func main() {
	log.Info("Starting...")

	lunar := &service.LunarUML{
		Config: loadConfig(),
	}
	lunar.InitUML()
	lunar.ParseCodeWithAST()
	lunar.PrintNodeForDebug()
	lunar.Inspect()
	lunar.OutputUML()
}

func loadConfig() models.YamlConfig {
	file, err := os.ReadFile(consts.RelativeYamlConfigFilePath)
	if err != nil {
		log.Error("Read file error")
		os.Exit(1)
	}
	config := models.YamlConfig{}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Error("Unmarshal error")
		os.Exit(1)
	}
	return config
}
