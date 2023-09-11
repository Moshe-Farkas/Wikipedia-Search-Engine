package src

import (
	"log"
)

var (
	dbConn *databaseConn
)

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

func addToIndex(docName string, tkns *tokenizedDoc) {
	dbConn.bufferNewDoc(docName, tkns)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// func populateFromStrings(docs map[string]string) {
// 	termsInsertStr :=
// 	`
// 	insert into terms(termname)
// 	values ($1)
// 	`
// 	insertTerm, err := db.Prepare(termsInsertStr)
// 	checkErr(err)
// 	defer insertTerm.Close()

// 	docInsertStr :=
// 	`
// 	insert into docs(docname)
// 	values ($1)
// 	`
// 	insertDoc, err := db.Prepare(docInsertStr)
// 	checkErr(err)
// 	defer insertDoc.Close()

// 	termEntryStr :=
// 	`
// 	insert into termentry(termname, docname, tfscore)
// 	values ($1, $2, $3)
// 	`
// 	insertTermEntry, err := db.Prepare(termEntryStr)
// 	checkErr(err)
// 	defer insertTermEntry.Close()

// 	for docName, docData := range docs {
// 		tokens := tempTokenize(docData)
// 		addToIndex(docName, tokens)
// 	}
// }

// func tempTokenize(doc string) *tokenizedDoc {
// 	result := newTokenizedDoc()
// 	for _, term := range strings.Split(doc, " ") {
// 		result.addToken(term)
// 	}
// 	return result
// }