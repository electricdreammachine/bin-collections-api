package submitflowchange

import (
	getconfigvalue "bin-collections-api/internal/pkg/get-config-value"
	getinpagemetadata "bin-collections-api/internal/pkg/get-in-page-metadata"
	"bytes"
	"fmt"
	"net/http"

	"github.com/gocolly/colly"
	"github.com/tidwall/sjson"
)

// Submit makes submit request
func Submit(additionalValues []getinpagemetadata.MetaDataItem) <-chan []string {
	c := colly.NewCollector()
	channel := make(chan []string)
	requiredMetaData := <-getinpagemetadata.GetTokens()

	data := make(map[string]string)

	fmt.Println(additionalValues)

	fmt.Println(append(requiredMetaData.MetaData, additionalValues...))

	for _, v := range append(requiredMetaData.MetaData, getinpagemetadata.Populate(nil, additionalValues)...) {
		data, _ = sjson.Set(data, v.Path, v.Value)
		fmt.Println(data)
	}

	c.OnResponse(func(e *colly.Response) {
		fmt.Println(e)
	})

	c.SetCookies(
		getconfigvalue.ByKey("DATES_COOKIE_DOMAIN"),
		[]*http.Cookie{
			&http.Cookie{
				Name:  requiredMetaData.Cookie[0],
				Value: requiredMetaData.Cookie[1],
			},
		},
	)

	c.Post(
		"https://iweb.itouchvision.com/portal/wwv_flow.accept",
		bytes.NewBuffer([]byte(data)),
	)

	go func() {
		channel <- requiredMetaData.Cookie
	}()

	return channel
}
