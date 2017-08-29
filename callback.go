package cap

import (
	"errors"
)

type Callback struct {
}

type CallbackInterface interface {
	Handle(msg interface{}) error
}

func (this *Callback) Handle(msg interface{}) error {
	return errors.New("Custome callback must impl CallbackInterface")
}
