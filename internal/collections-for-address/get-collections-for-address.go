package getcollectiondates

import (
	getconfigvalue "bin-collections-api/internal/pkg/get-config-value"
	getinpagemetadata "bin-collections-api/internal/pkg/get-in-page-metadata"
	submitflowchange "bin-collections-api/internal/pkg/submit-flow-change"
	"encoding/json"
	"log"
	"net/http"
)

type collectionSearchRequest map[string]interface{}

type additionalDataSchema struct {
	Values []getinpagemetadata.MetaDataItem `json:"values"`
}

// GetCollectionsForID blocks until tokens retrieved and then retrieves collections for the given ID
func GetCollectionsForID(w http.ResponseWriter, request *http.Request) {
	var parsedSchema additionalDataSchema
	metaDataSchema := getconfigvalue.ByKey("ADDITIONAL_COLLECTION_SEARCH_METADATA")
	err := json.Unmarshal([]byte(metaDataSchema), &parsedSchema)
	var collectionSearch collectionSearchRequest
	err2 := json.NewDecoder(request.Body).Decode(&collectionSearch)

	for i, v := range parsedSchema.Values {
		if len(v.ValueFromMap) > 0 {
			v.Value = collectionSearch[v.ValueFromMap]
		}
		parsedSchema.Values[i] = v
	}

	if err != nil {
		log.Fatal("err")
	}

	if err2 != nil {
		log.Fatal("err2")
	}

	cookie := <-submitflowchange.Submit(parsedSchema.Values)
	collectionDates := <-ForUniqueAddressID(cookie)

	json.NewEncoder(w).Encode(collectionDates)
}
