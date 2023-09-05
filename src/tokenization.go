package src

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

type tokenizedDoc struct {
	tokens  map[string]int  // term: tf
	docLen  int
}

func (tkns *tokenizedDoc) addToken(token string) {
	_, seen := tkns.tokens[token]	
	if seen {
		tkns.tokens[token]++
	} else {
		tkns.tokens[token] = 1
	}
	tkns.docLen++
}

func newTokenizedDoc() *tokenizedDoc {
	return &tokenizedDoc {
		map[string]int {},
		0,
	}
}

func tokenize(docText string) *tokenizedDoc {
	var tkns = newTokenizedDoc()
	re := regexp.MustCompile(`[-\s_]`)
	for _, word := range re.Split(docText, -1) {
		token, err := normalizeWord(word)
		if err == nil {
			tkns.addToken(token)
		}
	}
	return tkns
}

func normalizeWord(input string) (string, error) {
	if isStopWord(input) {
		return "", errors.New("idk")
	}
	for _, char := range input {
		if !unicode.IsPrint(char) {
			return "", errors.New("idk")
		}
	}
	var maxWordLength int = 15
	input = expandContraction(input)
	input = replaceCertainChars(input)
	input = strings.ToLower(input)
	if len(input) >= maxWordLength {
		return "", errors.New("idk")
	}
	return stem(input), nil
}

func isStopWord(str string) bool {
	stopWords := []string {
		"a", "the", "and", "going",
		"for", "able", "by", "if",
		"to", "be", "or", "not",
		"is",
	}
	for _, sw := range stopWords {
		if str == sw {
			return true
		}
	}
	return false
}

func expandContraction(input string) string {
	toStrip := []string{
		"n't",
		"'s",
		"'ll",
		"'d",
		"'ve",
	}
	for _, toMatch := range toStrip {
		input = strings.TrimSuffix(input, toMatch)
	}
	return input
}

func replaceCertainChars(input string) string {
	// return input
	re := regexp.MustCompile("[^a-zA-Z]+")
	return re.ReplaceAllString(input, "")
}

func stem(str string) string {
	// source: https://tartarus.org/martin/PorterStemmer/index.html
	if len(str) <= 2 {
		return str
	}
	str = step5B(step5A(step4(step3(step2(step1C((step1B(step1A(str)))))))))
	return str
}

func step1A(str string) string {
	if strings.HasSuffix(str, "sses") {
		return strings.TrimSuffix(str, "sses") + "ss"
	} else if strings.HasSuffix(str, "ies") {
		return strings.TrimSuffix(str, "ies") + "i"
	} else if strings.HasSuffix(str, "ss") {
		return str
	}
	return strings.TrimSuffix(str, "s") // only trims if contains
}

func step1B(str string) string {
	if strings.HasSuffix(str, "eed") {
		withoutSuffix := strings.TrimSuffix(str, "eed")
		if vcCount(withoutSuffix) > 0 {
			return withoutSuffix + "ee"
		}
		return str
	}
	var firstOrSecondStepSuccessfull bool
	if strings.HasSuffix(str, "ed") {
		withoutSuffix := strings.TrimSuffix(str, "ed")
		if containsVowel(withoutSuffix) {
			str = withoutSuffix
			firstOrSecondStepSuccessfull = true
		} else {
			firstOrSecondStepSuccessfull = false
		}
	} else if strings.HasSuffix(str, "ing") {
		withoutSuffix := strings.TrimSuffix(str, "ing")
		if containsVowel(withoutSuffix) {
			str = withoutSuffix
			firstOrSecondStepSuccessfull = true
		} else {
			firstOrSecondStepSuccessfull = false
		}
	}
	if !firstOrSecondStepSuccessfull {
		return str
	}

	if strings.HasSuffix(str, "at") {
		return strings.TrimSuffix(str, "at") + "ate"
	}
	if strings.HasSuffix(str, "bl") {
		return strings.TrimSuffix(str, "bl") + "ble"
	}
	if strings.HasSuffix(str, "iz") {
		return strings.TrimSuffix(str, "iz") + "ize"
	}

	if conditionD(str) {
		if !(str[len(str)-1] == 'l' || str[len(str)-1] == 's' || str[len(str)-1] == 'z') {
			return str[:len(str)-1]
		}
	}
	// (m=1 and *o) -> E
	if vcCount(str) == 1 && conditionO(str) {
		return str + "e"
	}
	return str
}

func step1C(str string) string {
	if strings.HasSuffix(str, "y") {
		withoutSuffix := strings.TrimSuffix(str, "y")
		if containsVowel(withoutSuffix) {
			return withoutSuffix + "i"
		}
	}
	return str
}

func step2(str string) string {
	type replaceTest struct {
		suffix      string
		replacement string
	}
	tests := []replaceTest{
		{"ational", "ate"},
		{"iveness", "ive"},
		{"fulness", "ful"},
		{"ousness", "ous"},
		{"ization", "ize"},
		{"tional", "tion"},
		{"biliti", "ble"},
		{"iviti", "ive"},
		{"entli", "ent"},
		{"ousli", "ous"},
		{"ation", "ate"},
		{"alism", "al"},
		{"aliti", "al"},
		{"enci", "ence"},
		{"anci", "ance"},
		{"izer", "ize"},
		{"abli", "able"},
		{"alli", "al"},
		{"ator", "ate"},
		{"eli", "e"},
	}
	for _, t := range tests {
		if strings.HasSuffix(str, t.suffix) {
			withoutSuffix := strings.TrimSuffix(str, t.suffix)
			if vcCount(withoutSuffix) > 0 {
				return withoutSuffix + t.replacement
			}
		}
	}
	return str
}

func step3(str string) string {
	type replaceTest struct {
		suffix      string
		replacement string
	}
	tests := []replaceTest{
		{"icate", "ic"},
		{"ative", ""},
		{"alize", "al"},
		{"iciti", "ic"},
		{"ness", ""},
		{"ical", "ic"},
		{"ful", ""},
	}
	for _, t := range tests {
		if strings.HasSuffix(str, t.suffix) {
			withoutSuffix := strings.TrimSuffix(str, t.suffix)
			if vcCount(withoutSuffix) > 0 {
				return withoutSuffix + t.replacement
			}
		}
	}
	return str
}

func step4(str string) string {
	suffexTests := []string{
		"al",
		"ance",
		"ence",
		"er",
		"ic",
		"able",
		"ible",
		"ant",
		"ement",
		"ment",
		"ent",
		"ou",
		"ism",
		"ate",
		"iti",
		"ous",
		"ive",
		"ize",
	}
	for _, t := range suffexTests {
		if strings.HasSuffix(str, t) {
			withoutSuffix := strings.TrimSuffix(str, t)
			if vcCount(withoutSuffix) > 1 {
				return withoutSuffix
			}
		}
	}
	if strings.HasSuffix(str, "ion") {
		withoutSuffix := strings.TrimSuffix(str, "ion")
		if vcCount(withoutSuffix) > 1 {
			if withoutSuffix[len(withoutSuffix)-1] == 's' || withoutSuffix[len(withoutSuffix)-1] == 't' {
				return withoutSuffix
			}
		}
	}
	return str
}

func step5A(str string) string {
	if strings.HasSuffix(str, "e") {
		withoutSuffix := strings.TrimSuffix(str, "e")
		if vcCount(withoutSuffix) > 1 {
			return withoutSuffix
		}
		if vcCount(withoutSuffix) == 1 && !conditionO(withoutSuffix) {
			return withoutSuffix
		}
	}
	return str
}

func step5B(str string) string {
	if vcCount(str) > 1 && conditionD(str) && strings.HasSuffix(str, "l") {
		return strings.TrimSuffix(str, "l")
	}
	return str
}

func isVowel(char rune) bool {
	switch char {
	case 'a':
		return true
	case 'e':
		return true
	case 'i':
		return true
	case 'o':
		return true
	case 'u':
		return true
	}
	return false
}

func containsVowel(str string) bool {
	for _, char := range str {
		if isVowel(char) {
			return true
		}
	}
	return false
}

func isConsonant(str string, charIndex int) bool {
	// A consonant is a letter other than the vowels
	// and other than a letter “Y” preceded by a consonant.
	// So in “TOY” the consonants are “T” and “Y”, and in
	// “SYZYGY” they are “S”, “Z” and “G”. - https://vijinimallawaarachchi.com/2017/05/09/porter-stemming-algorithm/
	if len(str) <= 2 {
		return false
	}
	if isVowel(rune(str[charIndex])) {
		return false
	}
	if str[charIndex] == 'y' {
		if charIndex == 0 {
			return true
		}
		return isVowel(rune(str[charIndex-1])) // supposed to be recursive? who knows...
	}
	return true
}

func vcCount(str string) int {
	// returns number of occurences of the sequence [V, C].
	// V and C mean sequence of vowel or consonant > 0. aka m
	var vcCount int
	var lastSeq rune
	if isConsonant(str, 0) {
		lastSeq = 'c'
	} else {
		lastSeq = 'v'
	}
	for i := 0; i < len(str); i++ {
		if isConsonant(str, i) {
			if lastSeq == 'v' {
				vcCount++
			}
			lastSeq = 'c'
		} else {
			lastSeq = 'v'
		}
	}
	return vcCount
}

func conditionO(str string) bool {
	if len(str) < 3 {
		return false
	}
	lastLet := rune(str[len(str)-1])
	return isConsonant(str, len(str)-1) &&
		isVowel(rune(str[len(str)-2])) &&
		isConsonant(str, len(str)-3) &&
		lastLet != 'w' &&
		lastLet != 'x' &&
		lastLet != 'y'
}

func conditionD(str string) bool {
	return isConsonant(str, len(str)-1) &&
		isConsonant(str, len(str)-2)
}
