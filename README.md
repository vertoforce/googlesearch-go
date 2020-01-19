# Google Search Go

[![Go Report Card](https://goreportcard.com/badge/github.com/vertoforce/googlesearch-go)](https://goreportcard.com/report/github.com/vertoforce/googlesearch-go)
[![Documentation](https://godoc.org/github.com/vertoforce/googlesearch-go?status.svg)](https://godoc.org/github.com/vertoforce/googlesearch-go)

This library fetches google results programmatically

## Usage

```go
results, _ := Query(context.Background(), &Search{Q: "Test"})
if len(results) > 0 {
    fmt.Println(results[0].Title)
}

// Output: Speedtest by Ookla - The Global Broadband Speed Test
```

### Usage with TryHard

Setting `TryHard=true` causes the query to keep trying proxies to get a `200 OK`.
It uses the [proxier library](https://github.com/vertoforce/proxier).

```go
results, _ := Query(context.Background(), &Search{Q: "Test", TryHard: true})
```

### Continuous results (going through google search result pages)

The library also supports navigating through google pages to get more results.  Make a continuous query like the following:

```go
results := QueryContinuous(context.Background(), &Search{Q: "Test"}, time.Second*10, time.Second*5)
for result := range results {
    fmt.Println(result.Title)
}
```

Note that google usually only returns approximately the first ~400 results, so it can not continue beyond that.

## How it works

It makes a request to google (`https://google.com/search?q=`) with the search query.  It uses a random chrome user-agent.

Google sometimes returns different types of responses, so it tries multiple parsing functions in `parse.go` to find one that returns results successfully.

The default method is it loops over each `div .g` divs with the following jquery parse commands:

- Title: `Find("a").Find('h3")`
- Description: `Find("div .s").Find("div:has(:not(div))")`.
- URL: `Find("a").Attr("href")`

## Why

All the go search libraries I found were one of the following

- Poorly written
- Poorly documented
- Paid
- Only get URLs
- Old / no longer worked (last commit years ago)

### Other libraries

- https://github.com/serpapi/google-search-results-golang (paid)
- https://github.com/henkman/google
- https://github.com/anners/image-search
- https://github.com/schollz/googleit
- https://github.com/nakabonne/WebCrawlerForSerps
- https://github.com/keyle/googlesearch
- https://github.com/WithGJR/google-search-parser
- https://github.com/captainju/goGoogleSearch
