package cap

import (
	"time"
)

type CapReceivedMessage struct{
	Id int
	Name string
	Group string
	Content string
	Added int64
	ExpiresAt int
	Retries int
	StatusName string
	LastWarnedTime int
	MessageId int64
	TransactionId int64
}

func NewCapReceivedMessage(context MessageContext) *CapReceivedMessage {
	return &CapReceivedMessage{
		Group: context.Group ,
		Name: context.Name ,
		Content: context.Content ,
		Added: time.Now().Unix() ,
		ExpiresAt: 0 ,
		Retries: 0 ,
	}
}