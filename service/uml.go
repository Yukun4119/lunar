package service

import (
	"bufio"
	"fmt"
	"log"
	"lunar_uml/consts"
	"os"
)

func (l *LunarUML) InitUML() {
	l.Participants = append(l.Participants, l.Config.LunarConfig.TargetInf)
	l.PlantUML = append(l.PlantUML, consts.UmlStartuml)
	l.PlantUML = append(l.PlantUML, consts.UmlSetAutonumber)
	log.Println("Finish init UML")
}

func (l *LunarUML) OutputUML() {
	// TODO: handle errors

	l.PlantUML = append(l.PlantUML, consts.UmlEnduml)
	outputFile, err := os.Create("output/output.puml")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	for _, line := range l.PlantUML {
		fmt.Fprintln(writer, line)
	}
	writer.Flush()

	log.Println("Finish writing uml file")
}
