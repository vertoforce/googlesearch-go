# Google Search Go

This library fetches google results programmatically

## Usage

```go
search := Search{
    Q: "Test",
}

results, _ := Query(context.Background(), &search)
if len(results) > 0 {
    fmt.Println(results[0].Title)
}

// Output: Speedtest by Ookla - The Global Broadband Speed Test
```

## How it works

It makes a request to google (`https://google.com/search?q=`) with the search query.  It uses a random chrome user-agent.

It loops over each `div .g` divs with the following jquery parse commands:

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