package entities

type FlatParser interface {
	Parser() FlatEntity
}

type Flat struct {
	Name string
	Link string
	Price string
	Rooms int
	Date string
	Headline string
	District string
}

type selectors struct {
	Price string
	Rooms string
	Date string
	Headline string
	District string
}

type FlatEntity struct {
	Fields Flat
	HtmlTags selectors
}

