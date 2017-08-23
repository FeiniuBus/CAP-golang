package cap

type ProcessorInterface interface{
	Process(context ProcessingContext)
}