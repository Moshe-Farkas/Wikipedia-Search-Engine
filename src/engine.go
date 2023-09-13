package src

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
		panic(err)
	}
}