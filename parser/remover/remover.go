package remover

import (
	"strings"

	"golang.org/x/net/html"
	"github.com/PuerkitoBio/goquery"
)

func removeArticle(articleTree *goquery.Selection, remClass string) {
    articleTree.Find("*").Each(func(i int, s *goquery.Selection) {
        class, exists := s.Attr("class")
        if exists {
            if strings.Contains(class, remClass) {
                s.Remove()
            }
        }
    })
}

func removeClassTree(articleTree []*html.Node, remClass string) {
	doc := goquery.NewDocumentFromNode(articleTree[0])

	doc.Find("." + remClass).Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})
}