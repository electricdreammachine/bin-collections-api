package scrapercontroller

import (
	addresssearch "bin-collections-api/internal/controllers/address-search"
	getcollectiondates "bin-collections-api/internal/controllers/collections-for-address"
	"bin-collections-api/internal/services/config"
	scraperservice "bin-collections-api/internal/services/scraper"
	"encoding/json"
	"net/http"
)

func GetCollectionDates(w http.ResponseWriter, request *http.Request) {
	json.NewEncoder(w).Encode(
		scraperservice.PerformScrape(
			request.Body,
			config.ByKey("ADDITIONAL_COLLECTION_SEARCH_METADATA"),
			getcollectiondates.ForUniqueAddressID,
		),
	)
}

func GetAddresses(w http.ResponseWriter, request *http.Request) {
	json.NewEncoder(w).Encode(
		scraperservice.PerformScrape(
			request.Body,
			config.ByKey("ADDITIONAL_ADDRESS_SEARCH_METADATA"),
			addresssearch.ForPostCode,
		),
	)
}
