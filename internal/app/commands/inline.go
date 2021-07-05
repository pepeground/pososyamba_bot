package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/cast"
	"github.com/thesunwave/pososyamba_bot/internal/app/external/tenor"
	"github.com/thesunwave/pososyamba_bot/internal/app/fakenews"
	"github.com/thesunwave/pososyamba_bot/internal/app/string_builder"
	"math/rand"
	"os"
)

func InlinePososyamba(preparedPhrases []string, sb string_builder.StringBuilder) tgbotapi.InlineQueryResultArticle {
	article := tgbotapi.NewInlineQueryResultArticle(
		cast.ToString(rand.Intn(1000000)),
		"Выпустить пососямбу в этот чат",
		preparedPhrases[rand.Intn(len(preparedPhrases))]+"\r\n\n"+sb.BuildPososyamba(),
	)
	return article
}

func InlineFakeNews() tgbotapi.InlineQueryResultArticle {
	title, _ := fakenews.FetchTitle()
	fakeNews := tgbotapi.NewInlineQueryResultArticle(
		cast.ToString(rand.Intn(1000000)),
		"Сгенерить фейкньюс",
		title,
	)
	return fakeNews
}

func InlineF() (tgbotapi.InlineQueryResultGIF, error) {
	url, preview, err := tenor.GetGifsByIDs(os.Getenv("FUNERAL_GIFS_IDS"))
	if err != nil {
		return tgbotapi.InlineQueryResultGIF{}, err
	}
	gif := tgbotapi.NewInlineQueryResultGIF(os.Getenv("FUNERAL_GIFS_IDS"), url)
	gif.ThumbURL = preview
	gif.Title = "F"

	return gif, err
}
