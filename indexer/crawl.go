package indexer

import (
	"io"
	"log"
	"net/http"
	"strings"
	"golang.org/x/net/html"
	"os"
)

func LoadHtml() []string {
	text := getHtml("https://en.wikipedia.org/wiki/Laptop")
	// text := readHtmlFromFile("rob-pike.html")
	return parse(text)
}

func writeToDisk(data []string, fileName string) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for _, line := range data {
		_, err := file.WriteString(line)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func readHtmlFromFile(filepath string) string {
	rawData, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	return string(rawData)
}

func getHtml(link string) string {
	response, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(response.Header.Get("Content-type"))

	return string(bytes)
}

func parse(text string) []string {
	tkn := html.NewTokenizer(strings.NewReader(text))
	var vals []string
	var containsValubleText bool
	for {
		tt := tkn.Next()
		switch {
		case tt == html.ErrorToken:
			return vals

		case tt == html.StartTagToken || tt == html.EndTagToken:
			t := tkn.Token()
			containsValubleText = !(t.Data == "script" || t.Data == "style")

		case tt == html.TextToken:
			if containsValubleText {
				t := tkn.Token()
				vals = append(vals, t.Data)
			}
		}
	}
}
