package cap

type CapPublishedMessage struct{
	Id int
	Name string
	Content string
	Added int
	ExpiresAt int
	Retries int
	StatusName string
	LastWarnedTime int
	MessageId int64
	TransactionId int64
}
