package models

type Collections struct {
	Dates []Collection      `json:"dates"`
	Types map[string]string `json:"types"`
}

type Collection struct {
	Type []string `json:"type"`
	Date string   `json:"date"`
}
