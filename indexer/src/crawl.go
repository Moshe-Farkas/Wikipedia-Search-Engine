package src

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"golang.org/x/net/html"
)

func Idk() {
    text := loadHtml("https://en.wikipedia.org/wiki/Rob_Pike")
	data := parse(text)
	writeToDisk(data, "temp.txt")
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

func loadHtml(link string) string {
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

func parse(text string) (data []string) {

    tkn := html.NewTokenizer(strings.NewReader(text))

    var vals []string
    var containsValubleText bool
    seperation := "\n********************************************************************\n"
    var within string 
    for {
        tt := tkn.Next()
        
        switch {

        case tt == html.ErrorToken:
            return vals

        case tt == html.StartTagToken || tt == html.EndTagToken:
            t := tkn.Token()
            var validTags = []string{
                "p",
                "b",
                "a",
                "i",
                "td",
                "tr",
                "h1",
                "h2",
                "h3",
                "h4",
                "h5",
                "ol",
                "li",
                "ul",
                "s",
                "blockquote",
                "pre",
                "code",
                "bold",
                "strong",
                "u",
            } 
            for _, tag := range validTags {
                if t.Data == tag {
                    containsValubleText = true
                    within = t.String()
                    break
                }
            }            

        case tt == html.TextToken:

            t := tkn.Token()

            if containsValubleText {
                vals = append(vals, within + "  --  " + t.Data +"\n" + seperation + "\n")
            }
            
            containsValubleText = false
        }
    }
}
