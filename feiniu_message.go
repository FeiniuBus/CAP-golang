package cap

// FeiniuBusMessage ...
type FeiniuBusMessage struct {
	MetaData FeiniuBusMessageMetaData `json:"meta"`
	Content  []byte]                   `json:"content"`
}

// FeiniuBusMessageMetaData ...
type FeiniuBusMessageMetaData struct {
	TransactionID int64 `json:"transaction_id"`
	MessageID     int64 `json:"message_id"`
}
