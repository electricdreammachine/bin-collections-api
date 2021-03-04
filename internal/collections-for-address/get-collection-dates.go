package getcollectiondates

import (
	"bin-collections-api/internal/models"
	config "bin-collections-api/internal/services/config"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// ForUniqueAddressID gets all available collection dates for a single address
func ForUniqueAddressID(collector colly.Collector) <-chan models.Collections {
	collectionsChannel := make(chan models.Collections)
	collectionTypesChannel := make(chan models.CollectionColourRegistry)

	collector.OnHTML(config.ByKey("KEY_ELEMENT"), func(e *colly.HTMLElement) {
		spacePattern := regexp.MustCompile(`\s|\p{Z}`)
		keyText := spacePattern.ReplaceAllString(e.Text, "")
		keyRegex := regexp.MustCompile(
			config.ByKey("KEY_REGEX"),
		)

		if keyRegex.Match([]byte(keyText)) {
			cells := e.DOM.Find("td")
			cellGroupSize, _ := strconv.Atoi(config.ByKey("KEY_CELLS_GROUP_SIZE"))
			var collectionTypes []models.CollectionColourRegistryEntry

			for i := 0; i < cells.Length(); i = i + cellGroupSize {
				targetAttribute, _ := cells.Eq(i).Children().Eq(0).Attr("style")
				fmt.Println(targetAttribute)
				colourRegex := regexp.MustCompile(
					config.ByKey("COLOR_FIND_REGEX"),
				)
				collectionTypes = append(collectionTypes,
					models.CollectionColourRegistryEntry{
						Colour:   string(colourRegex.Find([]byte(targetAttribute))),
						TypeName: strings.Title(strings.ToLower(strings.TrimSpace(cells.Eq(i + 1).Text()))),
					},
				)
			}

			go func() {
				collectionTypesChannel <- models.NewCollectionColourRegistry(collectionTypes)

				close(collectionTypesChannel)
			}()
		}
	})

	collector.OnHTML(config.ByKey("DATES_ELEMENT"), func(e *colly.HTMLElement) {
		fmt.Println("swag")
		scriptText := e.Text
		datesArrayRegex := regexp.MustCompile(
			config.ByKey("DATES_REGEX"),
		)

		delimitedDateGroups := strings.ReplaceAll(datesArrayRegex.FindStringSubmatch(scriptText)[1], "\"", "")
		splitRegex := regexp.MustCompile(",{2,}")
		unprocessedDateValues := splitRegex.Split(delimitedDateGroups, -1)
		types := <-collectionTypesChannel
		var dates models.Collections

		for _, v := range unprocessedDateValues {
			dateValueParts := strings.Split(v, ",")

			if (len(dateValueParts)) == 1 {
				continue
			}

			date, unprocessedTypes := dateValueParts[0], dateValueParts[1:]
			var typeIndices []string
			for _, unprocessedType := range unprocessedTypes {
				if len(strings.TrimSpace(unprocessedType)) > 0 {
					typeIndices = append(typeIndices, strconv.Itoa(types[unprocessedType].Index))
				}
			}

			dates.Dates = append(dates.Dates,
				models.Collection{
					Type: typeIndices,
					Date: date,
				},
			)
		}

		typesByIndex := make(map[string]string)

		for _, typeByIndex := range types {
			typesByIndex[strconv.Itoa(typeByIndex.Index)] = typeByIndex.TypeName
		}

		dates.Types = typesByIndex

		go func() {
			collectionsChannel <- dates

			close(collectionsChannel)
		}()
	})

	return collectionsChannel
}
