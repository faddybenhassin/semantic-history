package main

import (
	"fmt"
	"os"
	"semantic-history/cli/brain"
	"semantic-history/cli/history"
)

func main() {
	apiUrl := "http://0.0.0.0:8000"
	threshold := 0.5
	path, err := history.GetPath()
	if err != nil {
		fmt.Println("Could not find path: ", err)
	}
	lines, err := history.ReadLines(path)
	if err != nil {
		fmt.Println("Could not read lines: ", err)
	}

	for _, line := range lines {
		if err := brain.PostCommand(apiUrl, line); err != nil {
			fmt.Println("post error", err)
		}
	}

	if len(os.Args) != 2 {
		fmt.Println("Usage: semanticSearch <query>")
		return
	}

	// Grab the string at index 1
	query := os.Args[1]

	results, err := brain.FetchCommands(apiUrl, query, 3)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Print the results nicely
	var res *brain.SearchResp
	for _, r := range results {
		if res == nil || res.Score < r.Score {
			res = &r
		}
		// fmt.Printf("[%0.2f] %s\n", r.Score, r.Command)
	}
	if res != nil && res.Score > threshold {
		fmt.Println(res.Command)
	}
}
