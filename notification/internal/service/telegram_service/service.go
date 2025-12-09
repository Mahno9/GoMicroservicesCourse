package telegramservice

import (
	"context"

	clients "github.com/Mahno9/GoMicroservicesCourse/notification/internal/client/http"
	services "github.com/Mahno9/GoMicroservicesCourse/notification/internal/service"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
	"go.uber.org/zap"
)

type service struct {
	telegramClient clients.TelegramClient
	savedChatIds   []int64
}

func NewService(telegramClient clients.TelegramClient) services.TelegramService {
	s := &service{
		telegramClient: telegramClient,
	}

	// subscribe for new chats
	s.telegramClient.RegisterNewChatSubscriber(s)

	return s
}

func (s *service) BroadcastMessage(ctx context.Context, message string) error {
	for _, chatId := range s.savedChatIds {
		err := s.telegramClient.SendMessage(ctx, chatId, message)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) NewChatStarted(ctx context.Context, chatId int64) {
	s.savedChatIds = append(s.savedChatIds, chatId)
	logger.Info(ctx, "📝 New chat registered", zap.Int64("ChatId", chatId))

	s.telegramClient.SendMessage(ctx, chatId, "User registered.")
}
