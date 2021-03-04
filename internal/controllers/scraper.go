package scrapercontroller

import (
	"bin-collections-api/internal/services/config"
	scraperservice "bin-collections-api/internal/services/scraper"
	"encoding/json"
	"net/http"
)

func GetCollectionDates(w http.ResponseWriter, request *http.Request) {
	json.NewEncoder(w).Encode(scraperservice.PerformScrape(request.Body, config.ByKey("ADDITIONAL_COLLECTION_SEARCH_METADATA")))
}
