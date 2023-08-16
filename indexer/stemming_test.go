package indexer

import (
	"reflect"
	"testing"
)

func TestStep1A(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"caresses", "caress"},
		{"caress", "caress"},
		{"cares", "care"},
		{"care", "care"},
	}
	for _, tt := range tests {
		got := Stem(tt.input)
		if got != tt.want {
			t.Errorf("want %s but got %s\n", tt.want, got)
		}
	}
}

func TestIsConsonant(t *testing.T) {
	enumerateConsonants := func(str string) []rune {
		consonants := make([]rune, 0)
		for i, char := range str {
			if isConsonant(str, i) {
				consonants = append(consonants, char)
			}
		}
		return consonants
	}
	tests := []struct {
		input string
		want  []rune
	}{
		{"toy", []rune{'t', 'y'}},
		{"syzygy", []rune{'s', 'z', 'g'}},
		{"y", []rune{'y'}},
		{"ny", []rune{'n'}},
		{"yoy", []rune{'y', 'y'}},
		{"yyyy", []rune{'y'}},
		{"cyyy", []rune{'c'}},
	}
	for _, test := range tests {
		got := enumerateConsonants(test.input)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("for input %s want %c but got %c\n", test.input, test.want, got)
		}
	}
}

func TestVcCount(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"tree", 0},
		{"trouble", 1},
		{"troubles", 2},
		{"multidimensional", 6},
		{"allow", 2},
	}
	for _, test := range tests {
		got := vcCount(test.input)
		if test.want != got {
			t.Errorf("for %s want %d but got %d\n", test.input, test.want, got)
		}
	}
}

func TestStep1B(t *testing.T) {
	tests := []struct {
		input string
		want string
	} {
		{"tanned", "tan"},
		{"fizzed", "fizz"},
		{"filing", "file"},
		{"bled", "bled"},
		{"sing", "sing"},
		{"feed", "feed"},
		{"ifeedeed", "ifeedee"},
		{"agreed", "agree"},
		{"plastered", "plaster"},
		{"motoring", "motor"},
		{"hopping", "hop"},
		{"conflated", "conflate"},
		{"failing", "fail"},
		{"hissing", "hiss"},
	}
	for _, test := range tests {
		got := step1B(test.input)
		if test.want != got {
			t.Errorf("for %s want %s but got %s", test.input, test.want, got)
		}
	}
}

func TestStep1C(t *testing.T) {
	tests := []struct {
		input string
		want string
	} {
		{"happy", "happi"},
		{"sky", "sky"},
	}
	for _, test := range tests {
		got := step1C(test.input)
		if test.want != got {
			t.Errorf("for %s want %s but got %s", test.input, test.want, got)
		}
	}
}

func TestStep2(t *testing.T) {
	tests := []struct {
		input string
		want string
	} {
		{"conditional", "condition"},
		{"rational", "rational"},
		{"valenci", "valence"},
		{"hesitanci", "hesitance"},
		{"digitizer", "digitize"},
		{"conformabli", "conformable"},
		{"radicalli", "radical"},
		{"differentli", "different"},
		{"vileli", "vile"},
		{"analogousli", "analogous"},
		{"vietnamization", "vietnamize"},
		{"predication", "predicate"},
		{"operator", "operate"},
		{"feudalism", "feudal"},
		{"decisiveness", "decisive"},
		{"hopefulness", "hopeful"},
		{"callousness", "callous"},
		{"formaliti", "formal"},
		{"sensitiviti", "sensitive"},
		{"sensibiliti", "sensible"},
	}

	for _, test := range tests {
		got := step2(test.input)
		if test.want != got {
			t.Errorf("for %s want %s but got %s", test.input, test.want, got)
		}
	}
}

func TestStep3(t *testing.T) {
	tests := []struct {
		input string
		want string
	} {
		{"triplicate", "triplic"},
		{"formative", "form"},
		{"formalize", "formal"},
		{"electriciti", "electric"},
		{"electrical", "electric"},
		{"hopeful", "hope"},
		{"goodness", "good"},                        
	}
	for _, test := range tests {
		got := step3(test.input)
		if test.want != got {
			t.Errorf("for %s want %s but got %s", test.input, test.want, got)
		}
	}
}

func TestStep4(t *testing.T) {
	tests := []struct {
		input string
		want string
	} {
		{"adoption", "adopt"},
		{"revival", "reviv"},
		{"allowance", "allow"},
		{"inference", "infer"},
		{"airliner", "airlin"},
		{"gyroscopic", "gyroscop"},
		{"adjustable", "adjust"},
		{"defensible", "defens"},
		{"irritant", "irrit"},
		{"replacement", "replac"},
		{"adjustment", "adjust"},
		{"dependent", "depend"},
		{ "homologou", "homolog"},
		{ "communism", "commun"},
		{ "activate", "activ"},
		{ "angulariti", "angular"},
		{ "homologous", "homolog"},
		{ "effective", "effect"},
		{ "bowdlerize", "bowdler"},
	}
	for _, test := range tests {
		got := step4(test.input)
		if test.want != got {
			t.Errorf("for %s want %s but got %s", test.input, test.want, got)
		}
	}
}

func TestStep5a(t *testing.T) {
	tests := []struct {
		input string
		want string
	} {
		{"probate", "probat"},
		{"rate", "rate"},
		{"cease", "ceas"},                       
	}
	for _, test := range tests {
		got := step5A(test.input)
		if test.want != got {
			t.Errorf("for %s want %s but got %s", test.input, test.want, got)
		}
	}
}

func TestStep5B(t *testing.T) {
	tests := []struct {
		input string
		want string
	} {
		{"controll", "control"},
		{"roll", "roll"},
	}
	for _, test := range tests {
		got := step5B(test.input)
		if test.want != got {
			t.Errorf("for %s want %s but got %s", test.input, test.want, got)
		}
	}
}
