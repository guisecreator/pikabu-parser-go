package parser

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"github.com/PuerkitoBio/goquery"
)

func (p *Parser) getTree() bool {
	resp, err := http.Get(p.EntryURL)
	if err != nil {
		log.Printf("Failed to get page: %v", err)
		return false
	}
	defer resp.Body.Close()

	pageBytes, err := io.ReadAll(resp.Body)
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