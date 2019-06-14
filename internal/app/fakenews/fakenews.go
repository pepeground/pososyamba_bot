package fakenews

import (
	"encoding/json"
	"github.com/mb-14/gomarkov"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
	"github.com/thesunwave/pososyamba_bot/internal/app/cache"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	Documents map[string]Document `json:documents`
}

type Document struct {
	Title string `json:title`
}

func generateNews() (*[]string, error) {
	//Create a chain of order 2
	chain, err := loadModel()

	var titles *[]string

	if err != nil {
		chain, titles, err = BuildModel()
	}

	var newsList []string

	for i := 0; i < 1000; i++ {
		newsList = append(newsList, generateTitle(chain))
	}

	log.Print(newsList)
	err = saveToRedis(&newsList)

	cache.Redis().SRem("news_titles", *titles)

	if err != nil {
		log.Error().Err(err)
	}

	return &newsList, err
}

func BuildModel() (*gomarkov.Chain, *[]string, error) {
	chain := gomarkov.NewChain(1)
	titles := fetchTitles()
	for _, story := range *titles {
		chain.Add(strings.Split(story, " "))
	}

	jsonObj, _ := json.Marshal(chain)
	err := ioutil.WriteFile("model.json", jsonObj, 0644)

	if err != nil {
		log.Error().Err(err)
	}

	return chain, titles, err
}

func saveToRedis(titles *[]string) error {
	log.Print(len(*titles))
	return cache.Redis().SAdd("news_titles", *titles).Err()
}

func FetchTitle() (string, error) {
	redisObj := cache.Redis().SPop("news_titles")
	var err error

	log.Print(redisObj.Val())

	if redisObj.Err() != nil {
		log.Error().Err(redisObj.Err())

		err = os.Remove("model.json")
		if err != nil {
			log.Error().Err(err)
		}

		_, err := generateNews()

		if err != nil {
			log.Error().Err(err)
		}

		redisObj = cache.Redis().SPop("news_titles")
	}

	return redisObj.Val(), err
}

func loadModel() (*gomarkov.Chain, error) {
	var chain gomarkov.Chain
	data, err := ioutil.ReadFile("model.json")
	if err != nil {
		return &chain, err
	}
	err = json.Unmarshal(data, &chain)
	if err != nil {
		return &chain, err
	}
	return &chain, nil
}

func generateTitle(chain *gomarkov.Chain) string {
	tokens := []string{gomarkov.StartToken}
	for tokens[len(tokens)-1] != gomarkov.EndToken {
		next, _ := chain.Generate(tokens[(len(tokens) - 1):])
		tokens = append(tokens, next)
	}

	return strings.Join(tokens[1:len(tokens)-1], " ")
}

func fetchTitles() *[]string {
	var titles []string
	var response Response

	for i := 0; i < 10; i++ {
		resp, err := http.Get("https://meduza.io/api/v3/search?chrono=news&locale=ru&page=" + cast.ToString(i) + "&per_page=50")
		if err != nil {
			log.Error().Err(err)
		}

		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			log.Error().Err(err)
		}

		err = json.Unmarshal(body, &response)

		for _, v := range response.Documents {
			log.Print(v.Title)
			titles = append(titles, v.Title)
		}

		if err != nil {
			log.Error().Err(err)
		}
	}

	return &titles
}
