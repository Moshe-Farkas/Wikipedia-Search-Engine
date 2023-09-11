package src 

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"fmt"
)

type bufferedDocEntry struct {
	docName string
	tokens *tokenizedDoc	
}

const maxBufferedDocEntries = 100

type databaseConn struct {
	db *sql.DB

	seenTermsInDB map[string]struct {}
	seenDocsInDB map[string]struct {}

	newSeenTerms map[string]struct {}	// buffering new seen terms. will be copied into seenTermsInDB
	newSeenDocs map[string]struct {}    // buffering new seen docs. will be copied into seenDocsInDB

	bufferedDocsEntries []bufferedDocEntry

	// behind the scenes will write it's buffered data to db
	// need to load seen terms and seen docs for faster lookup time, only when in indexing mode
	// need to seperate old terms already in db and new ones
}

func initDbConn() *databaseConn {
	dbConn := databaseConn {}
	dbConn.newSeenDocs = map[string]struct{}{}
	dbConn.seenDocsInDB = map[string]struct{}{}
	dbConn.newSeenTerms = map[string]struct{}{}
	dbConn.seenTermsInDB = map[string]struct{}{}

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	connStr := fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", username, password, dbName)
	var err error
	dbConn.db, err = sql.Open("postgres", connStr)

	checkErr(err)

	if err = dbConn.db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("The database is connected")
	return &dbConn
}

func (db *databaseConn) Close() {
	db.Close()	
}

func (db *databaseConn) corpusCount() int {
	return len(db.seenDocsInDB) + len(db.newSeenDocs)
}

func (db *databaseConn) termsCount() int {
	return len(db.newSeenTerms) + len(db.seenTermsInDB)
}

func (db *databaseConn) loadTermsAndDocs() {
	// only needed when indexing 
	rows, err := db.db.Query("select termname from terms")
	checkErr(err)
	db.seenTermsInDB = make(map[string]struct{})
	for rows.Next() {
			var termName string
			rows.Scan(&termName)
			db.seenTermsInDB[termName] = struct{}{}
	}
	rows, err = db.db.Query("select docname from docs")
	db.seenDocsInDB = make(map[string]struct{})
	for rows.Next() {
			var doc string
			rows.Scan(&doc)
			db.seenDocsInDB[doc] = struct{}{}
	}
}

func (db *databaseConn) resetDB() {
	query := 
	`
	delete from termentry;
	delete from terms;
	alter sequence terms_termindex_seq restart with 1;
	delete from docs;
	`
	_, err := db.db.Exec(query)
	checkErr(err)
}

func (db *databaseConn) seenDoc(docName string) bool {
	_, seenInDB := db.seenDocsInDB[docName]
	if seenInDB {
		return true
	}
	_, seenThisSession := db.newSeenDocs[docName]
	return seenThisSession
} 

func (db *databaseConn) seenTerm(term string) bool {
	_, seenInDB := db.seenTermsInDB[term]
	if seenInDB {
		return true
	}
	_, seenThisSession := db.newSeenTerms[term]
	return seenThisSession
}

func (db *databaseConn) bufferNewDoc(docName string, tokens *tokenizedDoc) {
	be := bufferedDocEntry {
		docName,
		tokens,
	}
	db.bufferedDocsEntries = append(db.bufferedDocsEntries, be)

	// should write?
	if len(db.bufferedDocsEntries) >= maxBufferedDocEntries {
		fmt.Println("writing to db")
		db.writeBufferedDocs()
	}
}

func (db *databaseConn) writeBufferedDocs() {
	tx, err := db.db.Begin()
	checkErr(err)

	makePreparedStmnt := func(query string) *sql.Stmt {
		smtn, err := tx.Prepare(query)
		checkErr(err)
		return smtn
	}
	addTerm := makePreparedStmnt(
		`
		insert into terms(termname, containingcount)
		values ($1, 1)
		`,
	)
	addDoc := makePreparedStmnt(
		`
		insert into docs(docname)
		values ($1)
		`,
	)
	updateContainingCount := makePreparedStmnt(
		`
		update terms
		set containingcount = containingcount + 1
		where termname=$1
		`,
	)
	addTermEntry := makePreparedStmnt(
		`
		insert into termentry(termname, docname, tfscore)
		values ($1, $2, $3)
		`,
	)
	passOrFail := func(err error) {
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}
	for _, be := range db.bufferedDocsEntries { 
		if db.seenDoc(be.docName) {
			continue
		}
		if len(be.docName) >= 100 {
			continue
		}
		_, err := addDoc.Exec(be.docName)
		db.newSeenDocs[be.docName] = struct{}{}
		passOrFail(err)
		for token, tf := range be.tokens.tokens {
			if !seenTerm(token) {
				_, err := addTerm.Exec(token)
				db.newSeenTerms[token] = struct{}{}
				passOrFail(err)
			} else {
				_, err := updateContainingCount.Exec(token)
				passOrFail(err)
			}
			tfScore := float64(tf) / float64(be.tokens.docLen)
			_, err := addTermEntry.Exec(token, be.docName, tfScore)
			passOrFail(err)
		}
	}
	err = tx.Commit()
	checkErr(err)
	db.bufferedDocsEntries = make([]bufferedDocEntry, 0)
}