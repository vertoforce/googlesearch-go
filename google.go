package google

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"

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
	defer body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}
	html, err := doc.Html()
	ioutil.WriteFile("out2.html", []byte(html), 0644)

	results := []Result{}

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

	return results, nil
}
