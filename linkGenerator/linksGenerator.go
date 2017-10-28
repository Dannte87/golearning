package linkGenerator

import (
	"github.com/opesun/goquery"
)

func Generate(link string, selector string) map[int]string {
	x, err := goquery.ParseUrl(link)

	links := make(map[int]string)

	if err == nil {
		for i, v := range x.Find(selector).Attrs("href") {
			links[i] = v
		}
	}

	return links
}


