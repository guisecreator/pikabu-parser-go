package recipient

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func getListArticles() []map[string]string {
	listArticlesID := make([]map[string]string, 0)
	if getTree() {
		if isBlocks() {
			for _, block := range getBlocks() {
				if !notMissingArticle(block) {
					continue
				}

				articleID := getArticleID(block)
				if articleID == "" {
					continue
				}

				link := normalizeURL(getArticleLink(block))
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

func getArticleTags(articleTree *goquery.Selection) []string {
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

func getArticle(articleLink string) *goquery.Document {
	data, err := parser.HTTPClient.Get(articleLink)
	if err != nil {
		dbLog(fmt.Sprintf("Error: %v", err))
		return nil
	}

	defer data.Body.Close()

	doc, err := goquery.NewDocumentFromReader(data.Body)
	if err != nil {
		dbLog(fmt.Sprintf("Error: %v", err))
		return nil
	}

	return doc

}

func normalizeURL(url string) string {
	psex := parser.ParsRegularExp
	if !strings.Contains(url, parser.BaseURL) && !strings.Contains(url, "http") {
		if url[:1] == "/" && url[:2] != "//" {
			return parser.BaseURL + url
		} else if url[:2] == "//" {
			return psex + url
		} else {
			return parser.BaseURL + "/" + url
		}
	}
	return url
}

func getArticleID(blockTree *goquery.Selection) string {
	// return ID Art
	return ""
}

func  getArticleLink(blockTree *goquery.Selection) string {
	//todo implement me
	// Returns a link to the article
	panic("not implemented")
}

func  getArticleTitle(articleTree *goquery.Document) string {
	//todo implement me
	// Returns the title of the article
	panic("not implemented")
}

func  getArticleDate(articleTree *goquery.Document) time.Time {
	// Returns the date of the article
	return time.Now()
}

func notMissingArticle(block *goquery.Selection) bool {
	// Add conditions here
	return true
}

func ignoreArticle(articleTree *goquery.Document) bool {
	// Ignore article
	return false
}



