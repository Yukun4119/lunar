package main

import (
	"github.com/ryqdev/golang_utils/log"
	"lunar_uml/consts"
	"lunar_uml/service"
	"lunar_uml/util"
)

func main() {
	log.Info("Starting...")
	lunar := &service.LunarUML{
		Config: util.LoadConfig(consts.RelativeYamlConfigFilePath),
	}
	lunar.InitUML().ParseCodeWithAST().PrintNodeForDebug().Inspect().OutputUML()
}
