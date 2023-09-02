package src

import (
	"fmt"
	"math"
	"sort"
	"time"
)

func StartHandlingQueries() {
	// need to create tf-idf vectors
	initialize()
	// printSparseVectors(tfidfVectors)
	var q string = "inverted matrix index"
	// fmt.Println("query")
	// printfCompleteSparseVector(vectorizeQuery(q))
	// printGlobalTerms()
	fmt.Println("-------------------------------------------------------")
	fmt.Println(q)
	startTime := time.Now()

	result := rank(q)
	i := 0
	for _, item := range result {
		if i > 10 {
			break
		}
		fmt.Println("\t", item.name, item.val)
		i++
	}
	fmt.Printf("took %f seconds to query\n", time.Since(startTime).Seconds())
	fmt.Println("-------------------------------------------------------")
}

var (
	tfidfVectors map[string]sparseVector    // document name: sparse vector of its terms tf-idf scores
)

type sparseVector map[int] float64  // index: val - actual length is len(termsDatabase)

func initialize() {
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

func rank(query string) []struct{name string; val float64} {
	qv := vectorizeQuery(query)
	// var minDistance = float64(0)
	// var minDistanceName = "no match found"

	var data []struct {
		name string
		val float64
	}

	for doc, vec := range tfidfVectors {
		data = append(data, struct{name string; val float64}{doc, cosineSimilarity(qv, vec)})


		// cosineAngle := cosineSimilarity(qv, vec)
		// // fmt.Printf("%s: %f\n", doc, cosineAngle)
		// if cosineAngle > minDistance {
		// 	minDistance = cosineAngle
		// 	minDistanceName = doc
		// }
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].val > data[j].val
	})	
	return data
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


