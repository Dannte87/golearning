package sites

import (
	"parser/entities"
	"parser/linkGenerator"
	"parser/config"
	"github.com/opesun/goquery"
	"strconv"
	"strings"
	"fmt"
)

//Entity OLX https://www.olx.ua/
type Olx struct {
	Links linkGenerator.GetMethod
	entities.FlatEntity
}

func (olx Olx) Parser(c config.Content) entities.FlatEntity {
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

func (olx Olx) GetData(conf *config.Config) entities.FlatEntity {
	var g linkGenerator.Link

	g = &olx.Links

	links := g.Generate()


	if links != nil {
		for _, l := range links {
			olx.Fields.Link = conf.GetLink(l)

			var e entities.FlatParser

			e = olx

			go e.Parser(conf.GetContent())
		}
	}

	return olx.FlatEntity
}