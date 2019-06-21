package admin

import (
	"github.com/thesunwave/pososyamba_bot/internal/app/cache"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/thesunwave/pososyamba_bot/internal/app/analytics"
	"github.com/thesunwave/pososyamba_bot/internal/app/string_builder"
)

type RequiredParams struct {
	Update        *tgbotapi.Update
	StringBuilder *string_builder.StringBuilder
	Config        *viper.Viper
}

func (params *RequiredParams) FlushHotNews() *[]tgbotapi.MessageConfig {
	var messages []tgbotapi.MessageConfig

	message := params.Update.Message
	msg := tgbotapi.NewMessage(message.Chat.ID, "")

	go analytics.SendToInflux(message.From.String(), message.From.ID, message.Chat.ID, message.Chat.Title, "message", "flush_hot_news")

	if message.From.ID == cast.ToInt(os.Getenv("ADMIN_ROOM")) {
		err := cache.Redis().Del("news_titles").Err()
		if err != nil {
			log.Error().Err(err)
			msg.Text = "Something went wrong"
		} else {
			msg.Text = "Done"
		}
	} else {
		msg.Text = "SOSNOOLEY"
	}

	messages = append(messages, msg)

	return &messages
}

func (params *RequiredParams) CountNews() *[]tgbotapi.MessageConfig {
	var messages []tgbotapi.MessageConfig

	message := params.Update.Message

	go analytics.SendToInflux(message.From.String(), message.From.ID, message.Chat.ID, message.Chat.Title, "message", "count_news")

	if message.From.ID == cast.ToInt(os.Getenv("ADMIN_ROOM")) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "")
		msg.Text = cache.Redis().SCard("news_titles").String()
		messages = append(messages, msg)
	}

	return &messages
}

func (params *RequiredParams) ChangeGayID() *[]tgbotapi.MessageConfig {
	var messages []tgbotapi.MessageConfig
	var clientID, newGayID string
	var err error

	message := params.Update.Message
	msg := tgbotapi.NewMessage(message.Chat.ID, "")
	arguments := strings.Split(message.CommandArguments(), " ")

	log.Print(message.From.ID, cast.ToInt(os.Getenv("ADMIN_ROOM")))

	go analytics.SendToInflux(message.From.String(), message.From.ID, message.Chat.ID, message.Chat.Title, "message", "change_gay_id")

	if message.From.ID != cast.ToInt(os.Getenv("ADMIN_ROOM")) {
		return &messages
	}

	if message.ReplyToMessage != nil {
		clientID = cast.ToString(message.ReplyToMessage.From.ID)

		if len(arguments) < 1 {
			log.Error().Msgf("Error arguments nubmer: %d", len(arguments))
			msg.Text = "Error arguments nubmer: " + cast.ToString(len(arguments))
			messages = append(messages, msg)
			return &messages
		}

		newGayID = arguments[0]

		if len(newGayID) == 0 {
			msg.Text = "Empty message"
			messages = append(messages, msg)
			return &messages
		}
	} else {
		if len(arguments) != 2 {
			log.Error().Msgf("Error arguments nubmer: %d", len(arguments))
			msg.Text = "Error arguments nubmer: " + cast.ToString(len(arguments))
			messages = append(messages, msg)
			return &messages
		}

		clientID = arguments[0]
		newGayID = arguments[1]
	}

	log.Print(clientID, newGayID)

	msg.Text = newGayID

	err = cache.Redis().Set(clientID, msg.Text, 0).Err()
	if err != nil {
		log.Error().Err(err)
		msg.Text = err.Error()
	}
	messages = append(messages, msg)

	return &messages
}
