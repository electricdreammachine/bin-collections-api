package getcollectiondates

import (
	"fmt"
	"net/http"
	"bin-collections-api/internal/pkg/get-tokens"
)

// GetCollectionsForID blocks until tokens retrieved and then retrieves collections for the given ID
func GetCollectionsForID(w http.ResponseWriter, request *http.Request) {
	//TODO: get id from request body
	fmt.Println("starting")
	instanceTokens := <-gettokens.GetTokens()
	collectionDates := <-ForUniqueAddressID(instanceTokens, "250056153")

	//TODO: construct http response
	fmt.Println(collectionDates)
}