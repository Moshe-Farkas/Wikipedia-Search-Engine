package src

import (
	"time"
	"fmt"
)

func StartHandlingQueries() {
	// need to create tf-idf vectors
	for {
		fmt.Println("mock handle query")
		time.Sleep(time.Second)
	}
}



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

