package parser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type Parser struct {
	EntryURL     string
	BaseURL      string
	Pre          string
	MissTags     []string
	PublicTags   bool
	
    EntryTree    *goquery.Document
    HTTPClient   *http.Client
}

func NewParser(url string, 
	missTags []string,
	publicTags bool) *Parser {

	re := regexp.MustCompile(`((http[s]?)://[^/]+)`)

	BaseURL := re.FindStringSubmatch(url)[1]

	Pre := re.FindStringSubmatch(url)[2]
	return &Parser{
		EntryURL:   url,
		BaseURL:    BaseURL,
		Pre:        Pre,
		MissTags:   missTags,
		PublicTags: publicTags,
	}
}

func (p *Parser) GetArticles() []map[string]interface{} {
	listArticles := p.getListArticles()
	// reverse list of articles
	for i, j := 0, len(listArticles)-1; i < j; i, j = i+1, j-1 {
		listArticles[i], listArticles[j] = listArticles[j], listArticles[i]
	}

	// articles := []map[string]interface{}{}

	for _, article := range listArticles {
		articleLink := strings.Replace(fmt.Sprintf(
			"%v", article["Link"]), 
			"old_string", "new_string", -1)

		errorText  :=    []string{}
		formatedText :=  []string{}
		articleTags :=   []string{}

		public := true

		fmt.Printf("Parsing: %s\n", articleLink)

		articleTree, err := p.getListArticles(articleLink)
		if err != nil {
			fmt.Printf("Error in ID: %s\n", articleLink)
			continue
		}

		p.timer(0)

		// articleDate := p.getArticleDate(articleTree)

		articleTitle := p.getArticleTitle(articleTree)
		if articleTitle == "" {
			fmt.Printf("Error in title: %s\n", articleLink)
			continue
		}

		articleTags = p.getArticleTags(articleTree)
		public = p.missToTags(articleTags)
		if !public {
			errorText = append(errorText, fmt.Sprintf("There is a tag from the list'%s' ", strings.Join(p.MissTags, ",")))
		}

		// if p.publicTags && len(articleTags) > 0 {
		// 	formatedText = append(formatedText, "\r\n\r\n", strings.Join(articleTags, " "))
		// }

		if len(strings.Join(formatedText, "")) >= 7000 {
			public = false
			fmt.Printf("The text is long, post id: %d\n", article["Id"])
			errorText = append(errorText, fmt.Sprintf("The text is long, post id: %d", article["Id"]))
		}

		// if p.igonreArticle(articleTree) {
		// 	fmt.Printf("Ignore, post id: %d\n", article["Id"])
		// 	formatedArticles = append(formatedArticles, map[string]interface{}{
		// 		"Id":        article["Id"], 
		// 		"Link":      articleLink,
		// 	})
		// }
	}
	return nil
}

func (p *Parser) getListArticles() []map[string]string {
	listArticlesID := make([]map[string]string, 0)
	if p.getTree() {
		if p.isBlocks() {
			for _, block := range p.getBlocks() {
				if !p.notMissingArticle(block) {
					continue
				}
				articleID := p.getArticleID(block)
				if articleID == "" {
					continue
				}
				link := p.normalizeURL(p.getArticleLink(block))
				listArticlesID = append(listArticlesID, map[string]string{"Id": articleID, "Link": link})
			}
		}
	}
	return listArticlesID
}

func (p *Parser) normalizeURL(url string) string {
	pre := p.Pre
	if !strings.Contains(url, p.BaseURL) && !strings.Contains(url, "http") {
		if url[:1] == "/" && url[:2] != "//" {
			return p.BaseURL + url
		} else if url[:2] == "//" {
			return pre + url
		} else {
			return p.BaseURL + "/" + url
		}
	}
	return url
}

func (p *Parser) excludePosts(listArticlesID []map[string]string) {
	toRem := []map[string]string{}
	for _, articleID := range listArticlesID {
		if p.dbExistArticle(articleID["Id"]) {
			toRem = append(toRem, articleID)
			p.dbLog(fmt.Sprintf("Has already: %v", articleID["Id"]))
		}
	}
}

func removeArticle(articleTree *goquery.Selection, remclass string) {
    articleTree.Find("*").Each(func(i int, s *goquery.Selection) {
        class, exists := s.Attr("class")
        if exists {
            if strings.Contains(class, remclass) {
                s.Remove()
            }
        }
    })
}

// func (p *Parser) getTree() bool {
//     data, err := requests.Get(p.EntryURL, headers)
//     if err != nil {
//         p.dbLog(fmt.Sprintf("Ошибка: %v", err))
//         return false
//     }
//     p.EntryURL = html.FromString(data.Content())
//     return true
// }

func (p *Parser) getTree() bool {
    resp, err := p.HTTPClient.Get(p.EntryURL)
    if err != nil {
        log.Printf("Failed to get page: %v", err)
        return false
    }
    defer resp.Body.Close()

    pageBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Failed to read response body: %v", err)
        return false
    }

    pageReader := bytes.NewReader(pageBytes)
    doc, err := goquery.NewDocumentFromReader(pageReader)
    if err != nil {
        log.Printf("Failed to parse HTML: %v", err)
        return false
    }

    p.EntryTree = doc
    return true
}

func (p *Parser) getArticleTags(articleTree *html.Node) []string {
	var tags []string
	tagsNodes := articleTree.Find(".tag-list").Find("a")
	tagsNodes.Each(func(_ int, tagNode *goquery.Selection) {
		tagText := strings.TrimSpace(tagNode.Text())
		if len(tagText) > 0 {
			tags = append(tags, tagText)
		}
	})
	return tags

}

func (p *Parser) getArticle(articleLink string) *goquery.Document {
	// Returns the article object
	data, err := p.HTTPClient.Get(articleLink)
	if err != nil {
		p.dbLog(fmt.Sprintf("Ошибка: %v", err))
		return nil
	}
	defer data.Body.Close()

	doc, err := goquery.NewDocumentFromReader(data.Body)
	if err != nil {
		p.dbLog(fmt.Sprintf("Error: %v", err))
		return nil
	}
	return doc
}

func (p *Parser) missToTags(articleTags []string) bool {
    if len(p.MissTags) > 0 {
        for _, atags := range articleTags {
            for _, mtags := range p.MissTags {
                if strings.Contains(strings.ToLower(atags), strings.ToLower(mtags)) {
                    fmt.Sprintf("Не публикуем, присутствует тэг \"%s\"", mtags)
                    return false
                }
            }
        }
    }
    return true
}

func (p *Parser) removeClassTree(articleTree []*html.Node, remClass string) {
	doc := goquery.NewDocumentFromNode(articleTree[0])

	doc.Find("." + remClass).Each(func(i int, 
		s *goquery.Selection) {
		s.Remove()
	})
}

func (p *Parser) getArticleID(blockTree *goquery.Selection) string {
	// return ID Art
    return ""
}

func (p *Parser) notMissingArticle(block *html.Node) bool {
	// Add conditions here
	return true
	}

func (p *Parser) clearText(text string) string {
    return regexp.MustCompile(`(\r\n|\n)`).ReplaceAllString(text, " ")
}

func (p *Parser) timer(seconds int) {
    time.Sleep(time.Duration(seconds) * time.Second)
}

func (p *Parser) dbExistArticle(articleID string) bool {
	// Checks if there is an article in the database
	return false
}

func (p *Parser) dbLog(logText string) {
	// Logg
	fmt.Println(logText)
}

func (p *Parser) isBlocks() bool {
	// Checks if the object is a valid block
	return false
}

func (p *Parser) ignoreArticle(articleTree *goquery.Document) bool {
	// Ignore article
	return false
}

func (p *Parser) getBlocks() []*goquery.Selection {
	// Returns an array of goquery blocks
	panic("not implemented")
}

func (p *Parser) getArticleLink(blockTree *goquery.Selection) string {
	// Returns a link to the article
	panic("not implemented")
}


func (p *Parser) getArticleTitle(articleTree *goquery.Document) string {
	// Returns the title of the article
	panic("not implemented")
}

func (p *Parser) getFormattedText(articleTree *goquery.Document) string {
	// Returns the text of the article
	panic("not implemented")
}

func (p *Parser) getArticleImages(articleTree *goquery.Document) []map[string]string {
	// Returns an array of found images
	return []map[string]string{{"src":  "img_alt"}}
}

func (p *Parser) getArticleVideos(articleTree *goquery.Document) []string {
	// Returns an array of found video links
	panic("not implemented")
}

func (p *Parser) getArticleDate(articleTree *goquery.Document) time.Time {
	// Returns the date of the article
	return time.Now()
}

