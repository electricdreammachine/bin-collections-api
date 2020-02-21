package addresssearch

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"bin-collections-api/internal/pkg/get-config-value"
)

// JSON generic json type
type JSON map[string]interface{}

// FindAddressByPostCode decodes json request body for a postcode used to search for possible address entities
func FindAddressByPostCode(w http.ResponseWriter, r *http.Request) {
	type AddressSearch struct {
		PostCode string `json:"postCode"`
	}

	requestBody := json.NewDecoder(r.Body)

	var search AddressSearch
	err := requestBody.Decode(&search)

	spaceClient := http.Client{
		Timeout: time.Second * 15, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, getconfigvalue.ByKey("ADDRESS_SEARCH_URL"), nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Add(getconfigvalue.ByKey("ADDRESS_SEARCH_PARAM_1_KEY"), getconfigvalue.ByKey("ADDRESS_SEARCH_PARAM_1_VALUE"))
	q.Add(getconfigvalue.ByKey("ADDRESS_SEARCH_PARAM_2_KEY"), getconfigvalue.ByKey("ADDRESS_SEARCH_PARAM_2_VALUE"))
	q.Add(getconfigvalue.ByKey("ADDRESS_SEARCH_PARAM_3_KEY"), getconfigvalue.ByKey("ADDRESS_SEARCH_PARAM_3_VALUE"))
	q.Add(getconfigvalue.ByKey("ADDRESS_SEARCH_PARAM_4_KEY"), search.PostCode)
	q.Add("location", search.PostCode)
	q.Add(getconfigvalue.ByKey("ADDRESS_SEARCH_PARAM_5_KEY"), getconfigvalue.ByKey("ADDRESS_SEARCH_PARAM_5_VALUE"))
	q.Add("pageSize", "21")
	q.Add("startnum", "1")
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