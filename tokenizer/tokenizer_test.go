package tokenizer

import (
	"github.com/antham/goller/parser"
	"testing"
)

func TestTokenizeLineWithAParser(t *testing.T) {
	p := parser.NewParser("whi", []string{})

	tok := NewTokenizer(p)

	tokens, err := tok.Tokenize("[2016-01-08 20:16] [ALPM] transaction started")

	if err != nil {
		t.Error("Must not throws an error")
	}

	if tokens == nil {
		t.Error("Tokens can't be nil")
	}

	if len(tokens) != 5 {
		t.Errorf("Expected length is %v got %v", 5, len(tokens))
	}

	if tokens[0].Value != "[2016-01-08" {
		t.Errorf("We should retrieve %v got %v", "[2016-01-08", tokens[0].Value)
	}

	if tokens[4].Value != "started" {
		t.Errorf("We should retrieve %v at token 4, got %v", "started", tokens[4].Value)
	}
}

func TestTokenizeALineWithLessTokensThanFirstLine(t *testing.T) {
	p := parser.NewParser("whi", []string{})

	tok := NewTokenizer(p)

	tok.Tokenize("test1 test2 test3 test4")

	_, err := tok.Tokenize("test1 test2 test3")

	if err.Error() != "Wrong parsing strategy (based on first line tokenization), got 3 tokens instead of 4\nLine : test1 test2 test3\n" {
		t.Error("We must have an error when we try to tokenize two lines with different sizes")
	}
}

func TestTokenizeALineWithMoreTokensThanFirstLine(t *testing.T) {
	p := parser.NewParser("whi", []string{})

	tok := NewTokenizer(p)

	tok.Tokenize("test1 test2 test3 test4")

	_, err := tok.Tokenize("test1 test2 test3 test4 test5")

	if err.Error() != "Wrong parsing strategy (based on first line tokenization), got 5 tokens instead of 4\nLine : test1 test2 test3 test4 test5\n" {
		t.Error("We must have an error when we try to tokenize two lines with different sizes")
	}
}
