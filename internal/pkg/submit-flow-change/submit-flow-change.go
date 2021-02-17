package submitflowchange

import (
	getconfigvalue "bin-collections-api/internal/pkg/get-config-value"
	getinpagemetadata "bin-collections-api/internal/pkg/get-in-page-metadata"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
	"github.com/tidwall/sjson"
)

type redirectMetaData struct {
	Cookie      []string
	RedirectUrl string
}

// Submit makes submit request
func Submit(additionalValues []getinpagemetadata.MetaDataItem) <-chan redirectMetaData {
	c := colly.NewCollector()
	channel := make(chan redirectMetaData)
	requiredMetaData := <-getinpagemetadata.GetTokens()
	var flowChangeResponse map[string]string

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
		json.Unmarshal(e.Body, &flowChangeResponse)

		go func() {
			channel <- redirectMetaData{
				Cookie:      requiredMetaData.Cookie,
				RedirectUrl: flowChangeResponse["redirectURL"],
			}
		}()
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

	return channel
}
