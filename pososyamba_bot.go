package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/influxdata/influxdb1-client/v2"
	"log"
	"math/rand"
	"os"
	"strconv"
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

var gayName = []string{
	"faggot",
	"pidor",
	"mraz",
	"petushok",
	"cocksucker",
	"analsucker",
	"pussylicker",
	"weabo",
	"cock",
}

var gayAdjective = []string{
	"pretty",
	"fat",
	"cool",
	"sweety",
	"furryloving",
	"cuckholdy",
}

var redisdb = redis.NewClient(&redis.Options{
	Addr:     "redis:6379", // use default Addr
	Password: "",           // no password set
	DB:       0,            // use default DB
})

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
		log.Printf("%+v\n", update)

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

			query := update.InlineQuery

			go sendToInflux(query.From.UserName, query.From.ID, 0, "", "inline", "inline")

			_, err := bot.AnswerInlineQuery(inlineConf)

			log.Println(update)

			if err != nil {
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

		message := update.Message

		// Create a new MessageConfig. We don't have text yet,
		// so we should leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		rand.Seed(time.Now().Unix())

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "pososyamba":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.Text = preparePhrases[rand.Intn(len(preparePhrases))]

			go sendToInflux(message.From.UserName, message.From.ID, message.Chat.ID, message.Chat.Title, "message", "pososyamba")

			sendMessage(msg, bot)

			msg.Text = buildPososyamba()

			sendMessage(msg, bot)
		case "gay_id":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "")

			forwardedMessage := update.Message

			var gayID, username string
			var clientID int
			var err error

			if forwardedMessage.ReplyToMessage != nil {
				username = forwardedMessage.ReplyToMessage.From.UserName
				clientID = forwardedMessage.ReplyToMessage.From.ID
				gayID, err = redisdb.Get(strconv.Itoa(clientID)).Result()

				log.Println(clientID)
				msg.ReplyToMessageID = forwardedMessage.ReplyToMessage.MessageID
				log.Println(gayID)
			} else {
				username = forwardedMessage.From.UserName
				clientID = forwardedMessage.From.ID
				gayID, err = redisdb.Get(strconv.Itoa(clientID)).Result()
				log.Println(clientID)
				log.Println(gayID)
			}

			if err != nil {
				msg.Text = generateGayID()

				err := redisdb.Set(strconv.Itoa(clientID), msg.Text, 0).Err()

				if err != nil {
					log.Println(err)
				}
			} else {
				msg.Text = gayID
			}

			msg.Text = "@" + username + " has gay_id: #" + msg.Text

			go sendToInflux(message.From.UserName, message.From.ID, message.Chat.ID, message.Chat.Title, "message", "gay_id")

			sendMessage(msg, bot)

		case "renew_gay_id":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "")

			gayID := generateGayID()

			msg.Text = "@" + update.Message.From.UserName + " you have updated gay_id: #" + gayID

			err := redisdb.Set(strconv.Itoa(update.Message.From.ID), gayID, 0).Err()

			if err != nil {
				log.Println(err)
			}

			go sendToInflux(message.From.UserName, message.From.ID, message.Chat.ID, message.Chat.Title, "message", "renew_gay_id")

			sendMessage(msg, bot)
		}
	}
}

func generateGayID() string {
	return fmt.Sprintf("%s_%s_%v",
		gayAdjective[rand.Intn(len(gayAdjective))],
		gayName[rand.Intn(len(gayName))],
		rand.Intn(10000),
	)
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

func sendToInflux(username string, userID int, chatID int64, chatTitle, messageType, command string) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     os.Getenv("INFLUX_URL"),
		Username: os.Getenv("INFLUX_USERNAME"),
		Password: os.Getenv("INFLUX_PASSWORD"),
	})
	if err != nil {
		log.Printf("%+v\n", "Error creating InfluxDB Client: "+err.Error())
	}
	defer c.Close()

	// Create a new point batch
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database: "web_services",
	})

	// Create a point and add to batch
	tags := map[string]string{"command": command}
	fields := map[string]interface{}{
		"username":    username,
		"user_id":     userID,
		"chat_title":  chatTitle,
		"chat_id":     chatID,
		"messageType": messageType,
	}

	pt, err := client.NewPoint("pososyamba_usage", tags, fields, time.Now())
	if err != nil {
		log.Println("Error: ", err.Error())
	}

	if os.Getenv("ENVIRONMENT") == "production" {
		bp.AddPoint(pt)

		// Write the batch
		c.Write(bp)
	}
}
