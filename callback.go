package cap

import (
	"errors"
)

type Callback struct {
	
}

type CallbackInterface interface {
	Handle() error
}

func (this *Callback) Handle() error {
	return errors.New("Custome callback must impl CallbackInterface")
}