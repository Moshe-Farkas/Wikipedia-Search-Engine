package src

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
)

func stopCrawling() {
	for i := 0; i < consumerCount+producerCount; i++ {
		quitChn <- true
	}
}

func startCrawlingAndIndexing(wg *sync.WaitGroup, initialLink string) {
	quitChn = make(chan bool)
	urlsChn := make(chan string, 5)
	responseChn := make(chan dataToParse, 100)
	urlsChn <- initialLink
	internalWg := sync.WaitGroup{}
	s := &session{
		&internalWg,
		urlsChn,
		responseChn,
		wikipediaParser{},
	}
	for i := 0; i < producerCount; i++ {
		internalWg.Add(1)
		go produceData(s)
	}
	for i := 0; i < consumerCount; i++ {
		internalWg.Add(1)
		go consumeData(s)
	}

	internalWg.Wait()
	close(responseChn)
	close(urlsChn)
	fmt.Println("finshed by crawling")
	wg.Done()
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

const MAX_URL_BUFFER_LENGTH = 3000
const consumerCount = 1
const producerCount = 4

var quitChn chan bool

func produceData(s *session) {
	for {
		select {
		case <- quitChn:
			goto Finish

		case url := <-s.urlsChn:
			if seenDoc(url) {
				continue
			}
			response, err := http.Get(url)
			if err != nil {
				continue
			}
			if response.StatusCode == 429 {
				panic(errors.New("too fast. got a 429. deal with it"))
			}
			body, err := io.ReadAll(response.Body)
			if err != nil {
				continue
			}
			response.Body.Close()
			s.responseChn <- dataToParse{
				url,
				body,
			}
		}
	}
Finish:
	s.wg.Done()
}

func consumeData(s *session) {
	var bufferedUrls []string
	for {
		select {
		case <-quitChn:
			goto Finish
		case response := <-s.responseChn:
			ph := s.htmlParser.parse(string(response.rawData), response.url)
			if len(bufferedUrls) < MAX_URL_BUFFER_LENGTH {
				bufferedUrls = append(bufferedUrls, ph.urls...) // parsing strategy will filtered out unwanted urls
			}
			addToIndex(response.url, tokenize(ph.text))
			// fmt.Printf("%s. doc count: %d words: %d\n", response.url, dbConn.corpusLength(), dbConn.termsCount())

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
Finish:
	s.wg.Done()
}

type parsedHtml struct {
	text string
	urls []string
}
