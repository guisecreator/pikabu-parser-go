package main

import (
	"fmt"
	"time"

	"github.com/guisecreator/pikabu-parser-go/internal/pikabu_parser"
)

func main() {
	Url := "https://pikabu.ru"

	AppInfo := map[string]interface{}{
		"Url": Url,
	}
														  //тэги pikabu
	parser := pikabu_parser.NewParserPikabu(Url, []string{"#", "#", "#",}, AppInfo, nil, nil,)
	fmt.Printf("Parsing has begun...")

	posts := parser.GetPosts()

	for _, post := range posts{
		fmt.Println(post)
	}

	fmt.Printf("%s: Parsing completed\n", time.Now().Format("2005-04-02 12:12:12"))

}


