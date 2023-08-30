package src

import (
	"strings"
	"golang.org/x/net/html"
)

const BASE_URL = "https://en.wikipedia.org"

type wikipediaParser struct {} 

func (w wikipediaParser) parse(unparsedData string, thisUrl string) *parsedHtml {
	tkn := html.NewTokenizer(strings.NewReader(unparsedData))
	var ph parsedHtml
	var inP bool = false
	for  {
		tt := tkn.Next()
		switch tt {
		case html.ErrorToken:
			ph.urls = filterUrls(ph.urls, thisUrl)
			return &ph

		case html.StartTagToken:
			t := tkn.Token()
			if t.Data == "p" {
				inP = true
			} else if t.Data == "a" {
				for _, attr := range t.Attr {
					if attr.Key == "href" && len(attr.Val) > 0 && attr.Val[0] != '#' {
						ph.urls = append(ph.urls, attr.Val)
					}
				}
			}

		case html.EndTagToken:
			t := tkn.Token()
			if t.Data == "p" {
				inP = false
			}

		case html.TextToken: 

			if inP {
				ph.text += " " + tkn.Token().Data
			}

		}
	}
}

func filterUrls(unfilteredUrls []string, thisUrl string) []string {
	// need to check if url was already indexed
	// in case of relative urls need to add to base of thisUrl
	var result []string 
	for _, u := range unfilteredUrls {
		if seenDoc(thisUrl) {
			continue
		}
		if strings.Contains(u, ":") {
			continue
		}
		if strings.HasPrefix(u, BASE_URL) {
			result = append(result, u)
		} else if strings.HasPrefix(u, "/wiki/") {
			result = append(result, BASE_URL + u)
		}
	}
	return result
}