package src

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"sort"
	"encoding/json"
)

func StartHandlingQueries() {
	initTfidfvectors()
	mux := http.NewServeMux()
	mux.HandleFunc("/", queryHandle)	
	err := http.ListenAndServe(":1234", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed")
	} else {
		fmt.Println("could not start server")
		os.Exit(1)
	}
}

const RESULTS_AMOUNT = 30

func queryHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	query := r.URL.Query().Get("q")
	if query == "" {
		return
	}
	results := rank(query, RESULTS_AMOUNT)
	jsonData, err := json.Marshal(results)
	checkErr(err)

	fmt.Fprint(w, string(jsonData))
}

var (
	tfidfVectors map[string]sparseVector    // document name: sparse vector of its terms tf-idf scores
)

type sparseVector map[int] float64  // index: val - actual length is len(termsDatabase)

func initTfidfvectors() {
	// needs to create tf-idf vectors
	tfidfVectors = make(map[string]sparseVector)
	for _, termData := range globalTermsDatabase {

		var termIdf = termData.Idf
		var termIndex = termData.Index

		for docName, tf := range termData.Docs {
			_, encountered := tfidfVectors[docName]
			if !encountered {
				tfidfVectors[docName] = make(sparseVector)
			}
			docVec := tfidfVectors[docName]
			docVec[termIndex] = tf * termIdf
		}
	}
}

func rank(query string, n int) []string {
	// rank only upto n results
	qv := vectorizeQuery(query)
	type rankScore struct {
		name  string
		score float64
	}

	var ranks = make([]rankScore, n)

	for doc, vec := range tfidfVectors {
		ranks = append(ranks, rankScore{doc, cosineSimilarity(qv, vec)})
	}
	sort.Slice(ranks, func(i, j int) bool {
		return ranks[i].score > ranks[j].score
	})	
	var topDocNames  = []string {}
	for _, d := range ranks {
		if d.score == float64(0) {
			break
		}
		topDocNames = append(topDocNames, d.name)
	}
	return topDocNames
}

func vectorizeQuery(query string) sparseVector {
	// not calcing the query's tf now
	var qv = make(sparseVector)
	// var queryTokens = tokenize(query).tokens
	var queryTkns = tokenize(query)
	
	for token, tf := range queryTkns.tokens {
		_, encounteredTerm := globalTermsDatabase[token]
		if encounteredTerm {
			var token = globalTermsDatabase[token]
			qv[token.Index] = token.Idf * (float64(tf) / float64(queryTkns.docLen))
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