package fakenews

import (
	"encoding/json"
	"github.com/mb-14/gomarkov"
	"github.com/rs/zerolog/log"
	"github.com/thesunwave/pososyamba_bot/internal/app/cache"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type mMTitle struct {
	Title string
}

func generateNews() (*[]string, error) {
	//Create a chain of order 2
	chain, err := loadModel()

	if err != nil {
		chain, err = BuildModel()
	}

	var newsList []string

	for i := 0; i < 2000; i++ {
		newsList = append(newsList, generateTitle(chain))
	}

	result := removeDuplicatesUnordered(newsList)

	err = saveToRedis(&result)
	if err != nil {
		log.Error().Err(err)
	}

	return &result, err
}

func BuildModel() (*gomarkov.Chain, error) {
	chain := gomarkov.NewChain(1)
	titles := fetchTitles()
	for _, story := range titles {
		chain.Add(strings.Split(story.Title, " "))
	}

	jsonObj, _ := json.Marshal(chain)
	err := ioutil.WriteFile("model.json", jsonObj, 0644)

	if err != nil {
		log.Error().Err(err)
	}

	return chain, err
}

func saveToRedis(titles *[]string) error {
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

		result, err := generateNews()

		log.Print(result)
		if err != nil {
			log.Error().Err(err)
		}

		redisObj = cache.Redis().SPop("news_titles")
	}

	return redisObj.Val(), err
}

func removeDuplicatesUnordered(elements []string) []string {
	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key, _ := range encountered {
		result = append(result, key)
	}
	return result
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

func fetchTitles() []mMTitle {
	var titles []mMTitle
	var i interface{}

	resp, err := http.Get("https://meduza.io/api/v3/search?chrono=news&locale=ru&page=0&per_page=150")
	if err != nil {
		log.Error().Err(err)
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Error().Err(err)
	}

	err = json.Unmarshal(body, &i)

	if err != nil {
		log.Error().Err(err)
	}

	m := i.(map[string]interface{})
	z := m["documents"].(map[string]interface{})

	for _, v := range z {
		titles = append(titles, mMTitle{Title: v.(map[string]interface{})["title"].(string)})
	}

	return titles
}
