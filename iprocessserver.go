package cap

import (
	"sync"
)

// IProcessServer ...
type IProcessServer interface {
	Start()
	WaitForClose(wg *sync.WaitGroup)
}
