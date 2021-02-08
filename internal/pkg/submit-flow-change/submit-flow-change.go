package submitflowchange

import (
	getinpagemetadata "bin-collections-api/internal/pkg/get-in-page-metadata"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tidwall/sjson"
)

// Submit makes submit request
func Submit(additionalValues []getinpagemetadata.MetaDataItem) <-chan []string {
	channel := make(chan []string)
	requiredMetaData := <-getinpagemetadata.GetTokens()
	client := &http.Client{}

	data := ""

	fmt.Println(additionalValues)

	fmt.Println(append(requiredMetaData.MetaData, additionalValues...))

	for _, v := range append(requiredMetaData.MetaData, getinpagemetadata.Populate(nil, additionalValues)...) {
		data, _ = sjson.Set(data, v.Path, v.Value)
		fmt.Println(data)
	}

	req, err := http.NewRequest(
		"POST",
		"https://iweb.itouchvision.com/portal/wwv_flow.accept",
		bytes.NewBuffer([]byte(data)),
	)

	if err != nil {
		log.Fatal("woops")
	}

	req.Header.Add(requiredMetaData.Cookie[0], requiredMetaData.Cookie[1])

	resp, _ := client.Do(req)

	defer resp.Body.Close()

	ioutil.ReadAll(resp.Body)

	// log.Println(string(body))

	go func() {
		channel <- requiredMetaData.Cookie
	}()

	return channel
}
