package models

type JSON map[string]interface{}

type AdditionalDataSchema struct {
	Values []MetaDataItem `json:"values"`
}
