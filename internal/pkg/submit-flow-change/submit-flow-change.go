package submitflowchange

import (
	"bin-collections-api/internal/pkg/get-in-page-metadata"
	"net/http"
	"bytes"
	"log"
	"github.com/tidwall/sjson"
	"io/ioutil"
)

// Submit makes submit request
func Submit(additionalValues []getinpagemetadata.MetaDataItem) <-chan bool {
	channel := make(chan bool)
	requiredMetaData := <- getinpagemetadata.GetTokens()
	client := &http.Client{} 

	for _, v := range additionalValues {
		requiredMetaData.MetaData = append(requiredMetaData.MetaData, v)
	}

	data := ""

	for _, v := range requiredMetaData.MetaData {
		data, _ = sjson.Set(data, v.Path, v.Value)
	}

	req, err := http.NewRequest(
		"POST",
		"https://iweb.itouchvision.com/portal/wwv_flow.accept",
		bytes.NewBuffer([]byte(data)),
	)

	if (err != nil) {
		log.Fatal("woops")
	}

	req.Header.Add(requiredMetaData.Cookie[0], requiredMetaData.Cookie[1])

	resp, _ := client.Do(req)

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	log.Println(string(body))

	return channel
}
