package main

import (
	"parser/config"
	"fmt"
	"parser/sites"
	"parser/entities"
)

var (
	settings = map[string]string {
		"olx" : "/home/dante/GoProjects/bin/olx.yaml",
		"real-estate" : "/home/dante/GoProjects/bin/real-estate.yaml",
	}
)

func main() {
	for k, s := range settings {
		var d entities.FlatParser
		conf := config.Init(s)

		switch k {
		case "olx":
			olx := sites.Olx{}

			olx.Links.Link = conf.Link
			olx.Links.Selector = conf.Selector
			olx.Fields.Name = conf.Name

			d = olx

			break
		case "real-estate":
			re := sites.RealEstate{}
			re.Links.Link = conf.Link
			re.Links.Selector = conf.Selector
			re.Fields.Name = conf.Name

			d = re

			break
		}

		go d.GetData(&conf)
	}


	var input string
	fmt.Scanln(&input)
}

