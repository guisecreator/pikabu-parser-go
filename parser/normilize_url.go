package parser

import "strings"

func (p *Parser) normalizeURL(url string) string {
	psex := p.ParsRegularExp
	if !strings.Contains(url, p.BaseURL) && !strings.Contains(url, "http") {
		if url[:1] == "/" && url[:2] != "//" {
			return p.BaseURL + url
		} else if url[:2] == "//" {
			return psex + url
		} else {
			return p.BaseURL + "/" + url
		}
	}
	return url
}