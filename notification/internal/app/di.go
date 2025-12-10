package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-telegram/bot"

	clients "github.com/Mahno9/GoMicroservicesCourse/notification/internal/client/http"
	telegramClient "github.com/Mahno9/GoMicroservicesCourse/notification/internal/client/http/telegram"
	"github.com/Mahno9/GoMicroservicesCourse/notification/internal/config"
	kafkaDecoders "github.com/Mahno9/GoMicroservicesCourse/notification/internal/converter/kafka"
	kafkaDecoder "github.com/Mahno9/GoMicroservicesCourse/notification/internal/converter/kafka/decoder"
	services "github.com/Mahno9/GoMicroservicesCourse/notification/internal/service"
	orderPaidConsumerService "github.com/Mahno9/GoMicroservicesCourse/notification/internal/service/consumer/order_paid"
	shipAssembledConsumerService "github.com/Mahno9/GoMicroservicesCourse/notification/internal/service/consumer/ship_assembled"
	telegramService "github.com/Mahno9/GoMicroservicesCourse/notification/internal/service/telegram_service"
	kafkaWrapped "github.com/Mahno9/GoMicroservicesCourse/platform/pkg/kafka"
	kafkaWrappedConsumer "github.com/Mahno9/GoMicroservicesCourse/platform/pkg/kafka/consumer"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
)

type diContainer struct {
	config *config.Config

	orderPaidConsumerService services.OrderPaidConsumerService
	orderPaidConsumer        kafkaWrapped.Consumer
	orderPaidConsumerGroup   sarama.ConsumerGroup
	orderPaidDecoder         kafkaDecoders.OrderPaidDecoder

	shipAssembledConsumerService services.ShipAssembledConsumerService
	shipAssembledConsumer        kafkaWrapped.Consumer
	shipAssembledConsumerGroup   sarama.ConsumerGroup
	shipAssembledDecoder         kafkaDecoders.ShipAssembledDecoder

	telegramService services.TelegramService
	telegramClient  clients.TelegramClient
	telegramBot     *bot.Bot
}

func NewDIContainer(cfg *config.Config) *diContainer {
	return &diContainer{config: cfg}
}

func (c *diContainer) OrderPaidConsumerService(ctx context.Context) services.OrderPaidConsumerService {
	if c.orderPaidConsumerService == nil {
		c.orderPaidConsumerService = orderPaidConsumerService.NewService(
			c.OrderPaidConsumer(ctx),
			c.OrderPaidDecoder(),
			c.TelegramService(ctx),
		)
	}
	return c.orderPaidConsumerService
}

func (c *diContainer) OrderPaidConsumer(ctx context.Context) kafkaWrapped.Consumer {
	if c.orderPaidConsumer == nil {
		c.orderPaidConsumer = kafkaWrappedConsumer.NewConsumer(
			c.OrderPaidConsumerGroup(),
			[]string{c.config.OrderPaidConsumer.TopicName()},
			logger.Logger(),
		)
	}
	return c.orderPaidConsumer
}

func (c *diContainer) OrderPaidConsumerGroup() sarama.ConsumerGroup {
	if c.orderPaidConsumerGroup != nil {
		return c.orderPaidConsumerGroup
	}

	var err error
	c.orderPaidConsumerGroup, err = sarama.NewConsumerGroup(
		c.config.Kafka.Addresses(),
		c.config.OrderPaidConsumer.ConsumerGroupId(),
		c.config.OrderPaidConsumer.Config(),
	)
	if err != nil {
		panic(fmt.Sprintf("❗ failed to create consumer group: %v", err))
	}

	return c.orderPaidConsumerGroup
}

func (c *diContainer) OrderPaidDecoder() kafkaDecoders.OrderPaidDecoder {
	if c.orderPaidDecoder == nil {
		c.orderPaidDecoder = kafkaDecoder.NewOrderPaidDecoder()
	}

	return c.orderPaidDecoder
}

func (c *diContainer) TelegramClient() clients.TelegramClient {
	if c.telegramClient == nil {
		c.telegramClient = telegramClient.NewClient(c.TelegramBot())
	}
	return c.telegramClient
}

func (c *diContainer) TelegramBot() *bot.Bot {
	if c.telegramBot == nil {
		opts := []bot.Option{
			bot.WithHTTPClient(10*time.Second, &http.Client{
				Timeout: 10 * time.Second,
			}),
			bot.WithSkipGetMe(),
		}

		var err error
		c.telegramBot, err = bot.New(c.config.Telegram.BotToken(), opts...)
		if err != nil {
			panic(fmt.Errorf("❗ failed to create bot: %w", err))
		}
	}

	return c.telegramBot
}

func (c *diContainer) TelegramService(ctx context.Context) services.TelegramService {
	if c.telegramService == nil {
		c.telegramService = telegramService.NewService(ctx, c.TelegramClient(), c.config.Telegram)
	}
	return c.telegramService
}

func (c *diContainer) ShipAssembledConsumerService(ctx context.Context) services.ShipAssembledConsumerService {
	if c.shipAssembledConsumerService == nil {
		c.shipAssembledConsumerService = shipAssembledConsumerService.NewService(
			c.ShipAssembledConsumer(ctx),
			c.ShipAssembledDecoder(),
			c.TelegramService(ctx),
		)
	}

	return c.shipAssembledConsumerService
}

func (c *diContainer) ShipAssembledConsumer(ctx context.Context) kafkaWrapped.Consumer {
	if c.shipAssembledConsumer == nil {
		c.shipAssembledConsumer = kafkaWrappedConsumer.NewConsumer(
			c.ShipAssembledConsumerGroup(),
			[]string{c.config.ShipAssembledConsumer.TopicName()},
			logger.Logger(),
		)
	}
	return c.shipAssembledConsumer
}

func (c *diContainer) ShipAssembledConsumerGroup() sarama.ConsumerGroup {
	if c.shipAssembledConsumerGroup != nil {
		return c.shipAssembledConsumerGroup
	}

	var err error
	c.shipAssembledConsumerGroup, err = sarama.NewConsumerGroup(
		c.config.Kafka.Addresses(),
		c.config.ShipAssembledConsumer.ConsumerGroupId(),
		c.config.ShipAssembledConsumer.Config(),
	)
	if err != nil {
		panic(fmt.Sprintf("❗ failed to create consumer group: %v", err))
	}

	return c.shipAssembledConsumerGroup
}

func (c *diContainer) ShipAssembledDecoder() kafkaDecoders.ShipAssembledDecoder {
	if c.shipAssembledDecoder == nil {
		c.shipAssembledDecoder = kafkaDecoder.NewShipAssembledDecoder()
	}

	return c.shipAssembledDecoder
}
