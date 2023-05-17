package pikabu_parser

import (
	"fmt"
	
)

type ParserPikabu struct {
	Url             string
	APPID           string
	DB              string
	appInfo         string
	missTags        []string
	onlyWithMedia   bool
}

func NewParserPikabu(Url string, missTags []string, appInfo string, appid string, db string, onlyWithMedia bool) *ParserPikabu {
	return &ParserPikabu{
		Url:             Url,
		missTags:        missTags,
		APPID:           appid,
		DB:              db,
		appInfo:         appInfo,
		onlyWithMedia:   onlyWithMedia,
	}
}

func (p *ParserPikabu) parse() {
	fmt.Println("Parse Pikabu")
}