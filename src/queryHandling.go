package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"sort"
)

func printfCompleteSparseVector(v sparseVector) {
	// for i := 0; i < termsCount(); i++ {
	// 	fmt.Printf("\t%f\n", v[i])
	// }
}

func printSparseVectors(vectors map[string]sparseVector) {
	for doc, vec := range vectors {
		fmt.Println(doc)
		printfCompleteSparseVector(vec)
		fmt.Println()
	}
}

func StartHandlingQueries() {
	initTfidfvectors()
	fmt.Println("finished tf vectors...")
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
	tfidfVectors map[string]sparseVector // document name: sparse vector of its terms tf-idf scores
)

type sparseVector map[int]float64 // index: val - actual length is len(termsDatabase)

func initTfidfvectors() {
	cp := dbConn.corpusLength()
	tfidfVectors = make(map[string]sparseVector)
	rows := dbConn.termEntryRows()
	for rows.Next() {
		var termIndex uint32
		var containingcount int
		var termName string
		var docName string
		var tfScore float64
		err := rows.Scan(&termIndex, &containingcount, &termName, &docName, &tfScore)
		checkErr(err)
		_, encountered := tfidfVectors[docName]
		if !encountered {
			tfidfVectors[docName] = make(sparseVector)
		}
		docVec := tfidfVectors[docName]
		docVec[int(termIndex)] = tfScore * idf(cp, containingcount)
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
	var topDocNames = []string{}
	for _, d := range ranks {
		if d.score == float64(0) {
			break
		}
		topDocNames = append(topDocNames, d.name)
	}
	return topDocNames
}

func vectorizeQuery(query string) sparseVector {
	var qv = make(sparseVector)
	var queryTkns = tokenize(query)
	var cp = dbConn.corpusLength()
	for token, tf := range queryTkns.tokens {
		termIndex, containingCount, err := dbConn.termInfo(token)
		if err == nil {
			// aka seen term
			idfScore := idf(cp, containingCount)
			qv[int(termIndex)] = idfScore * float64(tf) / float64(queryTkns.docLen)
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

func idf(corpusLen, containingTermLen int) float64 {
	if containingTermLen == 0 {
		return 0
	}
	return math.Log10(float64(corpusLen) / float64(containingTermLen)) 
}
