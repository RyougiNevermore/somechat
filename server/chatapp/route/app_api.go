package route

type result struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Data map[string]interface{}	`json:"data"`
}
