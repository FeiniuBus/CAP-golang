package cap

import (
	"errors"
)

// Callback provide Sample Handler.
type Callback struct {
}

type CallbackInterface interface {
	Handle() error
}

func (this *Callback) Handle() error {
	return errors.New("Custome callback must impl CallbackInterface")
}
