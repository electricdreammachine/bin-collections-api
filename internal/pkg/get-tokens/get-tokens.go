package gettokens

import (
	"strings"
	"regexp"
	"github.com/gocolly/colly"
	"bin-collections-api/internal/pkg/get-config-value"
)

// Tokens contains both a cookie and csrf id
type Tokens struct {
	Cookie []string
	InstanceID string
}

// GetTokens will make a network request for required ephemeral tokens
func GetTokens() <-chan Tokens {
	c := colly.NewCollector()
	channel := make(chan Tokens)

	c.OnHTML(getconfigvalue.ByKey("TOKEN_ELEMENT"), func(e *colly.HTMLElement)  {
		headers := e.Response.Headers.Get(getconfigvalue.ByKey("TOKEN_HEADER"))
		instanceID := e.Attr(getconfigvalue.ByKey("TOKEN_ELEMENT_ATTRIBUTE"))
		appHeaderRegex := regexp.MustCompile(
			getconfigvalue.ByKey("TOKEN_REGEX"),
		)

		matchingSubGroups := appHeaderRegex.Find([]byte(headers))

		go func() {
			channel <- Tokens{
				strings.Split(string(matchingSubGroups), "="),
				instanceID,
			}

			close(channel)
		}()
	})

	c.Visit(getconfigvalue.ByKey("TOKEN_URL"))

	return channel
}