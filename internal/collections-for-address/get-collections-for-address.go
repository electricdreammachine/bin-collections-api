// package getcollectiondates

// import (
// 	"bin-collections-api/internal/models"
// 	config "bin-collections-api/internal/services/config"
// 	submitflowchange "bin-collections-api/internal/services/scraper"
// 	"encoding/json"
// 	"log"
// 	"net/http"
// )

// type collectionSearchRequest map[string]interface{}

// type additionalDataSchema struct {
// 	Values []models.MetaDataItem `json:"values"`
// }

// // GetCollectionsForID blocks until tokens retrieved and then retrieves collections for the given ID
// func GetCollectionsForID(w http.ResponseWriter, request *http.Request) {
// 	var parsedSchema additionalDataSchema
// 	metaDataSchema := config.ByKey("ADDITIONAL_COLLECTION_SEARCH_METADATA")
// 	err := json.Unmarshal([]byte(metaDataSchema), &parsedSchema)
// 	var collectionSearch collectionSearchRequest
// 	err2 := json.NewDecoder(request.Body).Decode(&collectionSearch)

// 	// fmt.Println("weewee")

// 	for i, v := range parsedSchema.Values {
// 		if len(v.ValueFromMap) > 0 {
// 			v.Value = collectionSearch[v.ValueFromMap]
// 		}
// 		parsedSchema.Values[i] = v
// 	}

// 	if err != nil {
// 		log.Fatal("err")
// 	}

// 	if err2 != nil {
// 		log.Fatal("err2")
// 	}

// 	submitResponse := <-submitflowchange.Submit(parsedSchema.Values)
// 	// fmt.Println(submitResponse)
// 	collectionDates := <-ForUniqueAddressID(submitResponse.RedirectUrl, submitResponse.Cookie)

// 	// fmt.Println(collectionDates)

// 	json.NewEncoder(w).Encode(collectionDates)
// }