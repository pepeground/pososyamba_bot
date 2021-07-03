package bot_client

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/thesunwave/pososyamba_bot/configs"
	"github.com/thesunwave/pososyamba_bot/internal/app/admin"
	"github.com/thesunwave/pososyamba_bot/internal/app/analytics"
	"github.com/thesunwave/pososyamba_bot/internal/app/cache"
	"github.com/thesunwave/pososyamba_bot/internal/app/commands"
	"github.com/thesunwave/pososyamba_bot/internal/app/external/tenor"
	"github.com/thesunwave/pososyamba_bot/internal/app/fakenews"
	"github.com/thesunwave/pososyamba_bot/internal/app/mrkshi"
	"github.com/thesunwave/pososyamba_bot/internal/app/string_builder"
	"gopkg.in/yaml.v2"
	"time"

	"io/ioutil"
	"math/rand"
	"os"
)

var mrkshi_phrases []string

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

	file, err := ioutil.ReadFile("configs/mrkshi_phrases.yml")
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	err = yaml.Unmarshal(file, &mrkshi_phrases)
	if err != nil {
		log.Error().Err(err).Msg("")
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

		if update.ChannelPost != nil {
			if update.ChannelPost.Text != "" { // if channel post is a plain text
				go mrkshi.UpdatePhrases(update.ChannelPost.Text, &mrkshi_phrases)
				continue
			}

			if update.ChannelPost.Caption != "" { // if channel post is a photo with caption
				go mrkshi.UpdatePhrases(update.ChannelPost.Caption, &mrkshi_phrases)
				continue
			}

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

	url, err := tenor.GetGifsByIDs(os.Getenv("FUNERAL_GIFS_IDS"))
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	gif := tgbotapi.NewInlineQueryResultGIF(os.Getenv("FUNERAL_GIFS_IDS"), url)

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		Results:       []interface{}{article, fakeNews, gif},
	}

	query := update.InlineQuery

	go analytics.SendToInflux(query.From.String(), query.From.ID, 0, "", "inline", "inline")

	_, err = bot.AnswerInlineQuery(inlineConf)

	log.Info().Interface("update", update)

	if err != nil {
		log.Error().Err(err)
	}
}

func messageCommandHandler(update *tgbotapi.Update, botClient *BotClient) {
	var messages interface{}
	//var messages *[]tgbotapi.MessageConfig

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
	case "start":
		messages = handlers.Start()
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
	case "count_news":
		messages = adminHandlers.CountNews()
	case "flush_hot_news":
		messages = adminHandlers.FlushHotNews()
	case "hot_news":
		messages = handlers.HotNews()
	case "f":
		messages = handlers.F()
	case "F":
		messages = handlers.NewF()
	case "MRKSHI", "mrkshi":
		messages = handlers.MRKSHI(&mrkshi_phrases)
	default:
		return
	}

	go botClient.sendMessage(messages)
}

func (c *BotClient) sendMessage(messages interface{}) {
	textMessages, ok := messages.(*[]tgbotapi.MessageConfig)
	if ok {
		for _, message := range *textMessages {
			if _, err := c.Bot.Send(message); err != nil {
				log.Fatal().Err(err).Msg("Something went wrong")
			}
		}
	}

	pictureMessages, ok := messages.(*[]tgbotapi.AnimationConfig)
	if ok {
		for _, message := range *pictureMessages {
			if _, err := c.Bot.Send(message); err != nil {
				log.Fatal().Err(err).Msg("Something went wrong")
			}
		}
	}

	videoNoteMessages, ok := messages.(*[]tgbotapi.VideoNoteConfig)
	if ok {
		for _, message := range *videoNoteMessages {
			msg, err := c.Bot.Send(message)

			if err != nil {
				log.Fatal().Err(err).Msg("Something went wrong")
			}
			cache.Redis().Set("funeral_video_id", msg.VideoNote.FileID, 100*time.Hour)
		}
	}
}
