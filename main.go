package main

import (
	"flag"
	"os"
	"os/signal"
	"project/src"
	"strings"
)

const (
	QUERY_MODE   = "query"
	INDEX_MODE   = "index"
)

var (
	mode         *string
)

func handleQuit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<- c
	if *mode == INDEX_MODE {
		src.StopIndexingSession()
	} else if *mode == QUERY_MODE {
		src.CloseDB()
		os.Exit(0)
	}
}

func main() {
	go handleQuit()
	mode = flag.String("mode", QUERY_MODE, "idk")
	flag.Parse()
	src.StartDB()
	switch *mode {
	case QUERY_MODE:
		src.StartHandlingQueries()
	case INDEX_MODE:
		initialLink := os.Args[len(os.Args)-1]
		initialLink = strings.ReplaceAll(initialLink, `'`, "")
		src.StartIndexingSession(initialLink)
		src.CloseDB()
	}
}
