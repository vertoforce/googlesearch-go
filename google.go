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
	"github.com/vertoforce/proxier"
)

const (
	baseURL = "https://google.com/search?q=%s&num=%d"
)

var (
	googleRegex = regexp.MustCompile(`https?:\/\/(\w*\.)?google`)
)

// Search A google search
type Search struct {
	// Query string
	Q string
	// Number of results per page
	Num int
	// TryHard if sent to true will keep trying proxies until we get a 200 OK response.
	// It uses the proxier library (github.com/vertoforce/proxier)
	TryHard bool
}

// Result from google
type Result struct {
	URL         string
	Title       string
	Description string
}

var p *proxier.Proxier

// Query makes a google query
func Query(ctx context.Context, search *Search) ([]Result, error) {
	// Get fake user agent
	fakeUserAgent := browser.Chrome()

	// Build request
	urlString := fmt.Sprintf(baseURL, url.QueryEscape(search.Q), search.Num)
	req, err := http.NewRequestWithContext(ctx, "GET", urlString, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", fakeUserAgent)
	var resp *http.Response

	// Do request
	if search.TryHard {
		// Use a proxy if default method fails
		if p == nil {
			p = proxier.New()
			p.TryNoProxyFirst = true
		}
		resp, err = p.DoRequest(ctx, req)
	} else {
		resp, err = http.DefaultClient.Do(req)
	}
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid response code: %d", resp.StatusCode)
	}

	return parse(resp.Body)
}

// parse a google body
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
