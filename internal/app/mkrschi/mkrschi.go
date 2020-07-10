package mkrschi

import (
    "github.com/rs/zerolog/log"
    "gopkg.in/yaml.v2"

    "io/ioutil"
)

func UpdatePhrases(mkrschi_phrase string) {

    file, err := ioutil.ReadFile("configs/mkrschi_phrases.yml")
    if err != nil {
        log.Error().Err(err).Msg("")
    }
    var mkrschi_phrases []string
    err = yaml.Unmarshal(file, &mkrschi_phrases)
    if err != nil {
        log.Error().Err(err).Msg("")
    }

    mkrschi_phrases = append(mkrschi_phrases, mkrschi_phrase)

    var bytes []byte

    bytes, err = yaml.Marshal(mkrschi_phrases)
    if err != nil {
        log.Error().Err(err).Msg("")
    }
    err = ioutil.WriteFile("configs/mkrschi_phrases.yml", bytes, 0644)
    if err != nil {
        log.Error().Err(err).Msg("")
    }

}
