package cap

type InfomationLevel struct{
	levelName string
}

func NewInfomationLevel() ILogLevel{
	level := &InfomationLevel{levelName : "Infomation"}
	return level
}

func (this *InfomationLevel) GetLevelName()string{
	return this.levelName
}