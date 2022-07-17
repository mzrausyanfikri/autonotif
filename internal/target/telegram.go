package target

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aimzeter/autonotif/entity"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const maxRetry = 5

type Telegram struct {
	channelID int64
	botclient *tgbotapi.BotAPI
}

func NewTelegram(secret string, channelID int64) (*Telegram, error) {
	bot, err := tgbotapi.NewBotAPI(secret)
	if err != nil {
		return nil, err
	}

	return &Telegram{
		channelID: channelID,
		botclient: bot,
	}, nil
}

func (t *Telegram) SendMessage(ctx context.Context, p entity.Proposal) error {
	return t.sendMessage(ctx, t.channelID, p)
}

func (t *Telegram) sendMessage(ctx context.Context, chatID int64, p entity.Proposal) error {
	text, err := composeMessageText(p)
	if err != nil {
		return err
	}

	if len(text) <= msgLimit {
		msg := tgbotapi.NewMessage(chatID, text)
		return t.botSend(msg)
	}

	lower := 0
	for {
		upper := lower + msgLimit
		if upper > len(text) {
			upper = len(text)
		}

		err := t.botSend(tgbotapi.NewMessage(chatID, text[lower:upper]))
		if err != nil {
			return err
		}

		lower = upper
		if upper >= len(text) {
			break
		}
	}

	return nil
}

func (t *Telegram) botSend(msg tgbotapi.MessageConfig) error {
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true

	var err error

	for attempt := 1; ; attempt++ {
		_, err = t.botclient.Send(msg)
		if err == nil {
			break
		}

		if attempt >= maxRetry {
			return fmt.Errorf("failed to send message: %s", err.Error())
		}
	}

	return nil
}

func composeMessageText(p entity.Proposal) (string, error) {
	var detail entity.ProposalRawData_COSMOS_v1

	jsonErr := json.Unmarshal([]byte(p.RawData), &detail)
	if jsonErr != nil {
		return "", jsonErr
	}

	typeSlice := strings.Split(detail.Proposal.Content.Type, ".")
	text := fmt.Sprintf(messageTemplate,
		detail.Proposal.ProposalID,
		detail.Proposal.Content.Title,
		typeSlice[len(typeSlice)-1],
		detail.Proposal.VotingStartTime.In(timeLoc).Format("2006-01-02 15:04:05")+" WIB",
		detail.Proposal.VotingEndTime.In(timeLoc).Format("2006-01-02 15:04:05")+" WIB",
		detail.Proposal.Content.Description,
		detail.Proposal.ProposalID,
	)

	return text, nil
}
