package cap

type IProcessor interface{
	Process(context *ProcessingContext) error
}