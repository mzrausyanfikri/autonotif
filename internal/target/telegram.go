package target

import (
	"context"
	"fmt"
	"strings"
	"time"

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

func (t *Telegram) SendMessage(ctx context.Context, p *entity.Proposal) error {
	return t.sendMessage(ctx, t.channelID, p)
}

func (t *Telegram) sendMessage(ctx context.Context, chatID int64, p *entity.Proposal) error {
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

func composeMessageText(p *entity.Proposal) (string, error) {
	messageStructure := p.ChainConfig.ChainMessageStructure

	typeStr := p.Data.Path(messageStructure.Type.JSONPath).Data().(string)
	typeSlice := strings.Split(typeStr, ".")

	st := p.Data.Path(messageStructure.StartTime.JSONPath).Data().(string)
	startTime, _ := time.Parse(time.RFC3339, st)

	et := p.Data.Path(messageStructure.EndTime.JSONPath).Data().(string)
	endTime, _ := time.Parse(time.RFC3339, et)

	text := fmt.Sprintf(messageTemplate,
		messageStructure.Name.Const,
		p.Data.Path(messageStructure.ProposalID.JSONPath).Data(),
		p.Data.Path(messageStructure.Title.JSONPath).Data(),
		p.Data.Path(messageStructure.Status.JSONPath).Data(),
		typeSlice[len(typeSlice)-1],
		startTime.In(timeLoc).Format("2006-01-02 15:04:05")+" WIB",
		endTime.In(timeLoc).Format("2006-01-02 15:04:05")+" WIB",
		messageStructure.ViewLink.Const,
		p.Data.Path(messageStructure.ProposalID.JSONPath).Data(),
		messageStructure.ViewLink.Const,
	)

	return text, nil
}
