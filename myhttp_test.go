package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

// Unit tests

func TestParseArgsWithoutParams(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"/path/to/exec.ext"}

	expectedJobCount := 10
	expectedRawUrls := []string{}

	actualJobCount, actualRawUrls := parseArgs()

 	if actualJobCount != expectedJobCount {
		t.Errorf("actualJobCount (%d) does not match expectedJobCount (%d)", actualJobCount, expectedJobCount)
	}

	if !reflect.DeepEqual(actualRawUrls, expectedRawUrls) {
		t.Errorf("actualRawUrls (%s) does not match expectedRawUrls (%s)", actualRawUrls, expectedRawUrls)
	}
}

func TestParseArgsWithParallel(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"/path/to/exec.ext", "-parallel", "35"}

	expectedJobCount := 35
	expectedRawUrls := []string{}

	actualJobCount, actualRawUrls := parseArgs()

 	if actualJobCount != expectedJobCount {
		t.Errorf("actualJobCount (%d) does not match expectedJobCount (%d)", actualJobCount, expectedJobCount)
	}

	if !reflect.DeepEqual(actualRawUrls, expectedRawUrls) {
		t.Errorf("actualRawUrls (%s) does not match expectedRawUrls (%s)", actualRawUrls, expectedRawUrls)
	}
}

func TestParseArgsWithUrls(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"/path/to/exec.ext", "example.com", "http://example.org"}

	expectedJobCount := 10
	expectedRawUrls := []string{"example.com", "http://example.org"}

	actualJobCount, actualRawUrls := parseArgs()

 	if actualJobCount != expectedJobCount {
		t.Errorf("actualJobCount (%d) does not match expectedJobCount (%d)", actualJobCount, expectedJobCount)
	}

	if !reflect.DeepEqual(actualRawUrls, expectedRawUrls) {
		t.Errorf("actualRawUrls (%s) does not match expectedRawUrls (%s)", actualRawUrls, expectedRawUrls)
	}
}

func TestParseArgsWithUrlsAndParallel(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"/path/to/exec.ext", "-parallel", "35", "example.com", "http://example.org"}

	expectedJobCount := 35
	expectedRawUrls := []string{"example.com", "http://example.org"}

	actualJobCount, actualRawUrls := parseArgs()

 	if actualJobCount != expectedJobCount {
		t.Errorf("actualJobCount (%d) does not match expectedJobCount (%d)", actualJobCount, expectedJobCount)
	}

	if !reflect.DeepEqual(actualRawUrls, expectedRawUrls) {
		t.Errorf("actualRawUrls (%s) does not match expectedRawUrls (%s)", actualRawUrls, expectedRawUrls)
	}
}

func TestGenerateMD5FromContent(t *testing.T) {
	content := []byte("Hello, World!")

	expectedMD5 := "65a8e27d8879283831b664bd8b7f0ad4"

	actualMD5 := generateMD5FromContent(content)

 	if actualMD5 != expectedMD5 {
		t.Errorf("actualMD5 (%s) does not match expectedMD5 (%s)", actualMD5, expectedMD5)
	}
}

func TestParseRawURLs(t *testing.T) {
	rawUrls := []string{"example.com", "http://example.org"}

	expectedUrls := []string{"http://example.com", "http://example.org"}

	actualUrls := parseRawURLs(rawUrls)

 	if !reflect.DeepEqual(actualUrls, expectedUrls) {
		t.Errorf("actualUrls (%s) does not match expectedUrls (%s)", actualUrls, expectedUrls)
	}
}

func TestValidateURLs(t *testing.T) {
	parsedUrls := []string{"http://example.com", "http://example"}

	expectedUrls := []string{"http://example.com"}

	actualUrls := validateURLs(parsedUrls)

 	if !reflect.DeepEqual(actualUrls, expectedUrls) {
		t.Errorf("actualUrls (%s) does not match expectedUrls (%s)", actualUrls, expectedUrls)
	}
}

func TestFetchURLContentValidURL(t *testing.T) {
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

 	if actual.url != expectedUrl {
		t.Errorf("actual.url (%s) does not match expectedUrl (%s)", actual.content, expectedUrl)
	}

 	if !reflect.DeepEqual(actual.content, expectedContent) {
		t.Errorf("actual.content (%s) does not match expectedContent (%s)", actual.content, expectedContent)
	}
}

func TestFetchURLContentInvalidProtocol(t *testing.T) {
	expectedUrl := "xxx://malformed-url"
	expectedContent := []byte("Get xxx://malformed-url: unsupported protocol scheme \"xxx\"")

	results := make(chan FetchedData)
	jobs := make(chan string)

	go fetchURLContent(jobs, results)

	jobs <- "xxx://malformed-url"
	close(jobs)

	actual := <-results
	close(results)

 	if actual.url != expectedUrl {
		t.Errorf("actual.url (%s) does not match expectedUrl (%s)", actual.content, expectedUrl)
	}

 	if !reflect.DeepEqual(actual.content, expectedContent) {
		t.Errorf("actual.content (%s) does not match expectedContent (%s)", actual.content, expectedContent)
	}
}

// Functional Tests

func TestStartWithProperParameters(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<html><body>Hello World!</body></html>")
	}))

	defer ts.Close()

	os.Args = []string{"/path/to/exec.ext", ts.URL}

	expectedReturn := 57

	actualResult := start()

 	if actualResult != expectedReturn {
		t.Errorf("actualResult (%d) does not match expectedReturn (%d)", actualResult, expectedReturn)
	}
}

func TestStartWithoutParameters(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<html><body>Hello World!</body></html>")
	}))

	defer ts.Close()

	os.Args = []string{"/path/to/exec.ext"}

	expectedReturn := 373

	actualResult := start()

 	if actualResult != expectedReturn {
		t.Errorf("actualResult (%d) does not match expectedReturn (%d)", actualResult, expectedReturn)
	}
}
