package mapaccess

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type item struct {
	typ itemType // The type of this item.
	val string   // The key value of this item.
}

// itemType identifies the type of lex items.
type itemType int

const (
	itemError itemType = iota
	itemEOF
	itemIdentifier
	itemArrayIndex
	itemDot
)

const eof = -1

// lexStateFn represents the state of the scanner as a function that returns the next state.
type lexStateFn func(*lexer) lexStateFn

// lexer holds the state of the scanner.
type lexer struct {
	input string    // the string being scanned
	pos   int       // current position in the input
	start int       // start position of this item
	width int       // width of last rune read from input
	items chan item // channel of scanned items
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) lexStateFn {
	l.items <- item{itemError, fmt.Sprintf(format, args...)}
	return nil
}

// nextItem returns the next item from the input.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) nextItem() item {
	return <-l.items
}

// lex creates a new scanner for the input string.
func lex(input string) *lexer {
	l := &lexer{
		input: input,
		items: make(chan item),
	}
	go l.run()
	return l
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	defer close(l.items)
	for state := lexStart; state != nil; {
		state = state(l)
	}
}

func lexStart(l *lexer) lexStateFn {
	next := l.peek()
	switch next {
	case '[':
		return lexArrayIndexAction
	default:
		return lexIdentifier
	}
}

// lexIdentifier scans an alphanumeric.
func lexIdentifier(l *lexer) lexStateFn {
	for {
		switch r := l.next(); {
		case isValidIdentifierRune(r):
			// absorb.
		default:
			l.backup()
			next := l.peek()
			if l.pos == l.start {
				if next != -1 {
					return l.errorf("expected identifier")
				}
				l.emit(itemEOF)
				return nil
			}
			l.emit(itemIdentifier)
			switch next {
			case '.':
				return lexDotAction
			case '[':
				return lexArrayIndexAction
			case -1:
				l.emit(itemEOF)
				return nil
			default:
				return l.errorf("expected <.> or <array index>")
			}
		}
	}
}

// lexDotAction expects a dot.
func lexDotAction(l *lexer) lexStateFn {
	r := l.next()
	switch r {
	case '.':
		l.emit(itemDot)
		next := l.peek()
		switch next {
		case '[':
			return lexArrayIndexAction
		default:
			return lexIdentifier
		}
	default:
		return l.errorf("bad character")
	}
}

// lexArrayIndexAction scans for [\d+].
func lexArrayIndexAction(l *lexer) lexStateFn {
	if !l.accept("[") {
		return l.errorf("missing closing bracket [ at array index <[]>")
	}
	if !l.accept("0123456789") {
		return l.errorf("missing digits in array index <[]>")
	}
	digits := "0123456789"
	l.acceptRun(digits)

	if !l.accept("]") {
		return l.errorf("missing closing bracket ] at array index <[]>")
	}
	l.emit(itemArrayIndex)
	next := l.peek()
	switch next {
	case '.':
		return lexDotAction
	case '[':
		return lexArrayIndexAction
	default:
		return lexEOF
	}
}

// lexEOF scans for EOF.
func lexEOF(l *lexer) lexStateFn {
	if l.next() == -1 {
		l.emit(itemEOF)
		return nil
	}
	return lexIdentifier
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isValidIdentifierRune(r rune) bool {
	return strings.ContainsRune("_-", r) || unicode.IsLetter(r) || unicode.IsDigit(r)
}
