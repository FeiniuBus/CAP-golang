package cap_rabbitmq

type RabbitMQOptions struct {
	ConnectionTimeoutSecs				int32
	Password							string
	UserName							string
	VirtualHost							string
	TopicExchangeName					string
	RequestedConnectionTimeout			int32
	HostName							string
	SocketReadTimeout					int32
	SocketWriteTimeout					int32
	Port								int32
	QueueMessageExpires 				int32
	ExchangeType						string
}

var (
	RabbitMQConfig *RabbitMQOptions
)

func init() {
	RabbitMQConfig = defaultRabbitMQConfig()
}

func defaultRabbitMQConfig() *RabbitMQOptions {
	var defaultTimeoutSecs int32 = 30 * 1000

	return &RabbitMQOptions {
		ConnectionTimeoutSecs: defaultTimeoutSecs,
		Password: "guest",
		UserName: "guest",
		VirtualHost: "/",
		TopicExchangeName: "cap.default.topic",
		RequestedConnectionTimeout: defaultTimeoutSecs,
		HostName: "localhost",
		SocketReadTimeout: defaultTimeoutSecs,
		SocketWriteTimeout: defaultTimeoutSecs,
		Port: -1,
		QueueMessageExpires: 864000000,
		ExchangeType: "topic",
	}
}

func (this *RabbitMQOptions) SetConnectionTimeout(timeoutSecs int32) {
	this.ConnectionTimeoutSecs = timeoutSecs
}

func (this *RabbitMQOptions) SetPassword(password string) {
	this.Password = password
}

func (this *RabbitMQOptions) SetUserName(userName string) {
	this.UserName = userName
}

func (this *RabbitMQOptions) SetVirtualHost(vHost string) {
	this.VirtualHost = vHost
}

func (this *RabbitMQOptions) SetTopicExchangeName(topicExchangeName string) {
	this.TopicExchangeName = topicExchangeName
}

func (this *RabbitMQOptions) SetRequestedConnectionTimeout(requestedConnectionTimeoutSecs int32) {
	this.RequestedConnectionTimeout = requestedConnectionTimeoutSecs
}

func (this *RabbitMQOptions) SetHostName(hostName string) {
	this.HostName = hostName
}

func (this *RabbitMQOptions) SetSocketReadTimeout(socketReadTimeoutSecs int32) {
	this.SocketReadTimeout = socketReadTimeoutSecs
}

func (this *RabbitMQOptions) SetSocketWriteTimeout(socketWriteTimeout int32) {
	this.SocketWriteTimeout = socketWriteTimeout
}

func (this *RabbitMQOptions) SetPort(port int32) {
	this.Port = port
}

func (this *RabbitMQOptions) SetQueueMessageExpires(queueMessageExpires int32) {
	this.QueueMessageExpires = queueMessageExpires
}

func (this *RabbitMQOptions) SetExchangeType(exchangeType string) {
	this.ExchangeType = exchangeType
}