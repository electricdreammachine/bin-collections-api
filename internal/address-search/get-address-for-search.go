package addresssearch

import (
	"bin-collections-api/internal/models"
	config "bin-collections-api/internal/services/config"
	"fmt"
	"net/http"

	"github.com/gocolly/colly"
)

type Address map[string]interface{}

// ForPostCode gets all available collection dates for a single address
func ForPostCode(url string, cookie models.Cookie) <-chan []Address {
	fmt.Println(cookie)
	c := colly.NewCollector()
	addressesChannel := make(chan []Address)
	fmt.Println(fmt.Sprintf("%v%v", config.ByKey("DATES_COOKIE_DOMAIN"), url))

	c.OnHTML("body", func(e *colly.HTMLElement) {
		// fmt.Println(e)
	})

	c.OnHTML(config.ByKey("ADDRESSES_SEARCH"), func(e *colly.HTMLElement) {
		fmt.Println("yes")

		go func() {
			addressesChannel <- []Address{nil}

			close(addressesChannel)
		}()
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(c.Cookies(config.ByKey("DATES_COOKIE_DOMAIN")))
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(string(r.Body))
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(err)
		fmt.Println(r)
	})

	c.SetCookies(
		config.ByKey("DATES_COOKIE_DOMAIN"),
		[]*http.Cookie{
			{
				Name:  cookie[0],
				Value: cookie[1],
			},
		},
	)

	c.Visit(fmt.Sprintf("%v/portal/%v", config.ByKey("DATES_COOKIE_DOMAIN"), url))

	return addressesChannel
}
