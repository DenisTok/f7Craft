package main

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
)

type Emj struct {
	Emoji string `json:"emoji"`
}

func main() {
	b, err := ioutil.ReadFile("emoji.json")
	if err != nil {
		log.Fatal().Err(err).Send()
		return
	}

	var Emjs []*Emj

	err = json.Unmarshal(b, &Emjs)
	if err != nil {
		log.Fatal().Err(err).Send()
		return
	}

	var s []string

	for i := 0; i < 256; i++ {
		s = append(s, fmt.Sprintf("\"%s\",", Emjs[i].Emoji))
	}

	err = ioutil.WriteFile("emoji.gen", []byte(fmt.Sprintf("%v", s)), 0777)
	if err != nil {
		log.Fatal().Err(err).Send()
		return
	}
}
