package entities

import "parser/config"

//Fields for describe flay entity
type Flat struct {
	Name string
	Link string
	Price string
	Rooms int
	Date string
	Headline string
	District string
}

//Interface that describe logic for parsing main data of flat
type FlatParser interface {
	Parser(c config.Content) FlatEntity
	GetData(conf *config.Config) FlatEntity
}



//Main structure that content all data for working with flat entity
type FlatEntity struct {
	Fields Flat
}

