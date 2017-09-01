package cap

// IProcessor bla.
type IProcessor interface {
	Process(context *ProcessingContext) (*ProcessResult, error)
}
