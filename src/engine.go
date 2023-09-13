package src

import (
	"fmt"
	"sync"
)

var (
	dbConn                 *databaseConn
	intermediateDocsBuffer bufferedDocs
	bufferedDocsChn        chan bufferedDocs
	shouldStopIndexing     chan bool
)

const maxBufferedDocEntries = 100

type bufferedDocs []bufferedDocEntry

func CloseDB() {
	dbConn.Close()
}

func StartDB() {
	dbConn = initDbConn()
}

func seenDoc(doc string) bool {
	return dbConn.seenDoc(doc)
}

func seenTerm(term string) bool {
	return dbConn.seenTerm(term)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func addToIndex(docName string, tkns *tokenizedDoc) {
	intermediateDocsBuffer = append(intermediateDocsBuffer, bufferedDocEntry{docName, tkns})
	// should write?
	if len(intermediateDocsBuffer) >= maxBufferedDocEntries {
		fmt.Printf("buffered %d parsed docs\n", maxBufferedDocEntries)
		bufferedDocsChn <- intermediateDocsBuffer // will block here if channel is >= 2
		intermediateDocsBuffer = make(bufferedDocs, 0)
	}
}

func StartIndexingSession(initialLink string) {
	dbConn.loadTermsAndDocs()
	if seenDoc(initialLink) {
		fmt.Println("seen first link already")
		return
	}
	wg := &sync.WaitGroup{}
	bufferedDocsChn = make(chan bufferedDocs, 2)
	shouldStopIndexing = make(chan bool)

	wg.Add(1)
	go dbIndexing(wg)
	wg.Add(1)
	go startCrawlingAndIndexing(wg, initialLink)

	wg.Wait()
	close(bufferedDocsChn)
	close(shouldStopIndexing)
}

func StopIndexingSession() {
	shouldStopIndexing <- true
	stopCrawling()
}

func dbIndexing(wg *sync.WaitGroup) {
	for {
		select {
		case <-shouldStopIndexing:
			goto End
		case bd := <-bufferedDocsChn:
			dbConn.writeBufferedDocs(bd)
		}
	}
End:
	fmt.Println("finshed by indexing")
	wg.Done()
}
