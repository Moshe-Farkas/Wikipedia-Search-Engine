package main

import (
	"fmt"
	"os"
	"project/src"
	"flag"
)

const (
	QUIT_MESSAGE = "q"
	QUERY_MODE = "query"
	INDEX_MODE = "index"
)

func handleQuit() {
	var temp string = ""
	for temp != QUIT_MESSAGE {
		fmt.Scanf("%s", &temp)
	}
	src.Cleanup()
	os.Exit(0)
}

func main() {
	go handleQuit()
	// src.OnStart()
	// src.StartCrawlingAndIndexing()
	var mode = flag.String("mode", QUERY_MODE, "idk")
	flag.Parse()
	src.EngineStart()
	switch *mode {
	case QUERY_MODE:
		// do query stuff
		src.StartHandlingQueries()
	case INDEX_MODE:
		// do index stuff
	}
}