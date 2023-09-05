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

var (
	mode         *string
)

func handleQuit() {
	var temp string = ""
	for temp != QUIT_MESSAGE {
		fmt.Scanf("%s", &temp)
	}
	if *mode == INDEX_MODE {
		src.StopIndexing()
	} else if *mode == QUERY_MODE {
		os.Exit(0)
	}
}

func main() {
	go handleQuit()
	mode = flag.String("mode", QUERY_MODE, "idk")
	flag.Parse()
	src.EngineStart()
	switch *mode {
	case QUERY_MODE:
		src.StartHandlingQueries()
	case INDEX_MODE:
		initialLink := os.Args[len(os.Args)-1]
		initialLink = strings.ReplaceAll(initialLink, `'`, "")
		src.StartCrawlingAndIndexing(initialLink)
		src.DoneIndexing()
	}
}
