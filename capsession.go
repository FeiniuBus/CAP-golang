package cap

type CapSession struct{
	Options CapOptions
}

func NewSession(options CapOptions) *CapSession{
	session := &CapSession{}
	session.Options = options
	return session
}