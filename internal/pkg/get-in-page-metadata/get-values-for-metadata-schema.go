package getinpagemetadata

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gocolly/colly"
)

func Populate(html *colly.HTMLElement, schema []MetaDataItem) []MetaDataItem {
	for i, v := range schema {
		var value string

		if v.Value != nil {
			value = fmt.Sprintf("%v", v.Value)
		} else {
			value, _ = html.DOM.Find(v.DomSelector).Attr("value")
		}

		v.Value = value

		if len(v.Format) > 0 {
			// use reflect to get field
			for formatKey, formatVal := range v.Format {
				if formatVal == "value" {
					v.Format[formatKey] = value
					continue
				}

				if len(formatVal) > 0 {
					rv := reflect.ValueOf(v)
					v.Format[formatKey] = reflect.Indirect(rv).FieldByName(strings.Title(strings.ToLower(strings.TrimSpace(formatVal)))).String()
					continue
				}
			}

			v.Value = v.Format
		}

		schema[i] = v
	}

	return schema
}
