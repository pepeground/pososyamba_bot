package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var preparePhrases = []string{
	"Выпускаю пососямбу",
	"Петусямба пососямба",
	"Писька пососямба",
	"Знаешь пососямбу?",
	"Пососямбы давно не было в чяте (",
	"паас - пососямба аз а сервис",
	"челябинская пососямба",
	"всё, а теперь пососямба",
	"пососямбу уже пускали в чат?",
}

var mainPososyamba = "ПОСОСЯМБА"

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))

	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.InlineQuery != nil {

			article := tgbotapi.NewInlineQueryResultArticle(
				update.InlineQuery.ID,
				"Выпустить пососямбу в этот чат",
				preparePhrases[rand.Intn(len(preparePhrases))]+"\r\n\n"+buildPososyamba(),
			)

			inlineConf := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				Results:       []interface{}{article},
			}

			if _, err := bot.AnswerInlineQuery(inlineConf); err != nil {
				log.Println(err)
			}
			continue
		}

		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we should leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		rand.Seed(time.Now().Unix())

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "pososyamba":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.Text = preparePhrases[rand.Intn(len(preparePhrases))]

			sendMessage(msg, bot)

			msg.Text = buildPososyamba()

			sendMessage(msg, bot)
		}
	}
}

func sendMessage(msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) {
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func buildPososyamba() string {
	text := []string{}

	for _, elem := range mainPososyamba {
		text = append(text, string(elem))
	}

	return strings.Join(text, "\r\n\n")
}
