package telegramservice

import (
	"bytes"
	"context"
	"text/template"

	"go.uber.org/zap"

	clients "github.com/Mahno9/GoMicroservicesCourse/notification/internal/client/http"
	"github.com/Mahno9/GoMicroservicesCourse/notification/internal/config"
	services "github.com/Mahno9/GoMicroservicesCourse/notification/internal/service"
	templates "github.com/Mahno9/GoMicroservicesCourse/notification/internal/service/telegram_service/templates"
	"github.com/Mahno9/GoMicroservicesCourse/notification/model"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
)

const (
	HelloMessage = "👋 Hi, I'm notification bot."
)

var (
	shipAssembledTemplate = template.Must(template.ParseFS(templates.FS, "ship_assembled_notification.tmpl"))
	orderPaidTemplate     = template.Must(template.ParseFS(templates.FS, "order_paid_notification.tmpl"))
)

type service struct {
	telegramClient clients.TelegramClient
	chatID         int64
}

func NewService(ctx context.Context, telegramClient clients.TelegramClient, cfg config.TelegramConfig) services.TelegramService {
	s := &service{
		telegramClient: telegramClient,
		chatID:         cfg.ChatID(),
	}

	err := s.telegramClient.SetStartHandler(ctx, s.startHandler)
	if err != nil {
		logger.Warn(ctx, "failed to set start handler", zap.Error(err))
	}

	return s
}

func (s *service) startHandler(ctx context.Context) error {
	err := s.telegramClient.SendMessage(ctx, s.chatID, HelloMessage)
	return err
}

func (s *service) SendShipAssembledMessage(ctx context.Context, event model.ShipAssembledEvent) error {
	message, err := s.formatShipAssembledMessage(event)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, s.chatID, message)
	return err
}

func (s *service) SendOrderPaidMessage(ctx context.Context, event model.OrderPaidEvent) error {
	message, err := s.formatOrderPaidMessage(event)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, s.chatID, message)
	return err
}

func (s *service) formatOrderPaidMessage(event model.OrderPaidEvent) (string, error) {
	var buf bytes.Buffer
	err := orderPaidTemplate.Execute(&buf, event)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *service) formatShipAssembledMessage(event model.ShipAssembledEvent) (string, error) {
	var buf bytes.Buffer
	err := shipAssembledTemplate.Execute(&buf, event)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
