package mapaccess

import (
	"testing"
)

type parseTest struct {
	input  string
	tokens []token
}

func mkToken(typ tokenType, text string) token {
	return token{
		typ: typ,
		val: text,
	}
}

var (
	tEnd   = mkToken(tokenEnd, "")
	tError = mkToken(tokenError, "")
)

var parseTests = []parseTest{
	{"", []token{tEnd}},
	{"somekey", []token{mkToken(tokenIdentifier, "somekey"), tEnd}},
	{"some6key9", []token{mkToken(tokenIdentifier, "some6key9"), tEnd}},
	{"keyone.keytwo", []token{mkToken(tokenIdentifier, "keyone"), mkToken(tokenIdentifier, "keytwo"), tEnd}},
	{"keyone[0]", []token{mkToken(tokenIdentifier, "keyone"), mkToken(tokenArrayIndex, "0"), tEnd}},
	{"keyone.keytwo[0]", []token{mkToken(tokenIdentifier, "keyone"), mkToken(tokenIdentifier, "keytwo"), mkToken(tokenArrayIndex, "0"), tEnd}},
	{"keyone[9].keytwo[0]", []token{mkToken(tokenIdentifier, "keyone"), mkToken(tokenArrayIndex, "9"), mkToken(tokenIdentifier, "keytwo"), mkToken(tokenArrayIndex, "0"), tEnd}},
	{"[0].test", []token{mkToken(tokenArrayIndex, "0"), mkToken(tokenIdentifier, "test"), tEnd}},
	{"keyo..", []token{mkToken(tokenIdentifier, "keyo"), tError}},
	{"keyo.[0].test", []token{mkToken(tokenIdentifier, "keyo"), tError}},
	{" somekey", []token{tError}},
	{"somekey ", []token{mkToken(tokenIdentifier, "somekey"), tError}},
	{"somekey.", []token{mkToken(tokenIdentifier, "somekey"), tEnd}},
}

// collectTokens gathers the emitted items into a slice.
func collectTokens(t *parseTest) (tokens []token) {
	p := parse(t.input)
	for {
		token := p.nextItem()
		tokens = append(tokens, token)
		if token.typ == tokenEnd || token.typ == tokenError {
			break
		}
	}
	return
}

func tokenEqual(i1, i2 []token) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			return false
		}
		if i1[k].val != i2[k].val && i1[k].typ != tokenError {
			return false
		}
	}
	return true
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		tokens := collectTokens(&test)
		if !tokenEqual(tokens, test.tokens) {
			t.Errorf("got\n\t%+v\nexpected\n\t%v", tokens, test.tokens)
		}
	}
}
