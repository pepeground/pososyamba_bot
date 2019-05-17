package bot_client

import (
	"github.com/spf13/viper"
	"github.com/thesunwave/pososyamba_bot/configs"
)

type BotClient struct {
	config *viper.Viper
}

func (c BotClient) Init() {
	configs.GetConfig()
}

func (c BotClient) SendMessage(message string) {

}
