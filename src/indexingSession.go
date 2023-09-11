package src

import (
	"errors"
	"io"
	"net/http"
	"fmt"
	"sync"
)

func StopIndexing() {
	for i := 0; i < consumerCount + producerCount; i++ {
		shouldStop <- true
	}
}

func StartCrawlingAndIndexing(initialLink string) {
	loadTermsAndDocs()
	shouldStop = make(chan bool)
	urlsChn := make(chan string, 5)
	responseChn := make(chan dataToParse, 100)
	if !validInitialLink(initialLink) {
		fmt.Println("seen first link already")
		return
	}
	urlsChn <- initialLink
	wg := sync.WaitGroup {}
	s := &session{
		&wg,
		urlsChn,
		responseChn,
		wikipediaParser{},
	}
	for i := 0; i < producerCount; i++ {
		wg.Add(1)
		go produceData(s)
	}
	for i := 0; i < consumerCount; i++ {
		wg.Add(1)
		go consumeData(s)
	}

	wg.Wait()
	close(responseChn)
	close(urlsChn)
}

func validInitialLink(initialLink string) bool {
	return !seenDoc(initialLink)
}

type session struct {
	wg 			*sync.WaitGroup
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
var shouldStop chan bool

func produceData(s *session) {
	for {
		select {
		case <- shouldStop:
			goto Finish
		
		case url := <- s.urlsChn: 
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
			s.responseChn <- dataToParse {
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
		case <- shouldStop:
			goto Finish
		case response := <- s.responseChn:
			ph := s.htmlParser.parse(string(response.rawData), response.url)
			if len(bufferedUrls) < MAX_URL_BUFFER_LENGTH {
				bufferedUrls = append(bufferedUrls, ph.urls...) 	// parsing strategy will filtered out unwanted urls
			}
			addToIndex(response.url, tokenize(ph.text))
			fmt.Printf("%s. doc count: %d words: %d\n", response.url, corpusCount(), termsCount())

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
