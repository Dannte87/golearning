package main

import (
	"github.com/opesun/goquery"
	"parser/config"
	"fmt"
	"strconv"
	"strings"
	"regexp"
	"parser/linkGenerator"
)

var (
	settings = map[string]string{
		"olx" : "/home/dante/GoProjects/bin/olx.yaml",
		"real-estate" : "/home/dante/GoProjects/bin/real-estate.yaml",
	}
)
/**
* Fields for describe flay entity
*/
type Flat struct {
	Name string
	Link string
	Price string
	Rooms int
	Date string
	Headline string
	District string
}

/**
* Html selectors for parsing content
*/
type selectors struct {
	Link string
	Price string
	Rooms string
	Date string
	Headline string
	District string
}

/**
* Main structure that content all data for working with flat entity
*/
type FlatEntity struct {
	Fields Flat
	HtmlTags selectors
}

/**
* Interface that describe logic for parsing main data of flat
*/
type FlatParser interface {
	Parser(c config.Content) FlatEntity
	GetData(conf *config.Config) FlatEntity
}

/**
* Entity OLX https://www.olx.ua/
*/
type Olx struct {
	links linkGenerator.GetMethod
	FlatEntity
}

/**
* Entity RealEstate https://www.real-estate.lviv.ua/
*/
type RealEstate struct {
	links linkGenerator.PostMethod
	FlatEntity
}

func (olx Olx) Parser(c config.Content) FlatEntity {
	x, err := goquery.ParseUrl(olx.Fields.Link)

	if err == nil {
		olx.Fields.Price = x.Find(c.Price).Text()

		html := x.Find(c.Rooms)
		r, err := strconv.ParseInt(strings.TrimSpace(html.HtmlAll()[2]), 10, 32)

		if err == nil {
			olx.Fields.Rooms = int(r)
		}

		text := strings.TrimSpace(x.Find(c.Date).Text())
		olx.Fields.Date = strings.TrimSpace(strings.Split(text,",")[1])

		olx.Fields.Headline = strings.TrimSpace(x.Find(c.Headline).Text())

		district := strings.Split(x.Find(c.District).Text(), ",")
		olx.Fields.District = district[2]
	}

	fmt.Println(olx.Fields)

	return olx.FlatEntity
}

func (olx Olx) GetData(conf *config.Config) FlatEntity  {
	var g linkGenerator.Link

	g = &olx.links

	links := g.Generate()


	if links != nil {
		for _, l := range links {
			olx.Fields.Link = conf.GetLink(l)

			var e FlatParser

			e = olx

			go e.Parser(conf.GetContent())
		}
	}

	return olx.FlatEntity
}

func (realEstate RealEstate) GetData(conf *config.Config) FlatEntity  {
	var g linkGenerator.Link

	g = &realEstate.links

	links := g.Generate()

	if links != nil {
		for _, l := range links {
			realEstate.Fields.Link = conf.GetLink(l)

			var e FlatParser

			e = realEstate

			go e.Parser(conf.GetContent())
		}
	}

	return realEstate.FlatEntity
}

func (realEstate RealEstate) Parser(c config.Content) FlatEntity  {
	x, err := goquery.ParseUrl(realEstate.Fields.Link)

	if err == nil {
		realEstate.Fields.Price = x.Find(c.Price).Html()

		r := x.Find(c.Rooms).HtmlAll()
		re := regexp.MustCompile("[0-9]+")
		rooms, _ := strconv.ParseInt(re.FindAllString(r[1], -1)[0], 10, 32)
		realEstate.Fields.Rooms = int(rooms)

		dateText := strings.Split(r[len(r) - 2], ":")
		realEstate.Fields.Date = strings.TrimSpace(dateText[1])

		realEstate.Fields.Headline = x.Find(c.Headline).Text()

		district := strings.Split(x.Find(c.District).HtmlAll()[1], " ")
		realEstate.Fields.District = district[0]
	}

	fmt.Println(realEstate.Fields)

	return realEstate.FlatEntity
}

func main() {

	for k, s := range settings {
		var d FlatParser
		conf := config.Init(s)

		switch k {
		case "olx":
			olx := Olx{}

			olx.links.Link = conf.Link
			olx.links.Selector = conf.Selector
			olx.Fields.Name = conf.Name

			d = olx

			break
		case "real-estate":
			re := RealEstate{}

			re.links.Link = conf.Link
			re.links.Selector = conf.Selector
			re.Fields.Name = conf.Name

			d = re

			break
		}

		go d.GetData(&conf)
	}


	var input string
	fmt.Scanln(&input)
}

