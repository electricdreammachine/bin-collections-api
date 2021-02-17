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
func ForPostCode(url string, cookie getinpagemetadata.Cookie) <-chan []Address {
	fmt.Println(cookie)
	c := colly.NewCollector()
	addressesChannel := make(chan []Address)
	fmt.Println(fmt.Sprintf("%v%v", getconfigvalue.ByKey("DATES_COOKIE_DOMAIN"), url))

	c.OnHTML("body", func(e *colly.HTMLElement) {
		// fmt.Println(e)
	})

	c.OnHTML(getconfigvalue.ByKey("ADDRESSES_SEARCH"), func(e *colly.HTMLElement) {
		fmt.Println("yes")

		go func() {
			addressesChannel <- []Address{nil}

			close(addressesChannel)
		}()
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(c.Cookies(getconfigvalue.ByKey("DATES_COOKIE_DOMAIN")))
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(string(r.Body))
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(err)
		fmt.Println(r)
	})

	c.SetCookies(
		getconfigvalue.ByKey("DATES_COOKIE_DOMAIN"),
		[]*http.Cookie{
			{
				Name:  cookie[0],
				Value: cookie[1],
			},
		},
	)

	c.Visit(fmt.Sprintf("%v/portal/%v", getconfigvalue.ByKey("DATES_COOKIE_DOMAIN"), url))

	return addressesChannel
}
