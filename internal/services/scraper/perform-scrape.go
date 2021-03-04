package scraperservice

import (
	"bin-collections-api/internal/models"
	"bin-collections-api/internal/services/config"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gocolly/colly"
)

func PerformScrape(requestBody io.Reader, additionalJSONSchemaForScrape string, scrapeCallback func(*colly.Collector) chan interface{}) interface{} {
	var parsedAdditionalSchema models.AdditionalDataSchema
	metaDataSchema := additionalJSONSchemaForScrape
	err := json.Unmarshal([]byte(metaDataSchema), &parsedAdditionalSchema)
	var requestModel models.JSON
	err2 := json.NewDecoder(requestBody).Decode(&requestModel)

	for i, v := range parsedAdditionalSchema.Values {
		if len(v.ValueFromMap) > 0 {
			v.Value = requestModel[v.ValueFromMap]
		}
		parsedAdditionalSchema.Values[i] = v
	}

	if err != nil {
		log.Fatal("Error getting required additional schema")
	}

	if err2 != nil {
		log.Fatal("Error getting required input from request")
	}

	submitFlowChangeResponse := <-Submit(parsedAdditionalSchema.Values)

	scrapeCollector := colly.NewCollector()

	scrapeCollector.SetCookies(
		config.ByKey("DATES_COOKIE_DOMAIN"),
		[]*http.Cookie{
			{
				Name:  submitFlowChangeResponse.Cookie[0],
				Value: submitFlowChangeResponse.Cookie[1],
			},
		},
	)

	dataChannel := scrapeCallback(scrapeCollector)

	scrapeCollector.Visit(fmt.Sprintf("%v/portal/%v", config.ByKey("DATES_COOKIE_DOMAIN"), submitFlowChangeResponse.RedirectUrl))

	dataFromPage := <-dataChannel

	return dataFromPage
}
