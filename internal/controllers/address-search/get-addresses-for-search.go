package addresssearch

import (
	"bin-collections-api/internal/models"
	config "bin-collections-api/internal/services/config"

	"github.com/gocolly/colly"
)

func ForPostCode(collector *colly.Collector) <-chan interface{} {
	addressesChannel := make(chan interface{})
	var addresses []models.Address

	collector.OnHTML(config.ByKey("ADDRESSES_SEARCH"), func(e *colly.HTMLElement) {
		addresses = append(addresses, models.Address{
			AddressStr: e.Text,
			ID:         e.Attr("value"),
		})
	})

	collector.OnScraped(func(r *colly.Response) {
		go func() {
			addressesChannel <- addresses

			close(addressesChannel)
		}()
	})

	return addressesChannel
}
