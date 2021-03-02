package getcollectiondates

import (
	collectiontypes "bin-collections-api/internal/pkg/collection-types"
	getconfigvalue "bin-collections-api/internal/pkg/get-config-value"
	getinpagemetadata "bin-collections-api/internal/pkg/get-in-page-metadata"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// Collections describes the collections available in the date range
type Collections struct {
	Dates []collection      `json:"dates"`
	Types map[string]string `json:"types"`
}

type collection struct {
	Type []string `json:"type"`
	Date string   `json:"date"`
}

// ForUniqueAddressID gets all available collection dates for a single address
func ForUniqueAddressID(url string, cookie getinpagemetadata.Cookie) <-chan Collections {
	c := colly.NewCollector()
	collectionsChannel := make(chan Collections)
	collectionTypesChannel := make(chan collectiontypes.CollectionColourRegistry)

	c.OnHTML("body", func(e *colly.HTMLElement) {
		// fmt.Println(e)
	})

	c.OnHTML(getconfigvalue.ByKey("KEY_ELEMENT"), func(e *colly.HTMLElement) {
		spacePattern := regexp.MustCompile(`\s|\p{Z}`)
		keyText := spacePattern.ReplaceAllString(e.Text, "")
		keyRegex := regexp.MustCompile(
			getconfigvalue.ByKey("KEY_REGEX"),
		)

		if keyRegex.Match([]byte(keyText)) {
			cells := e.DOM.Find("td")
			cellGroupSize, _ := strconv.Atoi(getconfigvalue.ByKey("KEY_CELLS_GROUP_SIZE"))
			var collectionTypes []collectiontypes.CollectionColourRegistryEntry

			for i := 0; i < cells.Length(); i = i + cellGroupSize {
				targetAttribute, _ := cells.Eq(i).Children().Eq(0).Attr("style")
				fmt.Println(targetAttribute)
				colourRegex := regexp.MustCompile(
					getconfigvalue.ByKey("COLOR_FIND_REGEX"),
				)
				collectionTypes = append(collectionTypes,
					collectiontypes.CollectionColourRegistryEntry{
						Colour:   string(colourRegex.Find([]byte(targetAttribute))),
						TypeName: strings.Title(strings.ToLower(strings.TrimSpace(cells.Eq(i + 1).Text()))),
					},
				)
			}

			go func() {
				collectionTypesChannel <- collectiontypes.NewCollectionColourRegistry(collectionTypes)

				close(collectionTypesChannel)
			}()
		}
	})

	c.OnHTML(getconfigvalue.ByKey("DATES_ELEMENT"), func(e *colly.HTMLElement) {
		scriptText := e.Text
		datesArrayRegex := regexp.MustCompile(
			getconfigvalue.ByKey("DATES_REGEX"),
		)

		delimitedDateGroups := strings.ReplaceAll(datesArrayRegex.FindStringSubmatch(scriptText)[1], "\"", "")
		splitRegex := regexp.MustCompile(",{2,}")
		unprocessedDateValues := splitRegex.Split(delimitedDateGroups, -1)
		// cellGroupSize, _ := strconv.Atoi(getconfigvalue.ByKey("DATES_GROUP_SIZE"))
		types := <-collectionTypesChannel
		var dates Collections

		fmt.Println(unprocessedDateValues)

		for i := 0; i < len(unprocessedDateValues); i = i + 1 {
			fmt.Println(unprocessedDateValues[i])
			unprocessedTypes := []string{unprocessedDateValues[i+1], unprocessedDateValues[i+2]}
			var typeIndices []string
			for _, unprocessedType := range unprocessedTypes {
				if len(strings.TrimSpace(unprocessedType[1:len(unprocessedType)-1])) > 0 {
					typeIndices = append(typeIndices, strconv.Itoa(types[unprocessedType[1:len(unprocessedType)-1]].Index))
				}
			}

			dates.Dates = append(dates.Dates,
				collection{
					Type: typeIndices,
					Date: unprocessedDateValues[i][1 : len(unprocessedDateValues[i])-1],
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

	return collectionsChannel
}
