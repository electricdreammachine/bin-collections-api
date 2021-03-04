package models

import (
	"encoding/json"
	"fmt"
)

// Cookie contains cookie headers
type Cookie []string

// MetaDataItem contains all the metadata needing to be added to requests
type MetaDataItem struct {
	Name         string
	Path         string
	Value        interface{}
	DomSelector  string
	Format       map[string]string
	ValueFromMap string
}

// MetaData f
type MetaData struct {
	Cookie   Cookie
	MetaData []MetaDataItem
}

// RedirectMetaData f
type RedirectMetaData struct {
	Cookie      []string
	RedirectUrl string
}

// UnmarshalJSON Custom unmarshaler
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
