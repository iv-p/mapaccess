package mapaccess

import (
	"testing"
)

type lexTest struct {
	name  string
	input string
	items []item
}

func mkItem(typ itemType, text string) item {
	return item{
		typ: typ,
		val: text,
	}
}

var (
	iDot   = mkItem(itemDot, ".")
	iEOF   = mkItem(itemEOF, "")
	iError = mkItem(itemError, "")
)

var lexTests = []lexTest{
	{"empty", "", []item{iEOF}},
	{"alpha key", "somekey", []item{mkItem(itemIdentifier, "somekey"), iEOF}},
	{"alphanum key", "some6key9", []item{mkItem(itemIdentifier, "some6key9"), iEOF}},
	{"multiple keys", "keyone.keytwo", []item{mkItem(itemIdentifier, "keyone"), iDot, mkItem(itemIdentifier, "keytwo"), iEOF}},
	{"array", "keyone[0]", []item{mkItem(itemIdentifier, "keyone"), mkItem(itemArrayIndex, "[0]"), iEOF}},
	{"nested array", "keyone.keytwo[0]", []item{mkItem(itemIdentifier, "keyone"), iDot, mkItem(itemIdentifier, "keytwo"), mkItem(itemArrayIndex, "[0]"), iEOF}},
	{"nested array", "keyone[9].keytwo[0]", []item{mkItem(itemIdentifier, "keyone"), mkItem(itemArrayIndex, "[9]"), iDot, mkItem(itemIdentifier, "keytwo"), mkItem(itemArrayIndex, "[0]"), iEOF}},
	{"root array", "[0].test", []item{mkItem(itemArrayIndex, "[0]"), iDot, mkItem(itemIdentifier, "test"), iEOF}},
	{"root array", "[0][1].test", []item{mkItem(itemArrayIndex, "[0]"), mkItem(itemArrayIndex, "[1]"), iDot, mkItem(itemIdentifier, "test"), iEOF}},
	{"error", "keyo..", []item{mkItem(itemIdentifier, "keyo"), iDot, iError}},
	{"error", "keyo.[0].test", []item{mkItem(itemIdentifier, "keyo"), iDot, mkItem(itemArrayIndex, "[0]"), iDot, mkItem(itemIdentifier, "test"), iEOF}},
	{"error", ".[0].test", []item{iError}},
	{"error", " somekey", []item{iError}},
	{"error", "somekey ", []item{mkItem(itemIdentifier, "somekey"), iError}},
	{"error", "somekey.", []item{mkItem(itemIdentifier, "somekey"), iDot, iEOF}},
}

// collect gathers the emitted items into a slice.
func collect(t *lexTest) (items []item) {
	l := lex(t.input)
	for {
		item := l.nextItem()
		items = append(items, item)
		if item.typ == itemEOF || item.typ == itemError {
			break
		}
	}
	return
}

func equal(i1, i2 []item, checkPos bool) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			return false
		}
		if i1[k].val != i2[k].val && i1[k].typ != itemError {
			return false
		}
	}
	return true
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		items := collect(&test)
		if !equal(items, test.items, false) {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%v", test.name, items, test.items)
		}
	}
}
