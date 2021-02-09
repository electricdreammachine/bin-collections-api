package submitflowchange

import (
	getconfigvalue "bin-collections-api/internal/pkg/get-config-value"
	getinpagemetadata "bin-collections-api/internal/pkg/get-in-page-metadata"
	"fmt"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
	"github.com/tidwall/sjson"
)

// Submit makes submit request
func Submit(additionalValues []getinpagemetadata.MetaDataItem) <-chan []string {
	c := colly.NewCollector()
	channel := make(chan []string)
	requiredMetaData := <-getinpagemetadata.GetTokens()

	data := make(map[string]string)

	for _, v := range append(requiredMetaData.MetaData, getinpagemetadata.Populate(nil, additionalValues)...) {
		if len(v.Path) == 0 {
			data[v.Name] = v.Value.(string)
		} else {
			pathSeparatedByParent := strings.SplitN(v.Path, ".", 2)
			data[pathSeparatedByParent[0]], _ = sjson.Set(data[pathSeparatedByParent[0]], pathSeparatedByParent[1], v.Value)
		}
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	})

	c.OnResponse(func(e *colly.Response) {
		fmt.Println(e)
	})

	c.SetCookies(
		getconfigvalue.ByKey("DATES_COOKIE_DOMAIN"),
		[]*http.Cookie{
			{
				Name:  requiredMetaData.Cookie[0],
				Value: requiredMetaData.Cookie[1],
			},
		},
	)

	c.Post(
		"https://iweb.itouchvision.com/portal/wwv_flow.accept",
		data,
	)

	go func() {
		channel <- requiredMetaData.Cookie
	}()

	return channel
}
