package google

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
	"github.com/vertoforce/proxier"
)

const (
	baseURL = "https://google.com/search?q=%s&num=%d&start=%d"
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
	// Result number to start at
	Start int
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

	if search.Num == 0 {
		search.Num = 10
	}

	// Build request
	urlString := fmt.Sprintf(baseURL, url.QueryEscape(search.Q), search.Num, search.Start)
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

	results := []Result{}
	for _, parserMethod := range defaultParserFunctions {
		results = parserMethod(doc)
		if len(results) != 0 {
			return results, nil
		}
	}

	return results, nil
}
