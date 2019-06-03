package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/thesunwave/pososyamba_bot/internal/app/analytics"
	"github.com/thesunwave/pososyamba_bot/internal/app/cache"
	"github.com/thesunwave/pososyamba_bot/internal/app/string_builder"
	"math/rand"
	"strconv"
)

type RequiredParams struct {
	Update        *tgbotapi.Update
	StringBuilder *string_builder.StringBuilder
	Config        *viper.Viper
}

func (params RequiredParams) Pososyamba() *[]tgbotapi.MessageConfig {
	var messages []tgbotapi.MessageConfig

	message := params.Update.Message

	preparedPhrases := params.Config.GetStringSlice("prepared_phrases")

	msg := tgbotapi.NewMessage(params.Update.Message.Chat.ID, "")
	msg.Text = preparedPhrases[rand.Intn(len(preparedPhrases))]
	messages = append(messages, msg)

	go analytics.SendToInflux(message.From.String(), message.From.ID, message.Chat.ID, message.Chat.Title, "message", "pososyamba")

	msg = tgbotapi.NewMessage(params.Update.Message.Chat.ID, "")
	msg.Text = params.StringBuilder.BuildPososyamba()
	messages = append(messages, msg)

	return &messages
}

func (params RequiredParams) GayID() *[]tgbotapi.MessageConfig {
	return getID("gay_id", &params)
}

func (params RequiredParams) MrazID() *[]tgbotapi.MessageConfig {
	return getID("mraz_id", &params)
}

func getID(message_type string, params *RequiredParams) *[]tgbotapi.MessageConfig {
	var messages []tgbotapi.MessageConfig
	var gayID, username string
	var clientID int
	var err error

	message := params.Update.Message

	msg := tgbotapi.NewMessage(params.Update.Message.Chat.ID, "")

	forwardedMessage := params.Update.Message

	if forwardedMessage.ReplyToMessage != nil {
		username = params.StringBuilder.FormattedUsername(forwardedMessage.ReplyToMessage)
		clientID = forwardedMessage.ReplyToMessage.From.ID
		gayID, err = cache.Redis().Get(strconv.Itoa(clientID)).Result()

		log.Info().Str("ClientID", string(clientID))
		msg.ReplyToMessageID = forwardedMessage.ReplyToMessage.MessageID
		log.Info().Str("ClientID", gayID)
	} else {
		username = params.StringBuilder.FormattedUsername(forwardedMessage)
		clientID = forwardedMessage.From.ID
		gayID, err = cache.Redis().Get(strconv.Itoa(clientID)).Result()
		log.Info().Str("ClientID", string(clientID))
		log.Info().Str("ClientID", gayID)
	}
	// TODO: Add a handler to have an opportunity to distinguish empty key from other errors
	if err != nil {
		msg.Text = params.StringBuilder.GenerateGayID()

		err := cache.Redis().Set(strconv.Itoa(clientID), msg.Text, 0).Err()

		if err != nil {
			log.Error().Err(err)
		}
	} else {
		msg.Text = gayID
	}

	msg.Text = username + " has " + message_type + ": #" + msg.Text

	go analytics.SendToInflux(message.From.String(), message.From.ID, message.Chat.ID, message.Chat.Title, "message", message_type)

	messages = append(messages, msg)

	return &messages
}

func (params RequiredParams) RenewGayID() *[]tgbotapi.MessageConfig {
	var messages []tgbotapi.MessageConfig

	message := params.Update.Message

	msg := tgbotapi.NewMessage(message.Chat.ID, "")

	gayID := params.StringBuilder.GenerateGayID()

	log.Info().Str("GAY ID:", gayID)

	msg.Text = params.StringBuilder.FormattedUsername(message) + " you have updated gay_id: #" + gayID

	err := cache.Redis().Set(strconv.Itoa(message.From.ID), gayID, 0).Err()

	if err != nil {
		log.Error().Err(err)
	}

	go analytics.SendToInflux(message.From.String(), message.From.ID, message.Chat.ID, message.Chat.Title, "message", "renew_gay_id")

	messages = append(messages, msg)

	return &messages
}
