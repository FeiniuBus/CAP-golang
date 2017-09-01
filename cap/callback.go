package cap

// Callback provide Sample Handler.
type Callback struct {
}

type CallbackInterface interface {
	Handle(msg interface{}) error
}
