package addresssearch

import (
	getconfigvalue "bin-collections-api/internal/pkg/get-config-value"
	getinpagemetadata "bin-collections-api/internal/pkg/get-in-page-metadata"
	submitflowchange "bin-collections-api/internal/pkg/submit-flow-change"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// JSON generic json type
type JSON map[string]interface{}

type AddressSearch map[string]interface{}

type additionalDataSchema struct {
	Values []getinpagemetadata.MetaDataItem
}

// FindAddressByPostCode decodes json request body for a postcode used to search for possible address entities
func FindAddressByPostCode(w http.ResponseWriter, r *http.Request) {
	var parsedSchema additionalDataSchema
	metaDataSchema := getconfigvalue.ByKey("ADDITIONAL_ADDRESS_SEARCH_METADATA")
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
	addresses := ForPostCode(submitResponse.Cookie)

	fmt.Println(addresses)
}
