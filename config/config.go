package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Detail struct {
	Link string
	Selector string
	Domain string
}

type Content struct {
	Price string
	Rooms string
	Date string
	Headline string
	District string
}

type Config struct {
	Name string
	Detail
	Content
}

func (c *Config) GetName() string {
	return c.Name
}

func (c *Config) GetDetail() Detail {
	return c.Detail
}

func (c *Config) GetContent() Content {
	return c.Content
}

func (c *Config) GetLink(l string) string {
	if c.Domain != "" {
		return c.Domain + l
	}

	return l
}

func Init(fn string) Config {
	var config Config
	source, err := ioutil.ReadFile(fn)
	
	if err != nil {
		panic(err)
	}
	
	err = yaml.Unmarshal(source, &config)
	
	if err != nil {
		panic(err)
	}

	return config
}

