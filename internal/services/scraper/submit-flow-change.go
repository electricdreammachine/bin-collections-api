package scraperservice

import (
	models "bin-collections-api/internal/models"
	config "bin-collections-api/internal/services/config"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
	"github.com/tidwall/sjson"
)

// Submit makes submit request
func Submit(additionalValues []models.MetaDataItem) <-chan models.RedirectMetaData {
	c := colly.NewCollector()
	channel := make(chan models.RedirectMetaData)
	requiredMetaData := <-GetTokens()
	var flowChangeResponse map[string]string

	data := make(map[string]string)

	fmt.Println(requiredMetaData)

	for _, v := range append(requiredMetaData.MetaData, Populate(nil, additionalValues)...) {
		if len(v.Path) == 0 {
			data[v.Name] = v.Value.(string)
		} else {
			pathSeparatedByParent := strings.SplitN(v.Path, ".", 2)
			data[pathSeparatedByParent[0]], _ = sjson.Set(data[pathSeparatedByParent[0]], pathSeparatedByParent[1], v.Value)
		}
	}

	fmt.Println(data)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	})

	c.OnResponse(func(e *colly.Response) {
		json.Unmarshal(e.Body, &flowChangeResponse)
		fmt.Println(flowChangeResponse)

		go func() {
			channel <- models.RedirectMetaData{
				Cookie:      requiredMetaData.Cookie,
				RedirectUrl: flowChangeResponse["redirectURL"],
			}
		}()
	})

	c.SetCookies(
		config.ByKey("DATES_COOKIE_DOMAIN"),
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
