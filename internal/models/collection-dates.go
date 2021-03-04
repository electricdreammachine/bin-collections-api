package models

// Collections describes the collections available in the date range
type Collections struct {
	Dates []Collection      `json:"dates"`
	Types map[string]string `json:"types"`
}

type Collection struct {
	Type []string `json:"type"`
	Date string   `json:"date"`
}
