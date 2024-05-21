package config

import "time"

type RabbitQueueSetting struct {
	QueueName                string        `json:"queueName" yaml:"queueName"`
	Durable                  bool          `json:"durable" yaml:"durable"`
	ConsumerChannelNumber    int           `json:"consumerChannelNumber" yaml:"consumerChannelNumber"`
	OfflineReconnectInterval time.Duration `json:"offlineReconnectInterval" yaml:"offlineReconnectInterval"`
	RetryCount               int           `json:"retryCount" yaml:"retryCount"`
	ExchangeType             string        `json:"exchangeType" yaml:"exchangeType"`
	ExchangeName             string        `json:"exchangeName" yaml:"exchangeName"`
	DelayedExchangeName      string        `json:"delayedExchangeName" yaml:"delayedExchangeName"`
}

type RabbitMQSetting struct {
	Enabled          bool               `json:"enabled"`
	Addr             string             `json:"addr"`
	HelloWorld       RabbitQueueSetting `json:"helloWorld" yaml:"helloWorld"`
	WorkQueue        RabbitQueueSetting `json:"workQueue" yaml:"workQueue"`
	PublishSubscribe RabbitQueueSetting `json:"publishSubscribe" yaml:"publishSubscribe"`
	Routing          RabbitQueueSetting `json:"routing" yaml:"routing"`
	Topic            RabbitQueueSetting `json:"topic" yaml:"topic"`
}

var Setting *RabbitMQSetting = &RabbitMQSetting{
	Enabled: false,
}
