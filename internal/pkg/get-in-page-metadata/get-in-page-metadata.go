package getinpagemetadata

import (
	getconfigvalue "bin-collections-api/internal/pkg/get-config-value"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
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

// GetTokens will make a network request for required ephemeral tokens
func GetTokens() <-chan MetaData {
	c := colly.NewCollector()
	channel := make(chan MetaData)
	metaDataSchema := getconfigvalue.ByKey("METADATA_IN_PAGE")

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(err)
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		var objmap map[string]json.RawMessage
		err := json.Unmarshal([]byte(metaDataSchema), &objmap)
		var values []MetaDataItem

		if err != nil {
			log.Fatal(err)
		}

		err2 := json.Unmarshal(objmap["values"], &values)

		// fmt.Println(values)
		values = Populate(e, values)
		// fmt.Println(values)
		if err2 != nil {
			log.Fatal(err2)
		}

		headers := e.Response.Headers.Get(getconfigvalue.ByKey("TOKEN_HEADER"))
		appHeaderRegex := regexp.MustCompile(
			getconfigvalue.ByKey("TOKEN_REGEX"),
		)

		matchingSubGroups := appHeaderRegex.Find([]byte(headers))

		go func() {
			channel <- MetaData{
				Cookie:   strings.Split(string(matchingSubGroups), "="),
				MetaData: values,
			}
		}()
	})

	c.Visit(getconfigvalue.ByKey("TOKEN_URL"))

	return channel
}
