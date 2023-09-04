package src

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"sort"
	// "time"
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
	qv := vectorizeQuery(query)
	var data []struct {
		name string
		val float64
	}
	for doc, vec := range tfidfVectors {
		data = append(data, struct{name string; val float64}{doc, cosineSimilarity(qv, vec)})
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].val > data[j].val
	})	
	var topDocNames  = []string {}
	for _, d := range data {
		if d.val == float64(0) {
			break
		}
		topDocNames = append(topDocNames, d.name)
	}
	return topDocNames
}

func printfCompleteSparseVector(v sparseVector) {
	for i := 0; i < len(globalTermsDatabase); i++ {
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

func printGlobalTerms() {
	for term, tData := range globalTermsDatabase {
		fmt.Printf("%s: index %d  idf %f   docs %v\n", term, tData.Index, tData.Idf, tData.Docs)
	}
}

func vectorizeQuery(query string) sparseVector {
	// not calcing the query's tf now
	var qv = make(sparseVector)
	var queryTokens = tokenize(query).tokens
	for token, tf := range queryTokens {
		_, encounteredTerm := globalTermsDatabase[token]
		if encounteredTerm {
			var token = globalTermsDatabase[token]
			qv[token.Index] = token.Idf * float64(tf)
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



// func main() {
//         mux := http.NewServeMux()
//         mux.HandleFunc("/", Root)
//         mux.HandleFunc("/hello", Hello)
//         mux.HandleFunc("/headers", headers)
//         mux.HandleFunc("/requests", requestCount)
//         mux.HandleFunc("/jsontest", jsonTest)
//         err := http.ListenAndServe(":1234", mux)
//         if errors.Is(err, http.ErrServerClosed) {
//                 fmt.Printf("server closed\n")
//         } else if err != nil {
//                 fmt.Printf("could not start server: %s\n", err)
//                 os.Exit(1)
//         }
// }

// func jsonTest(w http.ResponseWriter, r *http.Request) {
//         w.Header().Set("Content-Type", "application/json")
//         fmt.Fprintln(w, `{
//         "person": {
//                 "name": "John Doe",
//                 "age": 30,
//                 "email": "john@example.com",
//                 "address": {
//                 "street": "123 Main St",
//                 "city": "Anytown",
//                 "state": "CA",
//                 "postal_code": "12345"
//                 },
//                 "hobbies": ["reading", "hiking", "cooking"],
//                 "is_student": false
//                 }
//         } `)
// }

// func requestCount(w http.ResponseWriter, r *http.Request) {
//         globalRequestCount++
//         fmt.Fprintf(w, "you are the %d request\n", globalRequestCount)
// }

// func headers(w http.ResponseWriter, r *http.Request) {
//         for name, headers := range r.Header {
//                 for _, h := range headers {
//                         fmt.Fprintf(w, "%v: %v\n", name, h)
//                 }
//         }
//         globalRequestCount++
// }

// func Root(w http.ResponseWriter, r *http.Request) {
//         io.WriteString(w, "this is the root\n")
//         globalRequestCount++
// }

// func Hello(w http.ResponseWriter, r *http.Request) {
//         // ctx := r.Context()
//         // fmt.Println("function handler")
//         // defer fmt.Println("function handler ended")

// }