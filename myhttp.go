// A program that computes the MD5 of Websites from the URL.
// Done as a part of take-home task for Adjust.
package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// A simple struct to store the results of URL fetch.
type FetchedData struct {
	url     string
	content []byte
}

func main() {
	start()
}

// Main function driving the program.
func start() int {
	wg := sync.WaitGroup{}
	totalBytesWritten := 0
	results := make(chan FetchedData)
	jobs := make(chan string)

	if len(os.Args) == 1 {
		return displayUsage()
	}

	jobCount, rawUrls := parseArgs()

	parsedUrls := parseRawURLs(rawUrls)
	validatedUrls := validateURLs(parsedUrls)

	// Set up the parallel go routines to fetch URL content.
	for i := 0; i < jobCount; i++ {
		go fetchURLContent(jobs, results)
	}

	wg.Add(len(validatedUrls))
	// Add the URLs to the jobs queue.
	go func() {
		for _, validatedUrl := range validatedUrls {
			jobs <- validatedUrl
		}
		close(jobs)
	}()

	// Blocking wait until all the routines are done.
	go func() {
		wg.Wait()
		close(results)
	}()

	// Generate the MD5 and print the results.
	for result := range results {
		md5sum := generateMD5FromContent(result.content)
		bytesWritten, _ := fmt.Printf("%s: %s\n", result.url, md5sum)
		totalBytesWritten += bytesWritten
		wg.Done()
	}

	return totalBytesWritten
}

// Get the arguments from CLI.
func parseArgs() (int, []string) {
	jobsPtr := flag.Int("parallel", 10, "Number of parallel threads to run")

	flag.Parse()

	return *jobsPtr, flag.Args()
}

// Fetches the content of the URLs from the urls channel and returns the results
// in the Feteched Data channel.
func fetchURLContent(urls chan string, results chan FetchedData) {
	for url := range urls {
		defer func() {
			if r := recover(); r != nil {
				err := r.(error)
				results <- FetchedData{url, []byte(err.Error())}
			}
		}()

		resp, err := http.Get(url)

		if err != nil {
			panic(err)
		}

		// reads html as a slice of bytes
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		if resp != nil {
			resp.Body.Close()
		}

		results <- FetchedData{url, content}
	}
}

// Generates the MD5 of a byte array and returns the generated MD5 string.
func generateMD5FromContent(content []byte) string {
	md5 := md5.Sum(content)

	return hex.EncodeToString(md5[:])
}

// Parses the Raw URLs from CLI and returns a list of Parsed URLs.
func parseRawURLs(rawUrls []string) []string {
	parsedUrls := []string{}

	for _, rawUrl := range rawUrls {
		if !(strings.HasPrefix(rawUrl, "http://") || strings.HasPrefix(rawUrl, "https://")) {
			rawUrl = "http://" + rawUrl
		}
		parsedUrls = append(parsedUrls, rawUrl)
	}

	return parsedUrls
}

// Validates the Parsed URLs from the parseRawURLs and returns the list of Validated URLs.
func validateURLs(parsedUrls []string) []string {
	validatedUrls := []string{}

	for _, parsedUrl := range parsedUrls {
		_, err := url.ParseRequestURI(parsedUrl)
		if err != nil || !strings.Contains(parsedUrl, ".") {
			log.Printf("Skipping: [ %s ], unable to parse\n", parsedUrl)
			continue
		}
		validatedUrls = append(validatedUrls, parsedUrl)
	}

	return validatedUrls
}

// Displays a friendly help message.
func displayUsage() int {
	executableName := strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(filepath.Base(os.Args[0])))
	helpText := "Usage: " + executableName + " [-parallel JOBS] URL1 URL2...\n\n" +
		"Checks if a give URL is valid and generates the MD5 of the contents,\n" +
		"Example: " + executableName + " -parallel 1 http://adjust.com\n\n" +
		"Interpretation of parameters:\n" +
		"\tJOBS\t\tNumber simultaneous jobs to process the URLs. Defaults to 10.\n" +

		"\tURL\t\tA list of space separated URLs to fetch content and generate MD5.\n" +
		"\t\t\tFor example adjust.com http://google.com\n"

	bytesWritten, _ := fmt.Println(helpText)

	return bytesWritten
}
