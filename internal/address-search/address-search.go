package addresssearch

import (
	getconfigvalue "bin-collections-api/internal/pkg/get-config-value"
	getinpagemetadata "bin-collections-api/internal/pkg/get-in-page-metadata"
	submitflowchange "bin-collections-api/internal/pkg/submit-flow-change"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// JSON generic json type
type JSON map[string]interface{}

type AddressSearch map[string]interface{}

type additionalDataSchema struct {
	values []getinpagemetadata.MetaDataItem
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
	err2 := requestBody.Decode(&search)

	spaceClient := http.Client{
		Timeout: time.Second * 15,
	}

	for _, v := range parsedSchema.values {
		if len(v.ValueFromMap) > 0 {
			v.Value = search[v.ValueFromMap]
		}
	}

	submitflowchange.Submit(parsedSchema.values)

	req, _ := http.NewRequest(http.MethodGet, getconfigvalue.ByKey("ADDRESS_SEARCH_URL"), nil)
	if err2 != nil {
		log.Fatal(err2)
	}

	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var pdata JSON

	jsonErr := json.Unmarshal(body, &pdata)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	json.NewEncoder(w).Encode(pdata)
}
