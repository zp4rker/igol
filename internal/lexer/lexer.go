package lexer

import (
	"errors"
	"regexp"
	"strings"
)

var (
	number_regex  = compile_regex(`[0-9]+`)
	symbol_regex  = compile_regex(`[a-zA-Z_]+`)
	keyword_regex = compile_regex(`(?:(set)|(to)|(when)|(otherwise))`)
)

type Token struct {
	Type  string
	Value string
}

func Lex(str string) ([]Token, error) {
	var tokens []Token
	var token Token

	for i := 0; i < len(str); i++ {
		if strings.Contains(" \t", string(str[i])) {
			continue
		} else if strings.Contains(`"`, string(str[i])) {
			i++
			val, err := scan_delim(`"`, str[i:])
			if err != nil {
				return tokens, err
			}
			token = Token{Type: "string", Value: val}
			i += len(val)
		} else if number_regex.MatchString(string(str[i])) {
			val, err := scan_regex(number_regex, str[i:])
			if err != nil {
				return tokens, err
			}
			token = Token{Type: "number", Value: val}
			i += len(val) - 1
		} else if symbol_regex.MatchString(string(str[i])) {
			val, err := scan_regex(symbol_regex, str[i:])
			if err != nil {
				return tokens, err
			}
			if keyword_regex.MatchString(val) {
				token = Token{Type: "keyword", Value: val}
			} else {
				token = Token{Type: "symbol", Value: val}
			}
			i += len(val) - 1
		} else if strings.Contains("()", string(str[i])) {
			token = Token{Type: string(str[i]), Value: ""}
		} else {
			token = Token{Type: "unknown", Value: string(str[i])}
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

func scan_delim(delim string, str string) (string, error) {
	result := ""
	for _, char := range str {
		if string(char) == delim {
			return result, nil
		} else {
			result += string(char)
		}
	}
	return result, errors.New("unexpected end of string")
}

func scan_regex(re *regexp.Regexp, str string) (string, error) {
	result := ""
	for _, char := range str {
		if !re.MatchString(string(char)) {
			return result, nil
		} else {
			result += string(char)
		}
	}
	return result, nil
}

func compile_regex(str string) *regexp.Regexp {
	re, err := regexp.Compile(str)
	if err != nil {
		panic(err)
	}
	return re
}
