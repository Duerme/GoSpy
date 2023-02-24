package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {

	valInput := false
	reader := bufio.NewReader(os.Stdin) // initialize the reader to read user input
	var websiteAddress string
	for !valInput {

		fmt.Print("Enter website address (e.g. https://example.com): ")
		input, _ := reader.ReadString('\n')
		websiteAddress = strings.TrimSpace(input)

		parsedURL, err := url.Parse(websiteAddress)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
			fmt.Printf("Invalid website address: %s\n", websiteAddress)
		} else {
			valInput = true
		}
	}

	// fmt.Print("website input: " + websiteAddress)

	wordlistURL := "https://raw.githubusercontent.com/danielmiessler/SecLists/master/Discovery/Web-Content/raft-large-directories.txt"
	wordlistResponse, err := http.Get(wordlistURL)
	if err != nil {
		fmt.Printf("Error opening wordlist file: %s\n", err)
		return
	}
	defer wordlistResponse.Body.Close()

	scanner := bufio.NewScanner(wordlistResponse.Body)
	for scanner.Scan() {
		directoryName := scanner.Text()
		directoryURL := websiteAddress + "/" + directoryName
		response, err := http.Head(directoryURL)
		if err == nil && response.StatusCode == 200 {
			fmt.Printf("%s - %d\n", directoryURL, response.StatusCode)
		} else if err != nil {
			fmt.Printf("%s - Error: %s\n", directoryURL, err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading wordlist file: %s\n", err)
		return
	}
}
