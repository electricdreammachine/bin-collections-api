package models

import (
	"encoding/json"
	"fmt"
)

type Cookie []string

type MetaDataItem struct {
	Name         string
	Path         string
	Value        interface{}
	DomSelector  string
	Format       map[string]string
	ValueFromMap string
}

type MetaData struct {
	Cookie   Cookie
	MetaData []MetaDataItem
}

type RedirectMetaData struct {
	Cookie      []string
	RedirectURL string
}

func (t *MetaDataItem) UnmarshalJSON(data []byte) error {
	type metaDataItemAlias MetaDataItem
	var iteratee metaDataItemAlias

	_ = json.Unmarshal(data, &iteratee)

	if len(iteratee.DomSelector) <= 0 {
		iteratee.DomSelector = fmt.Sprintf("[name=\"%v\"]", iteratee.Name)
	}

	*t = MetaDataItem(iteratee)

	return nil
}
