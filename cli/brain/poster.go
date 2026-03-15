package brain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type IndexReq struct {
	Command string `json:"command"`
}

type IndexResp struct {
	Status string `json:"status"`
	ID     string `json:"id"`
}

func PostCommand(apiUrl, command string) error {
	reqBody := IndexReq{Command: command}
	b, _ := json.Marshal(reqBody)

	resp, err := http.Post(apiUrl+"/index", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %v", resp.Status)
	}

	var out IndexResp

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return err
	}

	// fmt.Println("indexed", out.ID)
	return nil
}
