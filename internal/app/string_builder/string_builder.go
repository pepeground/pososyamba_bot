package string_builder

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"math/rand"
	"strings"
)

type StringBuilder struct {
	Config *viper.Viper
}

func (sb StringBuilder) FormattedUsername(message *tgbotapi.Message) string {
	if message.From.UserName == "" {
		return message.From.String()
	}

	return "@" + message.From.UserName
}

func (sb StringBuilder) GenerateGayID() string {
	names := sb.Config.GetStringSlice("gay_names")
	adjectives := sb.Config.GetStringSlice("gay_adjectives")

	return fmt.Sprintf("%s_%s_%v",
		adjectives[rand.Intn(len(adjectives))],
		names[rand.Intn(len(names))],
		rand.Intn(10000),
	)
}

func (sb StringBuilder) BuildPososyamba() string {
	var text []string

	mainPososyamba := sb.Config.GetString("main_pososyamba")

	for _, elem := range mainPososyamba {
		text = append(text, string(elem))
	}

	return strings.Join(text, "\r\n\n")
}
