package src

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	_ "github.com/lib/pq"
)

var (
	db *sql.DB
	seenTerms map[string]struct {}
	seenDocs map[string]struct {}
)

func CloseDB() {
	db.Close()
}

func StartDB() {
	// loadIndex()

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	connStr := fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", username, password, dbName)
	var err error
	db, err = sql.Open("postgres", connStr)

	checkErr(err)

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("The database is connected")
	fmt.Println("num of docs:", corpusCount())
}

func populateFromStrings(docs map[string]string) {
	termsInsertStr :=
	`
	insert into terms(termname)
	values ($1)
	`
	insertTerm, err := db.Prepare(termsInsertStr)
	checkErr(err)
	defer insertTerm.Close()

	docInsertStr :=
	`
	insert into docs(docname)
	values ($1)
	`
	insertDoc, err := db.Prepare(docInsertStr)
	checkErr(err)
	defer insertDoc.Close()

	termEntryStr :=
	`
	insert into termentry(termname, docname, tfscore)
	values ($1, $2, $3)
	`
	insertTermEntry, err := db.Prepare(termEntryStr)
	checkErr(err)
	defer insertTermEntry.Close()

	for docName, docData := range docs {
		tokens := tempTokenize(docData)
		addToIndex(docName, tokens)
	}
}

func tempTokenize(doc string) *tokenizedDoc {
	result := newTokenizedDoc()
	for _, term := range strings.Split(doc, " ") {
		result.addToken(term)
	}
	return result
}

func resetDB() {
	query := 
	`
	delete from termentry;
	delete from terms;
	alter sequence terms_termindex_seq restart with 1;
	delete from docs;
	`
	_, err := db.Exec(query)
	checkErr(err)
}

func seenDoc(doc string) bool {
	query :=
		`
	select * from docs
	where docname=$1
	`
	row := db.QueryRow(query, doc)
	return row.Scan() != sql.ErrNoRows
}

func seenTerm(term string) bool {
	queryStmnt :=
		`
	select termname from terms
	where termname = $1	
	`
	row := db.QueryRow(queryStmnt, term)
	return row.Scan() != sql.ErrNoRows
}

func corpusCount() int {
	queryStmnt :=
		`
	select count(*) from docs	
	`
	row := db.QueryRow(queryStmnt)
	var len int
	row.Scan(&len)
	return len
}

func termsCount() int {
	queryStmnt :=
		`
	select count(*) from terms
	`
	row := db.QueryRow(queryStmnt)
	var count int
	row.Scan(&count)
	return count
}

func termDocsCount(term string) int {
	// returns the number of docs containing term
	query :=
		`
	select count(*) from termentry
	where termname=$1
	`
	row := db.QueryRow(query, term)
	var count int
	row.Scan(&count)
	return count
}

func addDoc(doc string) {
	insertQuery :=
	`
	insert into docs(docname)
	values ($1)
	`
	_, err := db.Exec(insertQuery, doc)
	checkErr(err)
}

func addTerm(term string) {
	insertQuery :=
	`
	insert into terms(termname, containingcount)
	values ($1, 1)
	`
	_, err := db.Exec(insertQuery, term)
	checkErr(err)
}

func addTermEntry(term string, doc string, tfScore float64) {
	insertQuery :=
	`
	insert into termentry(termname, docname, tfscore)
	values ($1, $2, $3)
	`
	_, err := db.Exec(insertQuery, term, doc, tfScore)
	checkErr(err)
}

func updateTermContainingCount(term string) {
	query :=
	`
	update terms
	set containingcount = containingcount + 1
	where termname=$1
	`
	_, err := db.Exec(query, term)
	checkErr(err)
}

func addToIndex(docName string, tkns *tokenizedDoc) {
	if seenDoc(docName) {
		return 
	}

	if len(docName) >= 100 {
		return 
	}
	addDoc(docName)
	for token, tf := range tkns.tokens {
		if !seenTerm(token) {
			addTerm(token)
		} else {
			updateTermContainingCount(token)
		}
		tfScore := float64(tf) / float64(tkns.docLen)
		addTermEntry(token, docName, tfScore)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
