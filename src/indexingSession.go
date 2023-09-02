package src

import (
	"errors"
	"io"
	"net/http"
	"sync"
	"fmt"
)

func StartCrawlingAndIndexing() {
	// indexFromString()
	urlsChn := make(chan string, 5)
	responseChn := make(chan dataToParse, 100)
	initialLink := "https://en.wikipedia.org/wiki/Tf%E2%80%93idf"

	wg := sync.WaitGroup{}
	urlsChn <- initialLink
	s := &session{
		&wg,
		urlsChn,
		responseChn,
		wikipediaParser{},
	}
	go produceData(s)
	wg.Add(1)
	go consumeData(s)

	wg.Wait()
	close(urlsChn)
	close(responseChn)
}

func printPage() {
	// should remove 
	re, err := http.Get("https://en.wikipedia.org/wiki/PageRank")
	checkErr(err)
	body, err := io.ReadAll(re.Body)
	checkErr(err)
	re.Body.Close()
	
	wp := wikipediaParser {}
	temp := wp.parse(string(body), "https://en.wikipedia.org/wiki/PageRank")
	for _, url := range temp.urls {
		fmt.Println(url)
	}
	fmt.Println("------------------------------------------------")
	// for _, t := range strings.Split(temp.text, " ") {
	// 	fmt.Println(t)
	// } 	
	// os.Exit(1)
}

func indexFromString() {
	docs := map[string]string {
		"doc1": "the brown cow",
		"doc2": "so the brown bag",
		"doc3": "and and and for",
		"doc4": "and and better",
	}
	for doc, docdata := range docs {
		addToIndex(doc, tokenize(docdata))	
	}
}

type session struct {
	wg          *sync.WaitGroup
	urlsChn     chan string
	responseChn chan dataToParse
	htmlParser  parsingStrategy
}

type dataToParse struct {
	url     string
	rawData []byte
}

type parsingStrategy interface {
	parse(string, string) *parsedHtml
}

var tempSeenDocs = map[string] struct {}{}

func produceData(s *session) {
	for url := range s.urlsChn {

		if _, seen := tempSeenDocs[url]; seen {
			continue
		}

		response, err := http.Get(url)
		tempSeenDocs[url] = struct{}{}

		checkErr(err)
		if response.StatusCode == 429 {
			panic(errors.New("too fast. got a 429. deal with it"))
		}
		body, err := io.ReadAll(response.Body)
		response.Body.Close()
		checkErr(err)
		s.responseChn <- dataToParse{
			url,
			body,
		}
	}
}

func consumeData(s *session) {
	var bufferedUrls []string
	var corpusL = 0
	for {
		select {
		case response := <-s.responseChn:
			ph := s.htmlParser.parse(string(response.rawData), response.url)
			bufferedUrls = append(bufferedUrls, ph.urls...) 	// parsing strategy will filtered out unwanted urls

			addToIndex(response.url, tokenize(ph.text))
			fmt.Printf("new index: %s. doc count: %d words: %d\n", response.url, corpusL, len(globalTermsDatabase))
			corpusL++

		default:
			if len(bufferedUrls) > 0 {
				select {
				case s.urlsChn <- bufferedUrls[0]:
					bufferedUrls = bufferedUrls[1:]
				default:
					break
				}
			}
		}
	}
	s.wg.Done()
}

type parsedHtml struct {
	text string
	urls []string
}
