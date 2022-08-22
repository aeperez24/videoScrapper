package service

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	config, err := LoadConfig("inputs")
	log.Println(err)
	assert.Nil(t, err, "Error parsing")
	SerieConfigurations := []SerieConfiguration{
		{SerieLink: "link1", SerieName: "name1", Provider: "provider"}, {SerieLink: "link2", SerieName: "name2", Provider: "provider"}}
	expected := AppConfiguration{SerieConfigurations: SerieConfigurations, OutputPath: "/output/"}
	assert.Equal(t, expected, config, "config is not equals to expected")
}
