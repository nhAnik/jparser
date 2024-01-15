package jparser

import (
	"errors"
	"fmt"
)

type (
	Value interface {
		value()
	}

	// Json represents the root element.
	Json struct {
		Value Value
	}

	// Object represents a json object like {"abc": "xyz"}
	Object struct {
		Elements []Element
	}

	// Array represents a json array like [1, "foo", true]
	Array struct {
		Values []Value
	}

	// Literal represents any valid json literal like string,
	// number and special value like true, false and null.
	Literal struct {
		Value string
		Type  TokenType
	}

	// Element represents a key value pair inside a json Object.
	Element struct {
		Key   string
		Value Value
	}
)

func (*Object) value()  {}
func (*Array) value()   {}
func (*Literal) value() {}

type parser struct {
	tokens chan Token
	cur    Token
	err    error
}

// Parse parses the given input and returns an AST representation
// of the json. In case of error, it returns an error.
func Parse(input []byte) (*Json, error) {
	p := parser{
		tokens: Lex(input),
	}
	j := p.parse()
	if p.err != nil {
		return nil, p.err
	}
	// check if all the tokens have been consumed
	if p.cur.Type != TokenEof {
		p.addErr(fmt.Errorf("expected TokenEof but found %s", p.cur.Type))
		return nil, p.err
	}
	return j, nil
}

func (p *parser) parseValue() Value {
	if p.err != nil {
		return nil
	}

	switch token := p.cur; token.Type {
	case TokenLbrace:
		return p.parseObject()
	case TokenLbrack:
		return p.parseArray()
	case TokenString, TokenNumber, TokenTrue, TokenFalse, TokenNull:
		return p.parseLiteral()
	default:
		p.addErr(fmt.Errorf("unexpected token %s", token.Value))
		return nil
	}
}

func (p *parser) parseObject() *Object {
	if p.err != nil {
		return nil
	}

	var elements []Element
	p.expect(TokenLbrace)

	if p.cur.Type != TokenRbrace {
		for {
			lit := p.parseLiteral()
			if lit == nil {
				return nil
			}

			var key string
			if lit.Type == TokenString {
				key = lit.Value
			} else {
				p.addErr(fmt.Errorf("expected string key but found %s", lit.Value))
				return nil
			}
			p.expect(TokenColon)
			if value := p.parseValue(); value != nil {
				elements = append(elements, Element{Key: key, Value: value})
			} else {
				return nil
			}

			if p.cur.Type == TokenRbrace {
				p.next()
				break
			}
			p.expect(TokenComma)
		}
	}

	return &Object{
		Elements: elements,
	}
}

func (p *parser) parseArray() *Array {
	if p.err != nil {
		return nil
	}

	var values []Value
	p.expect(TokenLbrack)

	if p.cur.Type != TokenRbrack {
		for {
			if value := p.parseValue(); value != nil {
				values = append(values, value)
			} else {
				return nil
			}

			if p.cur.Type == TokenRbrack {
				p.next()
				break
			}
			p.expect(TokenComma)
		}
	}

	return &Array{
		Values: values,
	}
}

func (p *parser) parseLiteral() *Literal {
	if p.err != nil {
		return nil
	}

	lit := &Literal{
		Value: p.cur.Value,
		Type:  p.cur.Type,
	}
	p.next()
	return lit
}

func (p *parser) next() {
	p.cur = <-p.tokens
	if p.cur.Type == TokenIllegal {
		p.addErr(errors.New(p.cur.Value))
	}
}

func (p *parser) expect(token TokenType) {
	if p.cur.Type == TokenIllegal {
		p.addErr(errors.New(p.cur.Value))
		return
	}
	if p.cur.Type != token {
		p.addErr(fmt.Errorf("expected %s but found %s", token, p.cur.Type))
	}
	p.next()
}

func (p *parser) addErr(err error) {
	if p.err == nil {
		p.err = err
	}
}

func (p *parser) parse() *Json {
	p.next()
	if p.err != nil {
		return nil
	}
	return &Json{
		Value: p.parseValue(),
	}
}
