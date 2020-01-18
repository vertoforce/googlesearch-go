# Google Search Go

This library fetches google results programmatically

## How it works

Simply, it makes a request to good (`https://google.com/search?q=`) with the search query.  However there's some interesting "gotchas" to be aware of.
First of all, when making a request from a Chrome user-agent, it changes the HTML it returns.  It instead returns `ping=` attributes that we look for.
Right now the code only supports the chrome user-agent response, but a TODO is to add support for the other responses, making it so we can send more user agents and more requests.

Once we get the HTML, we parse out the urls, titles, description, etc.

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