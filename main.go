package main

import (
	"github.com/opesun/goquery"
	"parser/config"
	"fmt"
	"strconv"
	"strings"
	"net/url"
	"bytes"
	"net/http"
	"io/ioutil"
	"regexp"
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
* General interface that describe logic for getting links on detail page of flat
*/
type Link interface {
	Generate() map[int]string
}

/**
* Interface that describe logic for parsing main data of flat
*/
type FlatParser interface {
	Parser(c config.Content) FlatEntity
	GetData(conf *config.Config) FlatEntity
}

/**
* Options for getting detail links from some source
*/
type source struct {
	sourceUri string
	detailHtmlSelector string
}

/**
* Entity OLX https://www.olx.ua/
*/
type Olx struct {
	source
	FlatEntity
}

/**
* Entity RealEstate https://www.real-estate.lviv.ua/
*/
type RealEstate struct {
	source
	FlatEntity
}

func (olx Olx) Generate() map[int]string {
	x, err := goquery.ParseUrl(olx.sourceUri)

	links := make(map[int]string)

	if err == nil {
		for i, v := range x.Find(olx.detailHtmlSelector).Attrs("href") {
			links[i] = v
		}
	}

	return links
}

func (realEstate RealEstate) Generate() map[int]string  {
	requestUrl := realEstate.sourceUri

	form := url.Values{
		"hash": {requestUrl},
	}

	body := bytes.NewBufferString(form.Encode())

	rsp, err := http.Post(requestUrl, "application/x-www-form-urlencoded", body)

	if err != nil {
		panic(err)
	}

	defer rsp.Body.Close()

	bodyByte, err := ioutil.ReadAll(rsp.Body)

	if err != nil {
		panic(err)
	}

	x, err := goquery.ParseString(string(bodyByte))
	links := make(map[int]string)

	if err == nil {
		for i, v := range x.Find(realEstate.detailHtmlSelector).Attrs("href") {
			links[i] = v
		}
	}

	return links
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
	var g Link

	g = olx

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
	var g Link

	g = realEstate

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

			olx.sourceUri = conf.Link
			olx.detailHtmlSelector = conf.Selector
			olx.Fields.Name = conf.Name

			d = olx

			break
		case "real-estate":
			re := RealEstate{}

			re.sourceUri = conf.Link
			re.detailHtmlSelector = conf.Selector
			re.Fields.Name = conf.Name

			d = re

			break
		}

		go d.GetData(&conf)
	}


	var input string
	fmt.Scanln(&input)
}

