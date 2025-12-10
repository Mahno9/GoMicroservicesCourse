package env

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
)

type orderPaidProducerConfigEnv struct {
	OrderPaidTopicName string `env:"ORDER_PAID_TOPIC_NAME,required"`
}

type orderPaidProducerConfig struct {
	orderPaidTopicName string
}

func NewOrderPaidProducerConfig() (*orderPaidProducerConfig, error) {
	var raw orderPaidProducerConfigEnv
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	err := raw.checkValues()
	if err != nil {
		return nil, err
	}

	return &orderPaidProducerConfig{
		orderPaidTopicName: raw.OrderPaidTopicName,
	}, nil
}

func (c *orderPaidProducerConfigEnv) checkValues() error {
	if c.OrderPaidTopicName == "" {
		return fmt.Errorf("%w: %s", model.ErrInvalidConfigValue, "ORDER_PAID_TOPIC_NAME")
	}

	return nil
}

func (c *orderPaidProducerConfig) TopicName() string {
	return c.orderPaidTopicName
}

func (c *orderPaidProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true // TODO: check wether it need

	return config
}
