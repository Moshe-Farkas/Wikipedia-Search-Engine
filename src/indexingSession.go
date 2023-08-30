package src

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
)

func StartCrawlingAndIndexing() {
	urlsChn := make(chan string, 5)
	responseChn := make(chan dataToParse, 100)
	initialLink := "https://en.wikipedia.org/wiki/Laptop"
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

func produceData(s *session) {
	for url := range s.urlsChn {
		response, err := http.Get(url)
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

			addToIndex(response.url, tokenize(ph))
			fmt.Printf("Consumer added %s to index\n", response.url)
			corpusL++
			fmt.Println("\t\t", len(globalTermsDatabase), "   docs: ", corpusL)

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
