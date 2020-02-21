package getcollectiondates

import (
	"github.com/gocolly/colly"
	"fmt"
	"regexp"
	"net/http"
	"bin-collections-api/internal/pkg/get-tokens"
	"bin-collections-api/internal/pkg/get-config-value"
)

// Dates describes the collections available in the date range
type Dates struct {
	Collections []collection
}

type collection struct {
	collectionType string
	collectionDate string
}

// ForUniqueAddressID gets all available collection dates for a single address 
func ForUniqueAddressID(t gettokens.Tokens, uniqueAddressID string) <-chan Dates {
	c := colly.NewCollector()
	channel := make(chan Dates)

	c.OnHTML(getconfigvalue.ByKey("DATES_ELEMENT"), func(e *colly.HTMLElement) {
		scriptText := e.Text
		datesArrayRegex := regexp.MustCompile(
			getconfigvalue.ByKey("DATES_REGEX"),
		)

		dates := datesArrayRegex.FindStringSubmatch(scriptText)[1]

		//TODO format and push this back through the channel
		fmt.Println(dates)
	})

	c.SetCookies(
		getconfigvalue.ByKey("DATES_COOKIE_DOMAIN"),
		[]*http.Cookie {
			&http.Cookie{
				Name: t.Cookie[0],
				Value: t.Cookie[1],
			},
		},
	)
	
	c.Visit(fmt.Sprintf(getconfigvalue.ByKey("DATES_URL"), t.InstanceID, uniqueAddressID))

	return channel
}