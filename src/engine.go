package src

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

var (
	globalTermsDatabase terms
	corpusLen           int
	indexedDocs         map[string]struct {}
	databasePath        string = "terms-data.gz"
)

// func OnStart() {
// 	// need to load in index
// 	loadIndex()
// 	// need to deserialize already pinged urls
// 	indexedDocs = map[string]struct {}{}
// }

func Cleanup() {
	serializeDatabase()
}

func EngineStart() {
	loadIndex()
}

func loadIndex() {
	err := deserializeDatabase()
	corpusLen = len(globalTermsDatabase)
	checkErr(err)
	fmt.Printf("%d terms in database\n", len(globalTermsDatabase))
}

func seenDoc(doc string) bool {
	_, seen := indexedDocs[doc]
	return seen
}

func seenToken(token string) bool {
	_, seen := indexedDocs[token]
	return seen
}

func addToIndex(docName string, tkns *tokenizedDoc) {
	indexedDocs[docName] = struct {}{}
	corpusLen++
	var currentTermIndex = 0
	for token, tf := range tkns.tokens {
		if !seenToken(token) {
			globalTermsDatabase[token] = &termData{
				currentTermIndex,
				0,
				map[string]float64 {docName: 0}, 
			}
			currentTermIndex++
		}
		globalTermsDatabase[token].Docs[docName] = float64(tf) / float64(tkns.docLen)
	}

	// update idf scores now that we know each term's docs len
	for _, tData := range globalTermsDatabase {
		tData.Idf = idf(corpusLen, len(tData.Docs))
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

// func QueryAndDisplay(query string) {
// 	globalTermsdatabase = deserializeData("temp.gz")
// 	fmt.Printf("%d terms in database\n", len(globalTermsdatabase))
// 	start := time.Now()
// 	vectors := createTfIdfVectors()
// 	fmt.Printf("%s: %s\n", query, top(vectors, query))
// 	fmt.Println("-------------------------------------------------------")
// 	fmt.Printf("query took %f seconds\n", time.Since(start).Seconds())
// }

// func top(tfidfVectors map[string]sparseVector, query string) string {
// 	qv := vectorizeQuery(query)
// 	var minDistance = float64(0)
// 	var minDistanceName = "no match found"
// 	for doc, vec := range tfidfVectors {
// 		cosineAngle := cosineSimilarity(qv, vec)
// 		fmt.Printf("%s: %f\n", doc, cosineAngle)
// 		if cosineAngle > minDistance {
// 			minDistance = cosineAngle
// 			minDistanceName = doc
// 		}
// 	}
// 	fmt.Println("-------------------------------------------------------")
// 	return minDistanceName
// }

// func vectorizeQuery(query string) sparseVector {
// 	// not calcing the query's tf now
// 	var qv = make(sparseVector)
// 	queryTokens, _ := tknzDoc(query)
// 	for token, tf := range queryTokens {
// 		_, encounteredTerm := globalTermsdatabase[token]
// 		if encounteredTerm {
// 			var token = globalTermsdatabase[token]
// 			qv[token.Index] = token.Idf * float64(tf)
// 		}
// 	}
// 	return qv
// }

// func cosineSimilarity(a, b sparseVector) float64 {
// 	// cosine simitlarity: (A dot B) / (||A|| * ||B||)
// 	aDotb := dotProduct(a, b)
// 	aMag := vectorMagnitude(a)
// 	bMag := vectorMagnitude(b)
// 	if aMag == 0 || bMag == 0 {
// 		return 0
// 	}
// 	return aDotb / (aMag * bMag)
// }

// func dotProduct(a, b sparseVector) float64 {
// 	// iterate over A or B. does not matter which one
// 	var product float64
// 	for index := range a {
// 		product += a[index] * b[index]
// 	}
// 	return product
// }

// func vectorMagnitude(vec sparseVector) float64 {
// 	var entriesSquared float64
// 	for _, val := range vec {
// 		entriesSquared += val * val
// 	}
// 	return math.Sqrt(entriesSquared)
// }

// type sparseVector map[int] float64  // index: val - actual length is len(termsDatabase)

type terms map[string]*termData

type termData struct {
	Index int
	Idf   float64
	Docs  map[string]float64
}

// func initGlobalTermsDatabase(docs map[string]string) terms {
// 	termsDatabase := make(terms)
// 	var currentTermIndex = 0
// 	for docName, docData := range docs {
// 		docTokens, docLen := tknzDoc(docData)
// 		for token, tf := range docTokens {
// 			_, encounteredterm := termsDatabase[token]
// 			if !encounteredterm {
// 				termsDatabase[token] = &termData{
// 					currentTermIndex,
// 					0,        // placeholder
// 					map[string]float64 {docName: 0},
// 				}
// 				currentTermIndex++
// 			}
// 			termsDatabase[token].Docs[docName] = float64(tf) / float64(docLen)
// 		}
// 	}
// 	// update idf scores now that we know each term's docs len
// 	for _, tData := range termsDatabase {
// 		tData.Idf = idf(corpusLen, len(tData.Docs))
// 	}
// 	return termsDatabase
// }

// func deserializeData(fileName string) terms {
// 	file, err := os.Open(fileName)
// 	checkErr(err)
// 	zr, err := gzip.NewReader(file)
// 	checkErr(err)
// 	var buff bytes.Buffer
// 	_, err = io.Copy(&buff, zr)
// 	checkErr(err)
// 	var termsDatabase = make(terms)
// 	json.Unmarshal(buff.Bytes(), &termsDatabase)
// 	return termsDatabase
// }

// func createTfIdfVectors() map[string]sparseVector {
// 	var tfidfVectors = make(map[string]sparseVector)
// 	for _, termData := range globalTermsdatabase {
// 		var termIdf = termData.Idf
// 		var termIndex = termData.Index
// 		for docSlug, tf := range termData.Docs {
// 			_, encountered := tfidfVectors[docSlug]
// 			if !encountered {
// 				tfidfVectors[docSlug] = make(sparseVector)
// 			}
// 			docVec := tfidfVectors[docSlug]
// 			docVec[termIndex] = float64(tf) * termIdf
// 		}
// 	}
// 	return tfidfVectors
// }

func idf(corpusLen, containingTermLen int) float64 {
	if containingTermLen == 0 {
		return 0
	}
	return math.Log10(float64(corpusLen) / float64(containingTermLen))
}

// func printfCompleteSparseVector(v sparseVector) {
// 	for i := 0; i < len(globalTermsdatabase); i++ {
// 		fmt.Printf("\t%f\n", v[i])
// 	}
// }

// func printSparseVectors(vectors map[string]sparseVector) {
// 	for doc, vec := range vectors {
// 		fmt.Println(doc)
// 		printfCompleteSparseVector(vec)
// 		fmt.Println()
// 	}
// }

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
