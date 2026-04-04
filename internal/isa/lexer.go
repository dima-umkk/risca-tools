package isa

import (
	"fmt"
	"strconv"
	"strings"
)

func Tokenize(origInput string) ([]Token, error) { // input line in uppercase, e.g. "ADD R1, R2"
	input := strings.ToUpper(origInput)
	var tokens []Token
	pos := 0
	for pos < len(input) {
		if isWhitespace(input, pos) {
			pos = skipWhitespace(input, pos)
		} else if isWord(input, pos) {
			word, newPos := readWord(input, pos)
			pos = newPos
			if tokenType, exists := mapWordToToken[word]; exists {
				tokens = append(tokens, Token{T: tokenType, Tk: word})
			} else {
				tokens = append(tokens, Token{T: TK_LABEL, Tk: word})
				//return nil, fmt.Errorf("unknown token: %s", word)
			}
		} else if isNumber(input, pos) {
			number, newPos, isParsed := readNumber(input, pos)
			if !isParsed {
				return nil, fmt.Errorf("Not a number on pos %d: %s>%s", pos, input[:pos], input[pos:])
			}
			tokens = append(tokens, Token{T: TK_NUMBER, Tk: input[pos:newPos], ValInt: number})
			pos = newPos
		} else {
			var token Token
			var found bool
			token, pos, found = readString(origInput, pos)
			if found {
				tokens = append(tokens, token)
				continue
			}
			token, pos, found = getBucks(input, pos)
			if found {
				tokens = append(tokens, token)
				continue
			}
			token, pos, found = getAt(input, pos)
			if found {
				tokens = append(tokens, token)
				continue
			}
			token, pos, found = getCommaOrColon(input, pos)
			if found {
				tokens = append(tokens, token)
				continue
			}
			token, pos, found = getPlusMinus(input, pos)
			if found {
				tokens = append(tokens, token)
				continue
			}
			token, pos, found = getBracket(input, pos)
			if found {
				tokens = append(tokens, token)
				continue
			}
			token, pos, found = getCompareOperator(input, pos)
			if found {
				tokens = append(tokens, token)
				continue
			}

			return nil, fmt.Errorf("Unexpeced symbol on pos %d: '%c'. %s>%s", pos, input[pos], input[:pos], input[pos:])
		}
	}
	return tokens, nil
}

func getCompareOperator(input string, pos int) (Token, int, bool) {
	switch input[pos] {
	case '!':
		pos++
		if pos < len(input) && input[pos] == '=' {
			pos++
			return Token{T: TK_CMP_NEQ, Tk: "!="}, pos, true
		}
	case '=':
		pos++
		if pos < len(input) && input[pos] == '=' {
			pos++
			return Token{T: TK_CMP_EQ, Tk: "=="}, pos, true
		}
	case '<':
		pos++
		if pos < len(input) && input[pos] == '=' {
			pos++
			return Token{T: TK_CMP_LTEQ, Tk: "<="}, pos, true
		} else {
			return Token{T: TK_CMP_LT, Tk: "<"}, pos, true
		}
	case '>':
		pos++
		if pos < len(input) && input[pos] == '=' {
			pos++
			return Token{T: TK_CMP_GTEQ, Tk: ">="}, pos, true
		} else {
			return Token{T: TK_CMP_GT, Tk: ">"}, pos, true
		}
	}
	return Token{}, pos, false
}

func getBucks(input string, pos int) (Token, int, bool) {
	if input[pos] == '$' {
		pos++
		return Token{T: TK_BUCKS, Tk: "$"}, pos, true
	}
	return Token{}, pos, false
}

func getAt(input string, pos int) (Token, int, bool) {
	if input[pos] == '@' {
		pos++
		return Token{T: TK_AT, Tk: "@"}, pos, true
	}
	return Token{}, pos, false
}

func getPlusMinus(input string, pos int) (Token, int, bool) {
	switch input[pos] {
	case '+':
		pos++
		return Token{T: TK_PLUS, Tk: "+"}, pos, true
	case '-':
		pos++
		return Token{T: TK_MINUS, Tk: "-"}, pos, true
	default:
		return Token{}, pos, false
	}

}

func getBracket(input string, pos int) (Token, int, bool) {
	switch input[pos] {
	case '[':
		pos++
		return Token{T: TK_L_SQBR, Tk: "["}, pos, true
	case ']':
		pos++
		return Token{T: TK_R_SQBR, Tk: "]"}, pos, true
	default:
		return Token{}, pos, false
	}
}

func getCommaOrColon(input string, pos int) (Token, int, bool) {
	switch input[pos] {
	case ':':
		pos++
		return Token{T: TK_COLON, Tk: ":"}, pos, true
	case ',':
		pos++
		return Token{T: TK_COMMA, Tk: ","}, pos, true
	default:
		return Token{}, pos, false
	}
}

func isWord(input string, pos int) bool {
	return (input[pos] >= 'A' && input[pos] <= 'Z') || (input[pos] == '.')
}

func isNumber(input string, pos int) bool {
	if isDigit(input, pos) {
		return true
	}
	if input[pos] == '0' && len(input) > pos+1 && input[pos+1] == 'X' {
		return true
	}
	return false
}

func isDigit(input string, pos int) bool {
	return input[pos] >= '0' && input[pos] <= '9'
}

func isWhitespace(input string, pos int) bool {
	return input[pos] == ' ' || input[pos] == '\t' || input[pos] == ';'
}

func skipWhitespace(input string, pos int) int {
	for pos < len(input) && (input[pos] == ' ' || input[pos] == '\t' || input[pos] == ';') {
		if input[pos] == ';' { // Skip line comment
			return len(input)
		} else {
			pos++
		}
	}
	return pos
}

func readWord(input string, pos int) (string, int) { // Read word started from letter, can contain letters, digits and dots (R1, R2, etc.)
	start := pos
	for pos < len(input) && (isWord(input, pos) || isDigit(input, pos) || input[pos] == '_') {
		pos++
	}
	return input[start:pos], pos
}

func isHexDigit(input string, pos int) bool {
	return (input[pos] >= '0' && input[pos] <= '9') || (input[pos] >= 'A' && input[pos] <= 'F')
}

func readNumber(input string, pos int) (int32, int, bool) { // Read number started from digit
	if input[pos] == '0' && len(input) > pos+1 && input[pos+1] == 'X' { //Read hex
		pos += 2 // Skip '0X' prefix
		start := pos
		for pos < len(input) && isHexDigit(input, pos) {
			pos++
		}
		value, _ := strconv.ParseInt(input[start:pos], 16, 64)
		return int32(value), pos, true
	} else { //Read decimal
		start := pos
		for pos < len(input) && isDigit(input, pos) {
			pos++
		}
		value, err := strconv.Atoi(input[start:pos])
		if err != nil {
			return 0, pos, false
		}
		return int32(value), pos, true
	}

}

func readString(input string, pos int) (Token, int, bool) {
	if input[pos] != '\'' {
		return Token{}, pos, false
	}
	pos++
	start := pos
	for pos < len(input) && input[pos] != '\'' {
		pos++
	}
	if input[pos] != '\'' {
		return Token{}, pos, false
	}
	str := string(input[start:pos])
	return Token{T: TK_STRING, Tk: str, ValStr: str}, pos + 1, true
}
