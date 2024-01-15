package jparser

import (
	"fmt"
)

//go:generate stringer -type=TokenType
type TokenType int

// The list of tokens
const (
	TokenIllegal TokenType = iota
	TokenEof

	TokenString // a valid string literal
	TokenNumber // a valid json number
	TokenTrue   // true
	TokenFalse  // false
	TokenNull   // null

	TokenLbrace // {
	TokenLbrack // [
	TokenRbrace // }
	TokenRbrack // ]
	TokenComma  // ,
	TokenColon  // :
)

// Token represents a lexical token of json.
type Token struct {
	Type  TokenType
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("<%s %s>", t.Type, t.Value)
}

const eof byte = 0

type lexer struct {
	input  []byte
	start  int
	pos    int
	tokens chan Token
}

func newLexer(input []byte) *lexer {
	return &lexer{
		input:  input,
		tokens: make(chan Token),
	}
}

// Lex reads the input and returns a channel of tokens.
func Lex(input []byte) chan Token {
	l := newLexer(input)
	go l.run()
	return l.tokens
}

func (l *lexer) next() byte {
	if l.pos >= len(l.input) {
		return eof
	}
	r := l.input[l.pos]
	l.pos++
	return r
}

func (l *lexer) prev() {
	l.pos--
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) emit(typ TokenType) {
	l.tokens <- Token{Type: typ, Value: string(l.input[l.start:l.pos])}
	l.start = l.pos
}

func (l *lexer) run() {
	for state := lex; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func (l *lexer) errorf(format string, args ...interface{}) stateFunc {
	errToken := Token{
		Type:  TokenIllegal,
		Value: fmt.Sprintf(format, args...),
	}
	l.tokens <- errToken
	return nil
}

type stateFunc func(*lexer) stateFunc

func lex(l *lexer) stateFunc {
	b := l.next()
	if isDigit(b) || (b == '-' && isDigit(l.next())) {
		return lexNumber
	} else if isWhiteSpace(b) {
		l.ignore()
		return lex
	}

	switch b {
	case eof:
		l.emit(TokenEof)
		return nil
	case 't':
		return lexValue("true", TokenTrue)
	case 'f':
		return lexValue("false", TokenFalse)
	case 'n':
		return lexValue("null", TokenNull)
	case '"':
		return lexString
	case '{':
		l.emit(TokenLbrace)
	case '}':
		l.emit(TokenRbrace)
	case '[':
		l.emit(TokenLbrack)
	case ']':
		l.emit(TokenRbrack)
	case ':':
		l.emit(TokenColon)
	case ',':
		l.emit(TokenComma)
	default:
		return l.errorf("unexpected character %c", b)
	}
	return lex
}

func lexValue(value string, token TokenType) stateFunc {
	return func(l *lexer) stateFunc {
		for idx, b := range []byte(value) {
			if idx == 0 {
				continue
			}
			if n := l.next(); n != b {
				if n == eof {
					unt := string(l.input[l.start:])
					return l.errorf("unterminated value %s, expected %s", unt, value)
				}
				return l.errorf("unexpected character %c, expected %s", n, value)
			}
		}
		l.emit(token)
		return lex
	}
}

func lexNumber(l *lexer) stateFunc {
	integer, fraction := true, false
	for {
		b := l.next()
		if b == '.' {
			if integer {
				integer = false
				fraction = true
				continue
			} else {
				return l.errorf("unexpected . in number")
			}
		}
		if (integer || fraction) && isLetter(b) {
			if b == 'E' || b == 'e' {
				if !isSign(l.next()) {
					l.prev()
				}
				integer = false
				fraction = false
				continue
			} else {
				return l.errorf("unexpected character %c in number", b)
			}
		}
		if !isDigit(b) {
			if b != eof {
				l.prev()
			}
			break
		}
	}
	l.emit(TokenNumber)
	return lex
}

func lexString(l *lexer) stateFunc {
	for {
		b := l.next()
		if b == '\\' {
			switch n := l.next(); n {
			case '"', '\\', 'b', 'f', 'n', 'r', 't':
			case 'u':
				for i := 0; i < 4; i++ {
					if h := l.next(); !isHex(h) {
						if h == eof {
							lit := string(l.input[l.start:])
							return l.errorf("unterminated string literal %s", lit)
						}
						return l.errorf("unexpected non-hex character %c", h)
					}
				}
			default:
				return l.errorf("unexpected escape character %c", n)
			}
		} else if b == '"' {
			break
		} else if b == eof {
			lit := string(l.input[l.start:])
			return l.errorf("unterminated string literal %s", lit)
		}
	}
	l.emit(TokenString)
	return lex
}

func isDigit(r byte) bool { return r >= '0' && r <= '9' }

func isSign(r byte) bool { return r == '+' || r == '-' }

func isLetter(r byte) bool { return (r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z') }

func isHex(r byte) bool {
	return isDigit(r) || r >= 'A' && r <= 'F' || r >= 'a' && r <= 'f'
}

func isWhiteSpace(r byte) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}
