package tags

import (
	"strings"
	"github.com/guisecreator/pikabu-parser-go/parser"
)

func missToTags(articleTags []string) bool {
	if len(parser.MissTags) > 0 {
		for _, atags := range articleTags {
			for _, mtags := range parser.MissTags {
				if strings.Contains(strings.ToLower(atags), strings.ToLower(mtags)) {
					// fmt.Sprintf("No publish, there is a tag \"%s\"", mtags)
					return false
				}
			}
		}
	}
	return true
}