package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type myhttpSuite struct{}

var _ = Suite(&myhttpSuite{})

// Unit tests

func (s *myhttpSuite) Test_parseArgs_WithoutParams(c *C) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"/path/to/exec.ext"}

	expectedJobCount := 10
	expectedRawUrls := []string{}

	actualJobCount, actualRawUrls := parseArgs()

	c.Assert(actualJobCount, Equals, expectedJobCount)
	c.Assert(actualRawUrls, DeepEquals, expectedRawUrls)
}

func (s *myhttpSuite) Test_parseArgs_WithParallel(c *C) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"/path/to/exec.ext", "-parallel", "35"}

	expectedJobCount := 35
	expectedRawUrls := []string{}

	actualJobCount, actualRawUrls := parseArgs()

	c.Assert(actualJobCount, Equals, expectedJobCount)
	c.Assert(actualRawUrls, DeepEquals, expectedRawUrls)
}

func (s *myhttpSuite) Test_parseArgs_WithUrls(c *C) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"/path/to/exec.ext", "example.com", "http://example.org"}

	expectedJobCount := 10
	expectedRawUrls := []string{"example.com", "http://example.org"}

	actualJobCount, actualRawUrls := parseArgs()

	c.Assert(actualJobCount, Equals, expectedJobCount)
	c.Assert(actualRawUrls, DeepEquals, expectedRawUrls)
}

func (s *myhttpSuite) Test_parseArgs_WithUrlsAndParallel(c *C) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"/path/to/exec.ext", "-parallel", "35", "example.com", "http://example.org"}

	expectedJobCount := 35
	expectedRawUrls := []string{"example.com", "http://example.org"}

	actualJobCount, actualRawUrls := parseArgs()

	c.Assert(actualJobCount, Equals, expectedJobCount)
	c.Assert(actualRawUrls, DeepEquals, expectedRawUrls)
}

func (s *myhttpSuite) Test_generateMD5FromContent(c *C) {
	content := []byte("Hello, World!")

	expectedMD5 := "65a8e27d8879283831b664bd8b7f0ad4"

	actualMD5 := generateMD5FromContent(content)

	c.Assert(actualMD5, Equals, expectedMD5)
}

func (s *myhttpSuite) Test_parseRawURLs(c *C) {
	rawUrls := []string{"example.com", "http://example.org"}

	expectedUrls := []string{"http://example.com", "http://example.org"}

	actualUrls := parseRawURLs(rawUrls)

	c.Assert(actualUrls, DeepEquals, expectedUrls)
}

func (s *myhttpSuite) Test_validateURLs(c *C) {
	parsedUrls := []string{"http://example.com", "http://example"}

	expectedUrls := []string{"http://example.com"}

	actualUrls := validateURLs(parsedUrls)

	c.Assert(actualUrls, DeepEquals, expectedUrls)
}

func (s *myhttpSuite) Test_fetchURLContent(c *C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<html><body>Hello World!</body></html>")
	}))

	defer ts.Close()

	expectedUrl := ts.URL
	expectedContent := []byte("<html><body>Hello World!</body></html>\n")

	results := make(chan FetchedData)
	jobs := make(chan string)

	go fetchURLContent(jobs, results)

	jobs <- ts.URL
	close(jobs)

	actual := <-results
	close(results)

	c.Assert(actual.url, Equals, expectedUrl)
	c.Assert(actual.content, DeepEquals, expectedContent)
}

func (s *myhttpSuite) Test_fetchURLContent_InvalidProtocol(c *C) {
	expectedUrl := "xxx://malformed-url"
	expectedContent := []byte("Get xxx://malformed-url: unsupported protocol scheme \"xxx\"")

	results := make(chan FetchedData)
	jobs := make(chan string)

	go fetchURLContent(jobs, results)

	jobs <- "xxx://malformed-url"
	close(jobs)

	actual := <-results
	close(results)

	c.Assert(actual.url, Equals, expectedUrl)
	c.Assert(actual.content, DeepEquals, expectedContent)
	c.Assert(1, Equals, 1)
}

// Functional Tests

func (s *myhttpSuite) Test_start_WithProperParameters(c *C) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<html><body>Hello World!</body></html>")
	}))

	defer ts.Close()

	os.Args = []string{"/path/to/exec.ext", ts.URL}

	expectedReturn := 57

	actualResult := start()

	c.Assert(actualResult, Equals, expectedReturn)
}

func (s *myhttpSuite) Test_start_WithoutParameters(c *C) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<html><body>Hello World!</body></html>")
	}))

	defer ts.Close()

	os.Args = []string{"/path/to/exec.ext"}

	expectedReturn := 373

	actualResult := start()

	c.Assert(actualResult, Equals, expectedReturn)
}
