package jparser

import "testing"

func TestLex(t *testing.T) {
	testCases := []struct {
		input string
		want  []Token
	}{
		{"true", []Token{{TokenTrue, "true"}}},
		{"false", []Token{{TokenFalse, "false"}}},
		{"null", []Token{{TokenNull, "null"}}},
		{`"string"`, []Token{{TokenString, `"string"`}}},
		{`""`, []Token{{TokenString, `""`}}},
		{`"ab\"cd\""`, []Token{{TokenString, `"ab\"cd\""`}}},
		{`"ab\\\""`, []Token{{TokenString, `"ab\\\""`}}},
		{`"15\u00f8C"`, []Token{{TokenString, `"15\u00f8C"`}}},
		{`"\uD83D\uDE10"`, []Token{{TokenString, `"\uD83D\uDE10"`}}},
		{"123", []Token{{TokenNumber, "123"}}},
		{"-123", []Token{{TokenNumber, "-123"}}},
		{"-1.23", []Token{{TokenNumber, "-1.23"}}},
		{"1.23e4", []Token{{TokenNumber, "1.23e4"}}},
		{"-1.23E5", []Token{{TokenNumber, "-1.23E5"}}},
		{"12E3", []Token{{TokenNumber, "12E3"}}},
		{"1.2E+3", []Token{{TokenNumber, "1.2E+3"}}},
		{"1e-3", []Token{{TokenNumber, "1e-3"}}},
		{`{"foo": 123}`,
			[]Token{
				{TokenLbrace, "{"},
				{TokenString, `"foo"`},
				{TokenColon, ":"},
				{TokenNumber, "123"},
				{TokenRbrace, "}"},
			},
		},
		{`{"foo": -1.23e4}`,
			[]Token{
				{TokenLbrace, "{"},
				{TokenString, `"foo"`},
				{TokenColon, ":"},
				{TokenNumber, "-1.23e4"},
				{TokenRbrace, "}"},
			},
		},
		{`{"foo": "bar"}`,
			[]Token{
				{TokenLbrace, "{"},
				{TokenString, `"foo"`},
				{TokenColon, ":"},
				{TokenString, `"bar"`},
				{TokenRbrace, "}"},
			},
		},
		{`{"foo": null}`,
			[]Token{
				{TokenLbrace, "{"},
				{TokenString, `"foo"`},
				{TokenColon, ":"},
				{TokenNull, "null"},
				{TokenRbrace, "}"},
			},
		},
		{`[123, true, "foo", null]`,
			[]Token{
				{TokenLbrack, "["},
				{TokenNumber, "123"},
				{TokenComma, ","},
				{TokenTrue, "true"},
				{TokenComma, ","},
				{TokenString, `"foo"`},
				{TokenComma, ","},
				{TokenNull, "null"},
				{TokenRbrack, "]"},
			},
		},
		{`[1.23, 2e3, -4.5E6, null]`,
			[]Token{
				{TokenLbrack, "["},
				{TokenNumber, "1.23"},
				{TokenComma, ","},
				{TokenNumber, "2e3"},
				{TokenComma, ","},
				{TokenNumber, "-4.5E6"},
				{TokenComma, ","},
				{TokenNull, "null"},
				{TokenRbrack, "]"},
			},
		},
		{`{"foo": [123, false]}`,
			[]Token{
				{TokenLbrace, "{"},
				{TokenString, `"foo"`},
				{TokenColon, ":"},
				{TokenLbrack, "["},
				{TokenNumber, "123"},
				{TokenComma, ","},
				{TokenFalse, "false"},
				{TokenRbrack, "]"},
				{TokenRbrace, "}"},
			},
		},
		{`{}`,
			[]Token{
				{TokenLbrace, "{"},
				{TokenRbrace, "}"},
			},
		},
		{`[]`,
			[]Token{
				{TokenLbrack, "["},
				{TokenRbrack, "]"},
			},
		},
		{`[  	]`,
			[]Token{
				{TokenLbrack, "["},
				{TokenRbrack, "]"},
			},
		},
		{`{
			"foo"
				:
					"bar"   }`,
			[]Token{
				{TokenLbrace, "{"},
				{TokenString, `"foo"`},
				{TokenColon, ":"},
				{TokenString, `"bar"`},
				{TokenRbrace, "}"},
			},
		},
	}

	for _, kase := range testCases {
		kase.want = append(kase.want, Token{TokenEof, ""})
		got := toSlice(Lex([]byte(kase.input)))
		if !isMatched(got, kase.want) {
			t.Fatalf("Expected: %v, got: %v", kase.want, got)
		}
	}
}

func toSlice(c chan Token) []Token {
	var tokens []Token
	for t := range c {
		tokens = append(tokens, t)
	}
	return tokens
}

func isMatched[T comparable](got []T, want []T) bool {
	if len(got) != len(want) {
		return false
	}
	for i := 0; i < len(got); i++ {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

func TestLexError(t *testing.T) {
	testCases := []struct{ input, want string }{
		{"True", "unexpected character T"},
		{"t", "unterminated value t, expected true"},
		{"falsee", "unexpected character e"},
		{"faLse", "unexpected character L, expected false"},
		{"nul", "unterminated value nul, expected null"},
		{`"str`, `unterminated string literal "str`},
		{`"`, `unterminated string literal "`},
		{`"ab\"`, `unterminated string literal "ab\"`},
		{`"ab\x"`, `unexpected escape character x`},
		{`"ab\x"`, `unexpected escape character x`},
		{`"\uD83D\uDE1"`, `unexpected non-hex character "`},
		{`"\u007\u0074"`, `unexpected non-hex character \`},
		{`"\u00`, `unterminated string literal "\u00`},
		{"1.2.3", "unexpected . in number"},
		{"+1.2", "unexpected character +"},
		{"-.1", "unexpected character -"},
		{"1.2f3", "unexpected character f in number"},
		{"1Y3", "unexpected character Y in number"},
		{"1%3", "unexpected character %"},
		{"1.2e*3", "unexpected character *"},
		{"1.2e3e3", "unexpected character e"},
		{`{"foo", 123)`, "unexpected character )"},
	}

	for _, kase := range testCases {
		tokens := toSlice(Lex([]byte(kase.input)))
		got := tokens[len(tokens)-1]
		if got.Type != TokenIllegal {
			t.Fatalf("Expected: TokenIllegal, got: %v", got.Type.String())
		}
		if got.Value != kase.want {
			t.Fatalf("Expected error: %v, got: %v", kase.want, got.Value)
		}
	}
}
