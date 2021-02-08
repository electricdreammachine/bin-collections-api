package addresssearch

import (
	getconfigvalue "bin-collections-api/internal/pkg/get-config-value"
	getinpagemetadata "bin-collections-api/internal/pkg/get-in-page-metadata"
	"fmt"
	"net/http"

	"github.com/gocolly/colly"
)

type Address map[string]interface{}

// ForPostCode gets all available collection dates for a single address
func ForPostCode(cookie getinpagemetadata.Cookie) <-chan []Address {
	fmt.Println(cookie)
	c := colly.NewCollector()
	addressesChannel := make(chan []Address)

	c.OnHTML("body", func(e *colly.HTMLElement) {
		fmt.Println(e)
	})

	c.OnHTML(getconfigvalue.ByKey("ADDRESSES_SEARCH"), func(e *colly.HTMLElement) {
		fmt.Println(e.Text)

		go func() {
			addressesChannel <- []Address{nil}

			close(addressesChannel)
		}()
	})

	c.SetCookies(
		getconfigvalue.ByKey("DATES_COOKIE_DOMAIN"),
		[]*http.Cookie{
			&http.Cookie{
				Name:  cookie[0],
				Value: cookie[1],
			},
		},
	)

	c.Visit(fmt.Sprintf(getconfigvalue.ByKey("DATES_URL")))

	return addressesChannel
}
