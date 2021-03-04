package scraperservice

import (
	"bin-collections-api/internal/models"
	config "bin-collections-api/internal/services/config"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

// GetTokens will make a network request for required ephemeral tokens
func GetTokens() <-chan models.MetaData {
	c := colly.NewCollector()
	channel := make(chan models.MetaData)
	metaDataSchema := config.ByKey("METADATA_IN_PAGE")

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(err)
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		var objmap map[string]json.RawMessage
		err := json.Unmarshal([]byte(metaDataSchema), &objmap)
		var values []models.MetaDataItem

		if err != nil {
			log.Fatal(err)
		}

		err2 := json.Unmarshal(objmap["values"], &values)
		values = Populate(e, values)

		if err2 != nil {
			log.Fatal(err2)
		}

		headers := e.Response.Headers.Get(config.ByKey("TOKEN_HEADER"))
		appHeaderRegex := regexp.MustCompile(
			config.ByKey("TOKEN_REGEX"),
		)

		matchingSubGroups := appHeaderRegex.Find([]byte(headers))

		go func() {
			channel <- models.MetaData{
				Cookie:   strings.Split(string(matchingSubGroups), "="),
				MetaData: values,
			}
		}()
	})

	c.Visit(config.ByKey("TOKEN_URL"))

	return channel
}
