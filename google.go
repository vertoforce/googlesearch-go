package google

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
)

const (
	baseURL = "https://google.com/search?q=%s"
)

var (
	googleRegex = regexp.MustCompile(`https?:\/\/(\w*\.)?google`)
)

// Search A google search
type Search struct {
	Q string
}

// Result Result from google
type Result struct {
	URL         string
	Title       string
	Description string
}

// Query Make a google query
func Query(ctx context.Context, search *Search) ([]Result, error) {
	// Get fake user agent
	fakeUserAgent := browser.Chrome()

	// Build request
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf(baseURL, url.QueryEscape(search.Q)), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", fakeUserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid response code: %d", resp.StatusCode)
	}

	return parse(resp.Body)
}

// parse Parse a google body
func parse(body io.ReadCloser) ([]Result, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}
	html, err := doc.Html()
	ioutil.WriteFile("out2.html", []byte(html), 0644)

	results := []Result{}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		URL, exists := s.Attr("href")
		if !exists {
			return
		}

		// Ignore google URLs, they usually are just buttons on the page or something
		if googleRegex.MatchString(URL) {
			return
		}

		// Ignore URLs starting with / or "#"
		if strings.Index(URL, "/") == 0 || URL == "#" {
			return
		}

		// Get the heading and replace extra spaces
		title := s.Find("[role='heading']").Text()
		// replace extra spaces and newlines
		title = regexp.MustCompile(`\s+`).ReplaceAllString(title, " ")
		title = strings.ReplaceAll(title, "\n", "")

		// Get description is probably the div that has "..." but does NOT have the > character, and does not have child divs
		description := s.Parent().Parent().Parent().Find("div:contains('...'):not(:contains('›')):not(:has(div))").Text()
		// replace extra spaces and newlines
		description = regexp.MustCompile(`\s+`).ReplaceAllString(description, " ")
		description = strings.ReplaceAll(description, "\n", "")

		result := Result{
			URL:         URL,
			Title:       title,
			Description: description,
		}
		results = append(results, result)
	})

	return results, nil
}
