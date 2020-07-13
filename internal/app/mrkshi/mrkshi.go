package mrkshi

import (
    "github.com/rs/zerolog/log"
    "gopkg.in/yaml.v2"

    "io/ioutil"
)

func UpdatePhrases(mrkshi_phrase string, mrkshi_phrases *[]string) {    
    *mrkshi_phrases = append(*mrkshi_phrases, mrkshi_phrase)

    bytes, err := yaml.Marshal(*mrkshi_phrases)
    if err != nil {
        log.Error().Err(err).Msg("")
    }
    
    err = ioutil.WriteFile("configs/mrkshi_phrases.yml", bytes, 0644)
    if err != nil {
        log.Error().Err(err).Msg("")
    }

}
