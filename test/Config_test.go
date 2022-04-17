package test

import (
	"aeperez24/animewatcher/service"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	config, err := service.LoadConfig("inputs")
	log.Println(err)
	assert.Nil(t, err, "Error parsing")
	animeConfigurations := []service.AnimeConfiguration{
		{AnimeLink: "link1", AnimeName: "name1"}, {AnimeLink: "link2", AnimeName: "name2"}}
	expected := service.AppConfiguration{AnimeConfigurations: animeConfigurations, OutputPath: "output"}
	assert.Equal(t, expected, config, "config is not equals to expected")
}
