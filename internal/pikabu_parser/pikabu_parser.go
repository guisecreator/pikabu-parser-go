package pikabu_parser

import (
	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/guisecreator/pikabu-parser-go/parser"
)

type ParserPikabu struct {
	Url             string
	missTags 		[]string
	AppINFO 		interface{}
	DataBase 		interface{}
	ID				interface{}
	Posts			[]map[string]interface{}

}

func NewParserPikabu(Url string, missTags []string, AppINFO interface{}, DataBase interface{}, ID interface{}) *ParserPikabu {
	return &ParserPikabu{
		Url:            Url,
		missTags:	    missTags,
		AppINFO:		AppINFO,
		DataBase: 		DataBase,
		ID: 			ID,		
	}
}

func (pars *ParserPikabu) GetPosts() []map[string]interface{} {
	// Posts := []map[string]interface{}{}

	doc, err := goquery.NewDocument(pars.Url)
	if err != nil {
		log.Println(err)
		return pars.Posts
	}

	doc.Find(".story").Each(func(i int, s *goquery.Selection) {
		postID := parser.GetArticleID(s) 
		if pars.existPost(postID){
			return
		}

		postLink := pars.getPostLink(s)
		postDoc, err := goquery.NewDocument(postLink)
		if err != nil {
			log.Println(err)
			return
		}

		if pars.ignorePost(postDoc.Selection){
			return
		}

		PostDate := pars.getPostDate(postDoc.Selection)
		PostTitle := pars.getPostTitle(postDoc.Selection)
		PostTags := pars.getPostTags(postDoc.Selection)

		post := map[string]interface{}{
			"PostDate": PostDate,
			"PostTitle": PostTitle,
			"PostTags": PostTags,

		}
		
		pars.Posts = append(pars.Posts, post)
	})
	
	return pars.Posts 
}

func (pars *ParserPikabu) getPostLink(link *goquery.Selection) string {
	postlink, _ := link.Find(".story__title-link").Attr("href")
	return pars.Url + postlink
}

func (pars *ParserPikabu) ignorePost(ignr *goquery.Selection) bool {
	return false
}

func (pars *ParserPikabu) getPostDate(date *goquery.Selection) string {
	postdate := date.Find(".story__datetime").Text()
	return postdate	 
}

func (pars *ParserPikabu) getPostTitle(title *goquery.Selection) string{
	postTitle := title.Find(".story__datetime").Text()
	return postTitle
}

func (pars *ParserPikabu) getPostTags(tag *goquery.Selection) []string {
	tags := []string{}
	tag.Find(".story__tags .tags__tag").Each(func(i int, tag *goquery.Selection) {
		tags = append(tags, tag.Text())
	})
	
	return tags
}

func (pars *ParserPikabu) existPost(postID string) bool {
	return false
}

