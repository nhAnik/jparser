package jparser

import "testing"

func TestParse(t *testing.T) {
	j, err := Parse([]byte(`{"foo": 123, "bar": [345, "pqr", true], "abc": "xyz"}`))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
		return
	}
	obj, ok := j.Value.(*Object)
	if !ok {
		t.Fatalf("Expected *Object but got %T", j.Value)
	}

	checkKeys(t, obj, []string{`"foo"`, `"bar"`, `"abc"`})

	val := obj.Elements[1].Value
	arr, ok := val.(*Array)
	if !ok {
		t.Fatalf("Expected *Array but got %T", val)
	}

	assert(t, 3, len(arr.Values))

	fst := arr.Values[0]
	if lit, ok := fst.(*Literal); ok {
		assert(t, "345", lit.Value)
	} else {
		t.Fatalf("Expected *Literal but got %T", fst)
	}

	snd := arr.Values[1]
	if lit, ok := snd.(*Literal); ok {
		assert(t, `"pqr"`, lit.Value)
	} else {
		t.Fatalf("Expected *Literal but got %T", snd)
	}

	trd := arr.Values[2]
	if lit, ok := trd.(*Literal); ok {
		assert(t, "true", lit.Value)
	} else {
		t.Fatalf("Expected *Literal but got %T", trd)
	}
}

func assert[T comparable](t *testing.T, expected T, got T) {
	if expected != got {
		t.Fatalf("Expected: %v, got: %v", expected, got)
	}
}

func checkKeys(t *testing.T, obj *Object, expected []string) {
	var got []string
	for _, el := range obj.Elements {
		got = append(got, el.Key)
	}
	if !isMatched(got, expected) {
		t.Fatalf("Expected: %v, got: %v", expected, got)
	}
}

func TestParseError(t *testing.T) {
	testCases := []struct{ input, want string }{
		// Test cases with lex error
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

		// Test cases with parse error
		{`{"foo", 123)`, "expected TokenColon but found TokenComma"},
		{`{"foo": 123]`, "expected TokenComma but found TokenRbrack"},
		{`{abc: 123}`, "unexpected character a"},
		{`{"foo": 123)`, "unexpected character )"},
		{`{true: 123}`, "expected string key but found true"},
		{`{123: 456}`, "expected string key but found 123"},
		{`[true: false]`, "expected TokenComma but found TokenColon"},
		{`[true: false}`, "expected TokenComma but found TokenColon"},
		{`[true, false}`, "expected TokenComma but found TokenRbrace"},
		{`{"foo": [true, false}}`, "expected TokenComma but found TokenRbrace"},
		{
			`{"foo": 123, "bar": [345, "pqr", true], "abc", "xyz"}`,
			"expected TokenColon but found TokenComma",
		},
		{
			`{"foo": 123} "bar": [345, "pqr", true], "abc": "xyz"}`,
			"expected TokenEof but found TokenString",
		},
		{`[1, 2, 3][1, 2, 3]`, "expected TokenEof but found TokenLbrack"},
	}

	for _, kase := range testCases {
		_, err := Parse([]byte(kase.input))
		if err == nil {
			t.Fatalf("Expected error but got no error")
		}
		if kase.want != err.Error() {
			t.Fatalf("Expected error: %s, got: %v", kase.want, err)
		}
	}
}
