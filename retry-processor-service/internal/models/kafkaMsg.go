package models

type Message struct {
	Subject    string   `json:"subject"`
	Body       string   `json:"body"`
	ToList     []string `json:"to_list"`
	RetryCount int      `json:"retry_count"`
}
