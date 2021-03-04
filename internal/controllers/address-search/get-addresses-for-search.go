package addresssearch

import (
	"bin-collections-api/internal/models"
	config "bin-collections-api/internal/services/config"
	"fmt"

	"github.com/gocolly/colly"
)

func ForPostCode(collector *colly.Collector) <-chan interface{} {
	addressesChannel := make(chan interface{})

	collector.OnHTML(config.ByKey("ADDRESSES_SEARCH"), func(e *colly.HTMLElement) {
		fmt.Println(e)
		go func() {
			addressesChannel <- []models.Address{}

			close(addressesChannel)
		}()
	})

	return addressesChannel
}
