package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"time"
)

func tempRead(fileName string) string {
	data, err := os.ReadFile(fileName)
	checkErr(err)
	return string(data)	
}

func main() {
	createAndCompileData()
	queryAndDisplay("krautrock")
}

func createAndCompileData() {
	fileName := "temp.gz"
	data := map[string]string {
		"tf-idf": tempRead("tf-idf.txt"),
		"laptop.txt": tempRead("laptop.txt"),
		"ice-cream.txt": tempRead("ice-cream.txt"),
		"addition.txt": tempRead("addition.txt"),
		"xor.txt": tempRead("xor.txt"),
		"kraftwerk.txt": tempRead("kraftwerk.txt"),
	}
	compileIndex(fileName, data)		
}

func queryAndDisplay(query string) {
	globalTermsdatabase = deserializeData("temp.gz")
	fmt.Printf("%d terms in database\n", len(globalTermsdatabase))
	start := time.Now()
	vectors := createTfIdfVectors()
	fmt.Printf("%s: %s\n", query, top(vectors, query))
	fmt.Println("-------------------------------------------------------")
	fmt.Printf("query took %f seconds\n", time.Since(start).Seconds())
}

var globalTermsdatabase terms
var corpusLen int

func top(tfidfVectors map[string]sparseVector, query string) string {
	qv := vectorizeQuery(query)
	var minDistance = float64(0)
	var minDistanceName = "no match found"
	for doc, vec := range tfidfVectors {
		cosineAngle := cosineSimilarity(qv, vec)
		fmt.Printf("%s: %f\n", doc, cosineAngle)
		if cosineAngle > minDistance {
			minDistance = cosineAngle
			minDistanceName = doc
		}	
	}
	fmt.Println("-------------------------------------------------------")
	return minDistanceName
}

func vectorizeQuery(query string) sparseVector {
	// not calcing the query's tf now
	var qv = make(sparseVector)
	for _, term := range strings.Split(query, " ") {
		_, encounteredTerm := globalTermsdatabase[term]
		if encounteredTerm {
			var term = globalTermsdatabase[term]
			qv[term.Index] = term.Idf
		}
	}
	return qv	
}

func cosineSimilarity(a, b sparseVector) float64 {
	// cosine simitlarity: (A dot B) / (||A|| * ||B||)
	aDotb := dotProduct(a, b)	
	aMag := vectorMagnitude(a)
	bMag := vectorMagnitude(b)
	if aMag == 0 || bMag == 0 {
		return 0
	}
	return aDotb / (aMag * bMag)
}

func dotProduct(a, b sparseVector) float64 {
	// iterate over A or B. does not matter which one
	var product float64
	for index := range a {
		product += a[index] * b[index]
	}	
	return product
}

func vectorMagnitude(vec sparseVector) float64 {
	var entriesSquared float64
	for _, val := range vec {
		entriesSquared += val * val
	}
	return math.Sqrt(entriesSquared)
}

type sparseVector map[int] float64  // index: val - actual length is len(termsDatabase)

type terms map[string] *termData

type termData struct {
	Index int
	Idf float64
	Docs map[string] float64
}

func compileIndex(fileName string, docs map[string]string) {
	corpusLen = len(docs)
	termsDatabase := initGlobalTermsDatabase(docs)
	jsonData, err := json.Marshal(termsDatabase)
	checkErr(err)
	var buff bytes.Buffer
	zw := gzip.NewWriter(&buff)
	_, err = zw.Write(jsonData)
	checkErr(err)
	err = zw.Close()
	checkErr(err)
	file, err := os.Create(fileName)
	checkErr(err)
	_, err = file.Write(buff.Bytes())
	checkErr(err)
}

func tokenizeDoc(doc string) (map[string]int, int) {
	tokens := make(map[string]int)
	var count int = 0
	for _, token := range strings.Split(doc, " ") {
		tokens[strings.ToLower(token)]++	
		count++	
	}
	return tokens, count
}

func initGlobalTermsDatabase(docs map[string]string) terms {
	termsDatabase := make(terms)
	var currentTermIndex = 0
	for docName, docData := range docs {
		docTokens, docLen := tokenizeDoc(docData)		
		for token, tf := range docTokens {
			_, encounteredterm := termsDatabase[token]	
			if !encounteredterm {
				termsDatabase[token] = &termData{
					currentTermIndex,
					0,        // placeholder
					map[string]float64 {docName: 0},
				}
				currentTermIndex++
			}
			termsDatabase[token].Docs[docName] = float64(tf) / float64(docLen)
		}
	}
	// update idf scores now that we know each term's docs len
	for _, tData := range termsDatabase {
		tData.Idf = idf(corpusLen, len(tData.Docs))
	}
	return termsDatabase
}

func deserializeData(fileName string) terms {
	file, err := os.Open(fileName)
	checkErr(err)
	zr, err := gzip.NewReader(file)
	checkErr(err)
	var buff bytes.Buffer
	_, err = io.Copy(&buff, zr)
	checkErr(err)
	var termsDatabase = make(terms)
	json.Unmarshal(buff.Bytes(), &termsDatabase)
	return termsDatabase
}

func createTfIdfVectors() map[string]sparseVector {
	var tfidfVectors = make(map[string]sparseVector)
	for _, termData := range globalTermsdatabase {
		var termIdf = termData.Idf
		var termIndex = termData.Index
		for docSlug, tf := range termData.Docs {
			_, encountered := tfidfVectors[docSlug]
			if !encountered {
				tfidfVectors[docSlug] = make(sparseVector)
			}
			docVec := tfidfVectors[docSlug]
			docVec[termIndex] = float64(tf) * termIdf
		}
	}
	return tfidfVectors
}

func idf(corpusLen, containingTermLen int) float64 {
	if containingTermLen == 0 {
		return 0
	}
	return math.Log10(float64(corpusLen) / float64(containingTermLen))
}

func printfCompleteSparseVector(v sparseVector) {
	for i := 0; i < len(globalTermsdatabase); i++ {
		fmt.Printf("\t%f\n", v[i])
	}	
}

func printSparseVectors(vectors map[string]sparseVector) {
	for doc, vec := range vectors {
		fmt.Println(doc)
		printfCompleteSparseVector(vec)
		fmt.Println()
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}