/*
Duerme, 2023
GoSpy, a directory enumeration tool.
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

	colorPurple := "\033[35m"
	colorReset := "\033[0m"	
	fmt.Println(string(colorPurple),`
	
 .d8888b.            .d8888b.                    
d88P  Y88b          d88P  Y88b                   
888    888          Y88b.                        
888         .d88b.   "Y888b.   88888b.  888  888 
888  88888 d88""88b     "Y88b. 888 "88b 888  888 
888    888 888  888       "888 888  888 888  888 
Y88b  d88P Y88..88P Y88b  d88P 888 d88P Y88b 888 
 "Y8888P88  "Y88P"   "Y8888P"  88888P"   "Y88888 
			       888           888 
			       888      Y8b d88P 
			       888       "Y88P"
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~			       		
			              By: Duerme`,string(colorReset))

	fmt.Print("\n\n")
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

	wordcount := 0
	var words[]string

	for scanner.Scan(){ // Calculate the amount of directory names
		wordcount++
		words = append(words,scanner.Text())
	}

	fmt.Printf("Loaded %d possible directory names. Enumeration will now begin.\n", wordcount)

	for i:= 0; i < len(words); i++ { // Loop through each line of the wordlist
		directoryName := words[i]
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
