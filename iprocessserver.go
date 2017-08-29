package cap

// IProcessServer ...
type IProcessServer interface {
	Start()
	WaitForClose()
}
