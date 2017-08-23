package cap

type DispatcherInterface interface{
	ProcessorInterface

	GetWaiting() bool;
}