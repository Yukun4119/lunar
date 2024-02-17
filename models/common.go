package models

type LunarConfig struct {
	CurService string `yaml:"curService"`
	TargetInf  string `yaml:"targetInf"`
	FilePath   string `yaml:"filePath"`
	IsDebug    bool   `yaml:"isDebug"`
}

type YamlConfig struct {
	LunarConfig LunarConfig `yaml:"lunarConfig"`
}
