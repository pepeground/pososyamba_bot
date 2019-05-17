package main

import (
	"github.com/thesunwave/pososyamba_bot/configs"
	botClient "github.com/thesunwave/pososyamba_bot/internal/app/bot"
)

func main() {
	configs.Init()
	botClient.Init()
}
