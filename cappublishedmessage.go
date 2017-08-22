package cap

type CapPublishedMessage struct{
	Id int
	Name string
	Content string
	Added int
	ExpiresAt int
	Retries int
	StatusName string
}
