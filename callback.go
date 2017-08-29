package cap

import (
	"errors"
)

// Callback provide Sample Handler.
type Callback struct {
}

type CallbackInterface interface {
	Handle(msg interface{}) error
}

func (this *Callback) Handle(msg interface{}) error {
	return errors.New("Custome callback must impl CallbackInterface")
}
