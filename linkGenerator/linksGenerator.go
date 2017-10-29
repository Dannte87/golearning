package linkGenerator

import (
	"github.com/opesun/goquery"
	"net/url"
	"bytes"
	"net/http"
	"io/ioutil"
)

var (
	links map[int]string
)

//General interface that describe logic for getting links on detail page of flat
type Link interface {
	Generate() map[int]string
}

//Structure for generate links via method GET
type GetMethod struct {
	Link string
	Selector string
}

//Structure for generate links via method POST
type PostMethod struct {
	Link string
	Selector string
}

//Make map of links and return it
func build(x []string) map[int]string {
	links := make(map[int]string)

	for i, v := range x {
		links[i] = v
	}

	return links
}

func (get *GetMethod) Generate() map[int]string {
	x, err := goquery.ParseUrl(get.Link)

	if err == nil {
		links = build(x.Find(get.Selector).Attrs("href"))
	}

	return links
}

func (post *PostMethod) Generate() map[int]string {
	form := url.Values{
		"hash": {post.Link},
	}

	body := bytes.NewBufferString(form.Encode())

	rsp, err := http.Post(post.Link, "application/x-www-form-urlencoded", body)

	if err != nil {
		panic(err)
	}

	defer rsp.Body.Close()

	bodyByte, err := ioutil.ReadAll(rsp.Body)

	if err != nil {
		panic(err)
	}

	x, err := goquery.ParseString(string(bodyByte))

	if err == nil {
		links = build(x.Find(post.Selector).Attrs("href"))
	}

	return links
}


