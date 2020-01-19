package google

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// File with different ways to parse depending on what google returns

// parserFunction is a function that can parse google results
type parserFunction func(doc *goquery.Document) []Result

var defaultParserFunctions = []parserFunction{
	parseMethod1,
	parsePingMethod,
}

// TODO: Add separate entry for SEO links on same entry
// TODO: Add cached page link

func parseMethod1(doc *goquery.Document) (results []Result) {
	doc.Find("div .g").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Find("a").Attr("href")
		if !exists {
			return
		}

		// Title
		title := s.Find("a").Find("h3").Text()
		// Remove extra spacing
		title = regexp.MustCompile(`\s+`).ReplaceAllString(title, " ")

		// Description
		description := s.Find("div .s").Find("div:has(:not(div))").Text()
		// Remove extra spacing
		description = regexp.MustCompile(`\s+`).ReplaceAllString(description, " ")

		result := Result{
			Title:       title,
			URL:         link,
			Description: description,
		}

		results = append(results, result)
	})

	return results
}

func parsePingMethod(doc *goquery.Document) (results []Result) {
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		googlePath, exists := s.Attr("ping")
		if !exists {
			return
		}

		// Parse the URL parameters
		vals, err := url.ParseQuery(googlePath)
		if err != nil {
			return
		}

		// Find the "url" parameter
		if URLs, ok := vals["url"]; ok {
			for _, URL := range URLs {

				// Ignore google URLs, they usually are just buttons on the page or something
				if googleRegex.MatchString(URL) {
					return
				}

				// Ignore URLs starting with / or "#"
				if strings.Index(URL, "/") == 0 || URL == "#" {
					return
				}

				// Get the heading and replace extra spaces
				title := s.Parent().Find("div:not(:has(span))").Text()
				// replace extra spaces and newlines
				title = regexp.MustCompile(`\s+`).ReplaceAllString(title, " ")
				title = strings.ReplaceAll(title, "\n", "")

				// Description is the span that does not have the link (cached link)
				description := s.Parent().Parent().Parent().Find("span:not(:has(a))").Text()
				// replace extra spaces and newlines
				description = regexp.MustCompile(`\s+`).ReplaceAllString(description, " ")
				description = strings.ReplaceAll(description, "\n", "")

				result := Result{
					URL:         URL,
					Title:       title,
					Description: description,
				}
				results = append(results, result)
			}
		}
	})

	return results
}
