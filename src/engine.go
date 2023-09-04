package src

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"math"
)

var (
	globalTermsDatabase terms
	seenDocs         map[string]struct {}
	databasePath        string = "terms-data.gz"
)

func Cleanup() {
	calcIdfScores()
	serializeDatabase()
}

func EngineStart() {
	loadIndex()
}

func loadIndex() {
	err := deserializeDatabase()
	loadSeenDocs()
	checkErr(err)
	fmt.Printf("%d terms in database\n", len(globalTermsDatabase))
}

func seenDoc(doc string) bool {
	_, seen := seenDocs[doc]
	return seen
}

func seenToken(token string) bool {
	_, seen := globalTermsDatabase[token]
	return seen
}

func corpusLen() int {
	return len(seenDocs)
}

func calcIdfScores() {
	for _, tData := range globalTermsDatabase {
		tData.Idf = idf(corpusLen(), len(tData.Docs))
	}
}

func addToIndex(docName string, tkns *tokenizedDoc) {
	seenDocs[docName] = struct {}{}
	var currentTermIndex = len(globalTermsDatabase) 
	for token, tf := range tkns.tokens {
		if !seenToken(token) {
			globalTermsDatabase[token] = &termData {
				currentTermIndex,
				0,
				map[string]float64 {docName: 0}, 
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

func deserializeDatabase() error {
	if _, err := os.Stat(databasePath); err != nil {
		globalTermsDatabase = make(terms)
		// this means the database of terms is empty
		return nil
	}
	file, err := os.Open(databasePath)	
	if err != nil {
		return err
	}
	zr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	var buff bytes.Buffer
	_, err = io.Copy(&buff, zr)
	if err != nil {
		return err
	}
	globalTermsDatabase = make(terms)
	err = json.Unmarshal(buff.Bytes(), &globalTermsDatabase)
	if err != nil {
		return err
	}
	return nil
}

func serializeDatabase() {
	jsonData, err := json.Marshal(globalTermsDatabase)
	checkErr(err)
	var buff bytes.Buffer
	zw := gzip.NewWriter(&buff)
	_, err = zw.Write(jsonData)
	checkErr(err)
	err = zw.Close()
	checkErr(err)
	file, err := os.Create(databasePath)
	checkErr(err)
	_, err = file.Write(buff.Bytes())
	checkErr(err)
}

func idf(corpusLen, containingTermLen int) float64 {
	if containingTermLen == 0 {
		return 0
	}
	return math.Log10(float64(corpusLen) / float64(containingTermLen))
}

type terms map[string]*termData

type termData struct {
	Index int
	Idf   float64
	Docs  map[string]float64
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
