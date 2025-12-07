package config

import "github.com/IBM/sarama"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type PostgresConfig interface {
	URI() string
	MigrationsDir() string
}

type HttpConfig interface {
	Host() string
	Port() string
	Address() string
}

type ClientsConfig interface {
	InventoryAddress() string
	PaymentAddress() string
}

type KafkaConfig interface {
	BrokersAddresses() []string
}

type OrderPaidProducerConfig interface {
	TopicName() string
	Config() *sarama.Config
}

type ShipAssembledConsumerConfig interface {
	TopicName() string
	ConsumerGroupID() string
	Config() *sarama.Config
}
