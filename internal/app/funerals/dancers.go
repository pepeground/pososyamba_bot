package funerals

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

func DancersFile() ([]byte, string, error) {
	fileInfo, err := ioutil.ReadDir("assets/dancers")
	if err != nil {
		log.Error().Err(err)
		return []byte{}, "", err
	}

	var onlyImages []os.FileInfo

	for _, f := range fileInfo {
		if strings.Split(f.Name(), ".")[1] == "gif" {
			onlyImages = append(onlyImages, f)
		}
	}

	if len(onlyImages) == 0 {
		return []byte{}, "", errors.New("no images found")
	}

	rand.Seed(time.Now().UnixNano())
	fileName := onlyImages[rand.Intn(len(onlyImages))]

	gif, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", "assets/dancers", fileName.Name()))
	if err != nil {
		log.Error().Err(err).Msg("")
		return []byte{}, "", err
	}

	return gif, fileName.Name(), nil
}
