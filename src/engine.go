package src

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"os"
)

var (
	globalTermsDatabase terms
	seenDocs            map[string]struct{}
	databasePath        string = "terms-data.gz"
)

func DoneIndexing() {
	calcIdfScores()
	fmt.Println("serializing index...")
	serializeDatabase()
}

func EngineStart() {
	loadIndex()
}

func loadIndex() {
	fmt.Println("loading index...")
	deserializeDatabase()
	loadSeenDocs()
	fmt.Printf("%d terms in index\n", len(globalTermsDatabase))
	fmt.Printf("%d docs in index\n", corpusCount())
}

func seenDoc(doc string) bool {
	_, seen := seenDocs[doc]
	return seen
}

func seenToken(token string) bool {
	_, seen := globalTermsDatabase[token]
	return seen
}

func corpusCount() int {
	return len(seenDocs)
}

func termsCount() int {
	return len(globalTermsDatabase)
}

func calcIdfScores() {
	for _, tData := range globalTermsDatabase {
		tData.Idf = idf(corpusCount(), len(tData.Docs))
	}
}

func addToIndex(docName string, tkns *tokenizedDoc) {
	seenDocs[docName] = struct{}{}
	var currentTermIndex = len(globalTermsDatabase)
	for token, tf := range tkns.tokens {
		if !seenToken(token) {
			globalTermsDatabase[token] = &termEntry{
				currentTermIndex,
				0,
				map[string]float64{docName: 0},
			}
			currentTermIndex++
		}
		tkn := globalTermsDatabase[token]
		tkn.Docs[docName] = float64(tf) / float64(tkns.docLen)
	}
}

func loadSeenDocs() {
	seenDocs = map[string]struct{}{}
	for _, tData := range globalTermsDatabase {
		for doc := range tData.Docs {
			if !seenDoc(doc) {
				seenDocs[doc] = struct{}{}
			}
		}
	}
}

func deserializeDatabase() {
	if _, err := os.Stat(databasePath); err != nil {
		globalTermsDatabase = make(terms)
		return
	}
	file, err := os.Open(databasePath)
	checkErr(err)
	zr, err := gzip.NewReader(file)
	checkErr(err)
	globalTermsDatabase = make(terms)
	gd := gob.NewDecoder(zr)
	globalTermsDatabase = make(terms)
	gd.Decode(&globalTermsDatabase)
}

func serializeDatabase() {
	file, err := os.Create(databasePath)
	checkErr(err)
	zw := gzip.NewWriter(file)
	ge := gob.NewEncoder(zw)
	err = ge.Encode(globalTermsDatabase)
	checkErr(err)
	zw.Close()
	file.Close()
}

func idf(corpusLen, containingTermLen int) float64 {
	if containingTermLen == 0 {
		return 0
	}
	return math.Log10(float64(corpusLen) / float64(containingTermLen))
}

type terms map[string]*termEntry

type termEntry struct {
	Index int
	Idf   float64
	Docs  map[string]float64
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
