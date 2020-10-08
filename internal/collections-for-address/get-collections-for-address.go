package getcollectiondates

import (
	"net/http"
	"log"
	"encoding/json"
	"bin-collections-api/internal/pkg/submit-flow-change"
	"bin-collections-api/internal/pkg/get-in-page-metadata"
)

type collectionSearchRequest struct {
	addressID string
}

// GetCollectionsForID blocks until tokens retrieved and then retrieves collections for the given ID
func GetCollectionsForID(w http.ResponseWriter, request *http.Request) {
	var collectionSearch collectionSearchRequest
	err := json.NewDecoder(request.Body).Decode(&collectionSearch)

	if err != nil {
		log.Fatal(err)
	}

	<-submitflowchange.Submit([])
	collectionDates := <-ForUniqueAddressID(collectionSearch.addressID)

	json.NewEncoder(w).Encode(collectionDates)
}