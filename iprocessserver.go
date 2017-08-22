package cap

type IProcessServer interface {
	Pulse()
	Start()
	Close()
}