package models

type KafkaMsg struct {
	ContactEmail    []string `json:"to_list"`
	TemplateTitle   string   `json:"subject"`
	TemplateMessage string   `json:"body"`
	RetryCount      int      `json:"retry_count"`
}
