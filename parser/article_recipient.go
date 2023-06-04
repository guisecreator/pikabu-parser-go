package parser

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/guisecreator/pikabu-parser-go/parser/recipient"
)

func (p *Parser) GetListArticles() []map[string]string {
	listArticlesID := make([]map[string]string, 0)
	if p.getTree() {
		if recipient.IsBlocks() {
			for _, block := range recipient.GetBlocks() {
				if !p.notMissingArticle(block) {
					continue
				}

				articleID := GetArticleID(block)
				if articleID == "" {
					continue
				}

				link := p.normalizeURL(p.GetArticleLink(block))
				listArticlesID = append(
					listArticlesID, map[string]string{
						"Id": articleID, 
						"Link": link,
					})
			}
		}
	}
	return listArticlesID
}

func (p *Parser) GetArticleTags(articleTree *goquery.Selection) []string {
	var tags []string
	tagsNodes := articleTree.Find(".tag-list").Find("a")

	tagsNodes.Each(func(_ int, 
		tagNode *goquery.Selection) {
		tagText := strings.TrimSpace(tagNode.Text())

		if len(tagText) > 0 {
			tags = append(tags, tagText)
		}
	})
	return tags

}

func (p *Parser) GetArticle(articleLink string) *goquery.Document {
	data, err := p.HTTPClient.Get(articleLink)
	if err != nil {
		recipient.DbLog(fmt.Sprintf("Error: %v", err))
		return nil
	}

	defer data.Body.Close()

	doc, err := goquery.NewDocumentFromReader(data.Body)
	if err != nil {
		recipient.DbLog(fmt.Sprintf("Error: %v", err))
		return nil
	}

	return doc

}

func GetArticleID(blockTree *goquery.Selection) string {
	return ""
}

func  (p *Parser) GetArticleLink(blockTree *goquery.Selection) string {
	//todo implement me
	// Returns a link to the article
	panic("not implemented")
}

func (p *Parser) GetArticleTitle(articleTree *goquery.Document) string {
	//todo implement me
	// Returns the title of the article
	panic("not implemented")
}

func (p *Parser) getArticleDate(articleTree *goquery.Document) time.Time {
	// Returns the date of the article
	return time.Now()
}

func (p *Parser) notMissingArticle(block *goquery.Selection) bool {
	// Add conditions here
	return true
}

func (p *Parser) IgnoreArticle(articleTree *goquery.Document) bool {
	// Ignore article
	return false
}



