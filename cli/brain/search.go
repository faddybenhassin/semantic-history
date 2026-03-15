package brain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type SearchReq struct {
	Query string `json:"query"`
}

type SearchResp struct {
	Command string  `json:"command"`
	Score   float64 `json:"score"`
}

func FetchCommands(apiUrl, query string, limit int) ([]SearchResp, error) {

	params := url.Values{}
	params.Add("query", query)
	params.Add("limit", strconv.Itoa(limit)) // Convert int to string

	resp, err := http.Get(apiUrl + "/search" + "?" + params.Encode())
	if err != nil {
		return nil, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %v", resp.Status)
	}

	var results []SearchResp

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}

	return results, nil
}
