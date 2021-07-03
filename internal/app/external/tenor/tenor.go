package tenor

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func GetGifsByIDs(query string) (string, error) {
	u, err := url.Parse("https://g.tenor.com/v1/gifs")
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("ids", query)
	q.Set("key", os.Getenv("TENOR_API_KEY"))
	q.Set("media_filter", "minimal")
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	jsonDataFromHttp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response Response
	err = json.Unmarshal(jsonDataFromHttp, &response)
	if err != nil {
		return "", err
	}

	media := response.Results[0].Media
	return media[0].Gif.Url, nil
}
