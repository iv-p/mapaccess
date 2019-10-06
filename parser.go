package mapaccess

import (
	"fmt"
	"strconv"
	"strings"
)

type token struct {
	typ tokenType
	val string
}

type parser struct {
	tokens chan token // channel of parsed token
	items  chan item  // chan to send tokens for client
	buf    *item      // have a buffer of 1 item for parser
	lex    *lexer     // input lexer
}

type tokenType int

const (
	tokenError tokenType = iota
	tokenEnd
	tokenIdentifier
	tokenArrayIndex
)

type parseStateFn func(*parser) parseStateFn

func parse(input string) *parser {
	p := &parser{
		tokens: make(chan token),
		lex:    lex(input),
	}
	go p.run()
	return p
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state.
func (p *parser) errorf(format string, args ...interface{}) parseStateFn {
	p.tokens <- token{tokenError, fmt.Sprintf(format, args...)}
	return nil
}

// emit passes a token back to the client.
func (p *parser) emit(t token) {
	p.tokens <- t
}

// nextItem returns the next item when it becomes available
func (p *parser) nextItem() token {
	return <-p.tokens
}

func (p *parser) run() {
	defer close(p.tokens)
	for state := parseStart; state != nil; {
		state = state(p)
	}
}

func (p *parser) getBufOrNext() item {
	if p.buf != nil {
		item := *p.buf
		p.buf = nil
		return item
	}
	return p.lex.nextItem()
}

// parseStart scans for either an identifier, or an array index
func parseStart(p *parser) parseStateFn {
	item := p.getBufOrNext()
	p.buf = &item
	switch item.typ {
	case itemArrayIndex:
		return parseArrayIndex
	case itemIdentifier:
		return parseIdentifier
	case itemEOF:
		p.emit(token{tokenEnd, ""})
		return nil
	default:
		return p.errorf("expected array index or identifier")
	}
}

// parseIdentifier scans for identifiers
func parseIdentifier(p *parser) parseStateFn {
	item := p.getBufOrNext()
	if item.typ == itemIdentifier {
		// we already did rune checking in the lexer, good to go
		p.emit(token{tokenIdentifier, item.val})
		next := p.lex.nextItem()
		p.buf = &next
		switch next.typ {
		case itemDot:
			return parseDot
		case itemArrayIndex:
			return parseArrayIndex
		case itemEOF:
			p.emit(token{tokenEnd, ""})
			return nil
		default:
			return p.errorf("expected dot or array index")
		}
	}

	if item.typ == itemEOF {
		p.emit(token{tokenEnd, ""})
		return nil
	}

	return p.errorf("expected identifier")
}

// parseDot scans for dots
func parseDot(p *parser) parseStateFn {
	item := p.getBufOrNext()
	if item.typ != itemDot {
		return p.errorf("expected dot")
	}
	// do nothing, ingest dot
	return parseIdentifier
}

// parseArrayIndex scans for dots
func parseArrayIndex(p *parser) parseStateFn {
	item := p.getBufOrNext()
	if item.typ != itemArrayIndex {
		return p.errorf("expected array index")
	}
	// lexer already checked that the val is starting and ending with brackets []
	index := strings.Trim(item.val, "[]")
	if _, err := strconv.Atoi(index); err != nil {
		return p.errorf("expected a integer")
	}
	p.emit(token{tokenArrayIndex, index})
	i := p.lex.nextItem()
	p.buf = &i
	switch p.buf.typ {
	case itemDot:
		return parseDot
	case itemArrayIndex:
		return parseArrayIndex
	case itemEOF:
		p.emit(token{tokenEnd, ""})
		return nil
	default:
		return p.errorf("expected dot or array index after array index")
	}
}
