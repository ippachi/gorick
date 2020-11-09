package http_parser

import (
	"fmt"
	"strings"
)

type ParsedHttpRequest struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

type ParseError struct{}

func (pe ParseError) Error() string {
	return fmt.Sprint("invalid format")
}

func Parse(request string) (ParsedHttpRequest, error) {
	eachLines := strings.Split(request, "\r\n")

	if len(eachLines) < 1 {
		return ParsedHttpRequest{}, ParseError{}
	}

	requestLineTokens := strings.Split(eachLines[0], " ")

	if len(requestLineTokens) != 3 {
		return ParsedHttpRequest{}, ParseError{}
	}

	method, target, version := requestLineTokens[0], requestLineTokens[1], requestLineTokens[2]

	if !isValidHttpMethod(method) {
		return ParsedHttpRequest{}, ParseError{}
	}

	_, err := originForm(target)

	if err != nil {
		return ParsedHttpRequest{}, ParseError{}
	}

	prh := ParsedHttpRequest{
		method,
		target,
		version,
	}

	return prh, nil
}

func isValidHttpMethod(method string) bool {
	HttpMethods := []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE"}

	for _, httpMethod := range HttpMethods {
		if httpMethod == method {
			return true
		}
	}

	return false
}

func originForm(target string) (int, error) {
	absolutePathLen := absolutePath(target, 0)
	if absolutePathLen == 0 {
		return 0, ParseError{}
	}
	return absolutePathLen, nil
}

func absolutePath(target string, i int) int {
	ret := 0
	for i < len(target) {
		if target[i] != '/' {
			return ret
		}

		i += 1
		ret += 1

		s := segment(target, i)

		i += s
		ret += s
	}

	return ret
}

func segment(target string, i int) int {
	ret := 0
	for n := pchar(target, i); n != 0; n = pchar(target, i) {
		ret += n
	}
	return ret
}

func pchar(target string, i int) int {
	c := target[i]

	if n := unreserved(target, i); n != 0 {
		return n
	}

	if n := pctEncoded(target, i); n != 0 {
		return n
	}

	if n := subDelims(target, i); n != 0 {
		return n
	}

	if c == ':' || c == '@' {
		return 1
	}

	return 0
}

func unreserved(target string, i int) int {
	c := target[i]
	if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' ||
		c == '-' || c == '.' || c == '_' || c == '~' {
		return 1
	}
	return 0
}

func pctEncoded(target string, i int) int {
	if len(target) < 3 && target[i] == '%' {
		return 0
	}

	for i := 1; i < 3; i++ {
		c := target[i]
		if 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' || '0' <= c && c <= '9' {
			continue
		}

		return 0
	}

	return 3
}

func subDelims(target string, i int) int {
	c := target[i]
	if c == '!' || c == '$' || c == '&' || c == '\'' || c == '(' ||
		c == ')' || c == '*' || c == '+' || c == ',' || c == ';' || c == '=' {
		return 1
	}
	return 0
}
