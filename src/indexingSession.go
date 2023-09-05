package src

import (
	"errors"
	"io"
	"net/http"
	"sync"
	"fmt"
)

func StartCrawlingAndIndexing(initialLink string) {
	// indexFromString()
	urlsChn := make(chan string, 5)
	responseChn := make(chan dataToParse, 100)
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
		"doc1": "vain brown cow",
		"doc2": "jump vain brown bag",
		"doc3": "placeholder same way",
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


func produceData(s *session) {
	for url := range s.urlsChn {

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
	for {
		select {
		case response := <- s.responseChn:
			ph := s.htmlParser.parse(string(response.rawData), response.url)
			bufferedUrls = append(bufferedUrls, ph.urls...) 	// parsing strategy will filtered out unwanted urls

			addToIndex(response.url, tokenize(ph.text))
			fmt.Printf("%s. doc count: %d words: %d\n", response.url, corpusLen(), len(globalTermsDatabase))

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
}

type parsedHtml struct {
	text string
	urls []string
}
