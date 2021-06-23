package scraperservice

import (
	models "bin-collections-api/internal/models"
	config "bin-collections-api/internal/services/config"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
	"github.com/tidwall/sjson"
)

func Submit(additionalValues []models.MetaDataItem) <-chan models.RedirectMetaData {
	c := colly.NewCollector()
	channel := make(chan models.RedirectMetaData)
	requiredMetaData := <-GetTokens()
	var flowChangeResponse map[string]string

	data := make(map[string]string)

	for _, v := range append(requiredMetaData.MetaData, Populate(nil, additionalValues)...) {
		if len(v.Path) == 0 {
			data[v.Name] = v.Value.(string)
		} else {
			pathSeparatedByParent := strings.SplitN(v.Path, ".", 2)
			
			if (v.Value == "[]") {
				value, _ := v.Value.(string)
				data[pathSeparatedByParent[0]], _ = sjson.SetRaw(data[pathSeparatedByParent[0]], pathSeparatedByParent[1], value)
			} else {
				data[pathSeparatedByParent[0]], _ = sjson.Set(data[pathSeparatedByParent[0]], pathSeparatedByParent[1], v.Value)
			}
		}
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	})

	c.OnResponse(func(e *colly.Response) {
		json.Unmarshal(e.Body, &flowChangeResponse)

		go func() {
			channel <- models.RedirectMetaData{
				Cookie:      requiredMetaData.Cookie,
				RedirectURL: flowChangeResponse["redirectURL"],
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
