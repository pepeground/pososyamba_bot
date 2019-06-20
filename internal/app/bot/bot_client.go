package bot_client

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/thesunwave/pososyamba_bot/configs"
	"github.com/thesunwave/pososyamba_bot/internal/app/admin"
	"github.com/thesunwave/pososyamba_bot/internal/app/analytics"
	"github.com/thesunwave/pososyamba_bot/internal/app/commands"
	"github.com/thesunwave/pososyamba_bot/internal/app/fakenews"
	"github.com/thesunwave/pososyamba_bot/internal/app/string_builder"

	"math/rand"
	"os"
)

type BotClient struct {
	Config *viper.Viper
	Bot    *tgbotapi.BotAPI
}

func Init() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))

	if err != nil {
		log.Fatal().Err(err).Msg("Telegram client init has error")
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	botClient := BotClient{
		Config: configs.GetConfig(),
		Bot:    bot,
	}

	botClient.run()
}

func (c BotClient) run() {
	sb := string_builder.StringBuilder{Config: c.Config}
	bot := c.Bot

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Error().Err(err).Timestamp()
	}

	preparedPhrases := sb.Config.GetStringSlice("prepared_phrases")

	for update := range updates {
		log.Printf("%+v\n", update)

		if update.InlineQuery != nil {
			go inlineQueryHandler(bot, update, preparedPhrases, sb)
			continue
		}

		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		go messageCommandHandler(&update, &c)
	}
}

func inlineQueryHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update, preparedPhrases []string, sb string_builder.StringBuilder) {
	article := tgbotapi.NewInlineQueryResultArticle(
		cast.ToString(rand.Intn(1000000)),
		"Выпустить пососямбу в этот чат",
		preparedPhrases[rand.Intn(len(preparedPhrases))]+"\r\n\n"+sb.BuildPososyamba(),
	)

	title, _ := fakenews.FetchTitle()
	fakeNews := tgbotapi.NewInlineQueryResultArticle(
		cast.ToString(rand.Intn(1000000)),
		"Сгенерить фейкньюс",
		title,
	)

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		Results:       []interface{}{article, fakeNews},
	}

	query := update.InlineQuery

	go analytics.SendToInflux(query.From.String(), query.From.ID, 0, "", "inline", "inline")

	_, err := bot.AnswerInlineQuery(inlineConf)

	log.Info().Interface("update", update)

	if err != nil {
		log.Error().Err(err)
	}
}

func messageCommandHandler(update *tgbotapi.Update, botClient *BotClient) {
	var messages *[]tgbotapi.MessageConfig

	handlers := commands.RequiredParams{
		Update:        update,
		StringBuilder: string_builder.GetBuilder(),
		Config:        configs.GetConfig(),
	}

	adminHandlers := admin.RequiredParams{
		Update:        update,
		StringBuilder: string_builder.GetBuilder(),
		Config:        configs.GetConfig(),
	}

	switch update.Message.Command() {
	case "pososyamba":
		messages = handlers.Pososyamba()
	case "gay_id":
		messages = handlers.GayID()
	case "mraz_id":
		messages = handlers.MrazID()
	case "renew_gay_id":
		messages = handlers.RenewGayID()
	case "change_gay_id":
		messages = adminHandlers.ChangeGayID()
	case "hot_news":
		messages = handlers.HotNews()
	}

	go botClient.sendMessage(*messages)
}

func (c *BotClient) sendMessage(messages []tgbotapi.MessageConfig) {
	for _, message := range messages {
		if _, err := c.Bot.Send(message); err != nil {
			log.Fatal().Err(err).Msg("Something went wrong")
		}
	}
}
