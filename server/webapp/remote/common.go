package remote

import "encoding/json"

type responseResult struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Data map[string]interface{}	`json:"data"`
}

func newResponseResult(p []byte) (*responseResult, error) {
	result := new(responseResult)
	err := json.Unmarshal(p, result)
	return result, err
}

