package http_parser

import (
	"testing"
)

func TestParseSuccessful(t *testing.T) {
	request := "GET / HTTP/1.1\r\n"

	prh, err := Parse(request)

	if err != nil {
		t.FailNow()
	}

	if prh.Method != "GET" || prh.RequestTarget != "/" || prh.HttpVersion != "HTTP/1.1" {
		t.FailNow()
	}
}

func TestParseFailureWithInvalidRequestLineItemCount(t *testing.T) {
	missingMethod := "/ HTTP/1.1\r\n"
	missingRequestTarget := "GET HTTP/1.1\r\n"
	missingHttpVersion := "GET /\r\n"

	for _, invalid := range []string{missingMethod, missingRequestTarget, missingHttpVersion} {
		_, err := Parse(invalid)

		if err == nil {
			t.Errorf("%s is invalid", invalid)
		}
	}
}

func TestParseHttpMethod(t *testing.T) {
	validMethods := []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE"}

	for _, validMethod := range validMethods {
		prh, err := Parse(validMethod + " / HTTP/1.1\r\n")

		if err != nil {
			t.Errorf("%s is valid method", validMethod)
		}

		if prh.Method != validMethod {
			t.Errorf("%s should eq %s", prh.Method, validMethod)
		}
	}

	invalidMethods := []string{"FETCH", "CATCH", "TOUCH"}

	for _, invalidMethod := range invalidMethods {
		_, err := Parse(invalidMethod + " / HTTP/1.1\r\n")

		if err == nil {
			t.Errorf("%s is invalid method", invalidMethod)
		}
	}
}

func TestParseRequestTarget(t *testing.T) {
	valids := []string{
		"/",
		"/test",
		"/test/test",
		"/test/%00%99%AA%ZZ%aa%zz",
		"/test/!$&'()*+,;=",
		"/test/:/@",
	}

	for _, valid := range valids {
		request := "GET " + valid + " HTTP/1.1"

		prh, err := Parse(request)

		if err != nil {
			t.Errorf("%s is valid", valid)
		}

		if prh.RequestTarget != valid {
			t.Errorf("%s should eq %s", prh.RequestTarget, valid)
		}
	}

	invalids := []string{
		"test/test",
		"/test/|",
		"/test%",
		"/test%A",
		"/test%AX",
	}

	for _, invalid := range invalids {
		request := "GET " + invalid + " HTTP/1.1"

		_, err := Parse(request)

		if err == nil {
			t.Errorf("%s is invalid", invalid)
		}
	}
}
