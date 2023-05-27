package parser

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/guisecreator/pikabu-parser-go/parser/recipient"
)

type Parser struct {
	BaseURL      	string
	ParsRegularExp  string
	MissTags     	[]string
	PublicTags   	bool
	EntryURL     	string
	HTTPClient   	*http.Client
	EntryTree    	*goquery.Document
}

func NewParser(url string, missTags []string, publicTags bool) *Parser {

	RegularExp := regexp.MustCompile(`((http[s]?)://[^/]+)`)

	BaseURL := RegularExp.FindStringSubmatch(url)[1]

	ParsRegularExp := RegularExp.FindStringSubmatch(url)[2]
	return &Parser{
		BaseURL:    		BaseURL,
		MissTags:   		missTags,
		ParsRegularExp:     ParsRegularExp,
		PublicTags: 		publicTags,
	}
}

func (p *Parser) GetPosts() []map[string]interface{} {
	listArticles := p.GetListArticles()
	p.excludePosts(listArticles)
	for i, j := 0, len(listArticles)-1; i < j; i, j = i+1, j-1 {
		listArticles[i], listArticles[j] = listArticles[j], listArticles[i]
	}
	
	formatedArticles := make([]map[string]interface{}, 0)
	for _, article := range listArticles {
		articleLink := strings.Replace(fmt.Sprintf(
			"%v", article["Link"]), 
			"old_string", "new_string", -1)

		errorText  :=    []string{}
		formatedText :=  []string{}
		articleTags :=   []string{}
		
		public := true

		fmt.Printf("Parsing: %s\n", articleLink)

		articleTree := p.GetArticle(articleLink)
		if articleTree == nil {
			fmt.Printf("Error in ID: %s\n", articleLink)
			continue
		}

		p.waiting_timer(0)

		articleTitle := p.getArticleTitle(articleTree)
		if articleTitle == "" {
			fmt.Printf("Error in title: %s\n", articleLink)
			continue
		}

		articleTags = p.GetArticleTags(articleTree.Selection)
		public = p.MissToTags(articleTags)
		if !public {
			errorText = append(errorText, fmt.Sprintf(
				"There is a tag from the list'%s' ", 
				strings.Join(p.MissTags, ",")))
		}

		if !public {
			errorText = append(errorText, fmt.Sprintf(
				"There is a tag from the list'%s' ", strings.Join(p.MissTags, ",")))
		}

		if p.PublicTags && len(articleTags) > 0 {
			formatedText = append(formatedText, 
				"\r\n\r\n", strings.Join(articleTags, " "))
		}

		//TODO: redo it
		// articleID, err := strconv.Atoi(article["Id"].(string)) 
		// if err != nil{
		// 	fmt.Printf("The text is long, post id: %d", articleID)
		// }

		if len(strings.Join(formatedText, "")) >= 7000 {
			public = false
			fmt.Printf("The text is long, post id: %d\n", article["Id"])

			
			errorText = append(errorText, 
				fmt.Sprintf("The text is long, post id: %d", article["Id"]))
		}

		getArticleDate := p.getArticleDate(articleTree)
		if p.ignoreArticle(articleTree) {
			fmt.Printf("Ignore, post id: %d\n", article["Id"])
				formatedArticles = append(formatedArticles, map[string]interface{}{
			        "Id":        article["Id"],
        			"Link":      articleLink,
        			"Text":      strings.Join(formatedText, ""),
        			"Published": getArticleDate,
        			"Public":    false,
        			"Error":     errorText,
			})
		}
	}

	return formatedArticles
}

func (p *Parser) excludePosts(listArticlesID []map[string]string) {
	toRem := make([]map[string]string, 0, len(listArticlesID))
	count := 0

	for _, articleID := range listArticlesID {
		if recipient.DbExistArticle(articleID["Id"]) {
			toRem[count] = articleID
			count++ 
			
			recipient.DbLog(fmt.Sprintf("Has already: %v", articleID["Id"],))
		}
	}

	if count < len(listArticlesID){
		toRem = toRem[:count]
		sort.Slice(toRem, func(i, j int) bool {
			return true
		})
	} 

}

func  (p *Parser) MissToTags(articleTags []string) bool {
	if len(p.MissTags) > 0 {
		for _, atags := range articleTags {
			for _, mtags := range p.MissTags {
				if strings.Contains(strings.ToLower(atags), strings.ToLower(mtags)) {
					recipient.DbLog(fmt.Sprintf("No publish, there is a tag \"%s\"", mtags))
					return false
				}
			}
		}
	}
	return true
}
 
func (p *Parser) waiting_timer(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}
