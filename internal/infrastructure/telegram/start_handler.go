package telegram

import (
	"context"
	"log/slog"
	"time"

	userapp "kino-notifier/internal/application/user"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const registrationTimeout = 5 * time.Second
const helloMessage = "👋Привет! Добро пожаловать в кино-нотификатор.\n\nВы успешно подписались на рассылку о выходе фильмов в прокат."

// UserRegistration — входной порт сценария, необходимый обработчику /start.
type UserRegistration interface {
	Execute(ctx context.Context, input userapp.RegisterUserInput) error
}

type startHandler struct {
	userRegistration UserRegistration
	responder        *Responder
	logger           *slog.Logger
}

func newStartHandler(
	userRegistration UserRegistration,
	responder *Responder,
	logger *slog.Logger,
) *startHandler {
	return &startHandler{
		userRegistration: userRegistration,
		responder:        responder,
		logger:           logger,
	}
}

func (h *startHandler) handle(ctx context.Context, _ *tgbot.Bot, update *models.Update) {
	if update == nil || update.Message == nil {
		return
	}
	message := update.Message
	if message.Chat.Type != models.ChatTypePrivate || message.From == nil || message.From.IsBot {
		return
	}
	from := message.From
	input := userapp.RegisterUserInput{
		TelegramUserID: from.ID,
		TelegramChatID: message.Chat.ID,
		Username:       from.Username,
		FirstName:      from.FirstName,
		LastName:       from.LastName,
		LanguageCode:   from.LanguageCode,
		IsBot:          from.IsBot,
	}

	registrationCtx, cancel := context.WithTimeout(ctx, registrationTimeout)
	err := h.userRegistration.Execute(registrationCtx, input)
	cancel()
	if ctx.Err() != nil {
		return
	}
	if err != nil {
		h.logger.Error("не удалось зарегистрировать пользователя",
			"telegram_user_id", from.ID,
			"update_id", update.ID,
			"error", err,
		)
		h.responder.Reply(ctx, message.Chat.ID,
			"Не удалось завершить регистрацию. Попробуй /start ещё раз.",
		)
		return
	}

	h.responder.Reply(ctx, message.Chat.ID, helloMessage)
}
