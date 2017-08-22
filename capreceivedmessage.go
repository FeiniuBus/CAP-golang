package cap

type CapReceivedMessage struct{
	Id int
	Name string
	Group string
	Content string
	Added int
	ExpiresAt int
	Retries int
	StatusName string
}