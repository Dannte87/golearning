package sites

import (
	"parser/linkGenerator"
	"parser/entities"
	"parser/config"
	"github.com/opesun/goquery"
	"regexp"
	"strconv"
	"strings"
	"fmt"
)

//Entity RealEstate https://www.real-estate.lviv.ua/
type RealEstate struct {
	Links linkGenerator.GetMethod
	entities.FlatEntity
}


func (realEstate RealEstate) GetData(conf *config.Config) entities.FlatEntity {
	var g linkGenerator.Link

	g = &realEstate.Links

	links := g.Generate()

	if links != nil {
		for _, l := range links {
			realEstate.Fields.Link = conf.GetLink(l)

			var e entities.FlatParser

			e = realEstate

			go e.Parser(conf.GetContent())
		}
	}

	return realEstate.FlatEntity
}

func (realEstate RealEstate) Parser(c config.Content) entities.FlatEntity  {
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