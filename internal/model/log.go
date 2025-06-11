package model

// Log はログデータを表す構造体
type Log struct {
	ID        string            `json:"id"`
	TraceID   string            `json:"traceId"`
	Timestamp string            `json:"timestamp"`
	Level     string            `json:"level"`
	Service   string            `json:"service"`
	Message   string            `json:"message"`
	Metadata  map[string]string `json:"metadata"`
}
