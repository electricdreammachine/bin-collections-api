package addresssearch

import (
	"bin-collections-api/internal/models"
	config "bin-collections-api/internal/services/config"
	submitflowchange "bin-collections-api/internal/services/scraper"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type AddressSearch map[string]interface{}

type additionalDataSchema struct {
	Values []models.MetaDataItem
}

// FindAddressByPostCode decodes json request body for a postcode used to search for possible address entities
func FindAddressByPostCode(w http.ResponseWriter, r *http.Request) {
	var parsedSchema additionalDataSchema
	metaDataSchema := config.ByKey("ADDITIONAL_ADDRESS_SEARCH_METADATA")
	err := json.Unmarshal([]byte(metaDataSchema), &parsedSchema)

	if err != nil {
		log.Fatal(err)
	}

	requestBody := json.NewDecoder(r.Body)

	var search AddressSearch
	requestDecodeError := requestBody.Decode(&search)

	if requestDecodeError != nil {
		fmt.Println(requestDecodeError)
	}

	for i, v := range parsedSchema.Values {
		if len(v.ValueFromMap) > 0 {
			v.Value = search[v.ValueFromMap]
		}
		parsedSchema.Values[i] = v
	}

	submitResponse := <-submitflowchange.Submit(parsedSchema.Values)
	addresses := ForPostCode(submitResponse.RedirectUrl, submitResponse.Cookie)

	fmt.Println(addresses)
}
