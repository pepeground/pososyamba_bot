package fakenews

import (
	"encoding/json"
	"github.com/mb-14/gomarkov"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"strings"
)

type mMTitle struct {
	Title string
}

func GenerateNews() (string, error) {
	//Create a chain of order 2
	chain, err := loadModel()

	if err != nil {
		chain, err = BuildModel()
	}

	return generateTitle(chain), err
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
