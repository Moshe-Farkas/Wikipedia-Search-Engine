package main

import (
	"fmt"
	"indexer/src"
	"strings"
	"regexp"
	"time"
)

func main() {
	start := time.Now()
	tokenizerTesting()
	fmt.Printf("finished in %v\n", time.Since(start))
}

func tokenizerTesting() {
	src.Idk()
}

func tokenized(input string) []string {
	tokens := []string {}
	for _, temp_token := range strings.Split(input, " ") {
		token := regexp.MustCompile("[^a-zA-Z]+").ReplaceAllString(temp_token, "")
		tokens = append(tokens, strings.ToLower(token))
	}
	return tokens
}

