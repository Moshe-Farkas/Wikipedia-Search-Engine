package main

import (
	"flag"
	"fmt"
	"os"
	"project/src"
	"strings"
)

const (
	QUIT_MESSAGE = "q"
	QUERY_MODE   = "query"
	INDEX_MODE   = "index"
)

func handleQuit() {
	var temp string = ""
	for temp != QUIT_MESSAGE {
		fmt.Scanf("%s", &temp)
	}
	if *mode == INDEX_MODE {
		src.FinishIndexing()
	}
	os.Exit(0)
}

var mode *string

func main() {
	go handleQuit()
	mode = flag.String("mode", QUERY_MODE, "idk")
	flag.Parse()
	src.EngineStart()
	switch *mode {
	case QUERY_MODE:
		// do query stuff
		src.StartHandlingQueries()
	case INDEX_MODE:
		// do index stuff
		initialLink := os.Args[len(os.Args)-1]
		initialLink = strings.ReplaceAll(initialLink, `'`, "")
		src.StartCrawlingAndIndexing(initialLink)
		src.FinishIndexing()
	}
}
