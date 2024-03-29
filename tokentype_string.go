// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package jparser

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TokenIllegal-0]
	_ = x[TokenEof-1]
	_ = x[TokenString-2]
	_ = x[TokenNumber-3]
	_ = x[TokenTrue-4]
	_ = x[TokenFalse-5]
	_ = x[TokenNull-6]
	_ = x[TokenLbrace-7]
	_ = x[TokenLbrack-8]
	_ = x[TokenRbrace-9]
	_ = x[TokenRbrack-10]
	_ = x[TokenComma-11]
	_ = x[TokenColon-12]
}

const _TokenType_name = "TokenIllegalTokenEofTokenStringTokenNumberTokenTrueTokenFalseTokenNullTokenLbraceTokenLbrackTokenRbraceTokenRbrackTokenCommaTokenColon"

var _TokenType_index = [...]uint8{0, 12, 20, 31, 42, 51, 61, 70, 81, 92, 103, 114, 124, 134}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
