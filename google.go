package google

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	browser "github.com/EDDYCJY/fake-useragent"
)

const (
	baseURL = "https://google.com/search?q=%s"
)

// Search A google search
type Search struct {
	Q string
}

// Result Result from google
type Result struct {
	URL string
}

// Query Make a google query
func Query(ctx context.Context, search *Search) ([]Result, error) {
	// Get fake user agent
	fakeUserAgent := browser.Random()

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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ioutil.WriteFile("out.html", body, 0644)

	return []Result{}, nil
}
