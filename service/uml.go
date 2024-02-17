package service

import "lunar_uml/consts"

func (l *LunarUML) InitUML() {
	l.Participants = append(l.Participants, l.Config.LunarConfig.TargetInf)
	l.PlantUML = append(l.PlantUML, consts.UmlStartuml)
	l.PlantUML = append(l.PlantUML, consts.UmlSetAutonumber)
}

func (l *LunarUML) PrintUML() {
	l.PlantUML = append(l.PlantUML, consts.UmlEnduml)

	for _, line := range l.PlantUML {
		// TODO: output to a file
		println(line)
	}
}
