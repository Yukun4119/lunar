package service

import (
	"bufio"
	"fmt"
	"github.com/ryqdev/golang_utils/log"
	"lunar_uml/consts"
	"os"
)

func (l *LunarUML) InitUML() *LunarUML {
	if l.Err != nil {
		return l
	}
	l.Participants = append(l.Participants, l.Config.LunarConfig.TargetInf)
	l.PlantUML = append(l.PlantUML, consts.UmlStartuml)
	l.PlantUML = append(l.PlantUML, consts.UmlSetAutonumber)
	log.Info("Finish init UML")
	return l
}

func (l *LunarUML) OutputUML() *LunarUML {
	if l.Err != nil {
		return l
	}
	l.PlantUML = append(l.PlantUML, consts.UmlEnduml)
	outputFile, err := os.Create("output/output.puml")
	if err != nil {
		fmt.Println("Error:", err)
		return l
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	for _, line := range l.PlantUML {
		fmt.Fprintln(writer, line)
	}
	writer.Flush()

	log.Info("Finish writing uml file")
	return l
}
