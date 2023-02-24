/*
Duerme, 2023
GoSpy, the multi purpose web scanning tool.
*/
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
	reader := bufio.NewReader(os.Stdin) // Initialize user input reader
	var websiteAddress string
	for valInput != true { // re-run prompt until user inputs a valid URL  

		fmt.Print("Enter website address (e.g. https://example.com): ")
		input, _ := reader.ReadString('\n') // Read user input, then start a new line
		websiteAddress = strings.TrimSpace(input) // cut any white out of the input

		parsedURL, err := url.Parse(websiteAddress)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" { // validate the URL
			fmt.Printf("Invalid website address: %s\n", websiteAddress)
	} else {
		valInput = true
	}

	}

	wordlistURL := "https://raw.githubusercontent.com/danielmiessler/SecLists/master/Discovery/Web-Content/raft-large-directories.txt" // Wordlist used to enumerate directories
	wordlistResponse, err := http.Get(wordlistURL) // If there is a problem retrieving the wordlist, end program
	if err != nil {
		fmt.Printf("Error opening wordlist file: %s\n", err)
		return
	}
	defer wordlistResponse.Body.Close() // Keep reading wordlist response until the scan is done

	scanner := bufio.NewScanner(wordlistResponse.Body) // Open the wordlist body with new Scanner
	for scanner.Scan() { // Loop through each line of the wordlist
		directoryName := scanner.Text()
		directoryURL := websiteAddress + "/" + directoryName
		response, err := http.Head(directoryURL)
		if err == nil && response.StatusCode == 200 { // Ignore any status return other than 200 and print out the directory
			fmt.Printf("%s - %d\n", directoryURL, response.StatusCode)
		} else if err != nil{
			fmt.Printf("%s - Error: %s\n", directoryURL, err) // If there is an error print it out to the console
		}
	}

	if err := scanner.Err(); err != nil { 
		fmt.Printf("Error reading wordlist file: %s\n", err)
		return
	}
}
