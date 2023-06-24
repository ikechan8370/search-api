package main

import (
	"context"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
	"github.com/gocolly/colly/v2/queue"
	"net/url"
	"strings"
)

func SearchBing(ctx context.Context, query, ua string, limit int, proxyAddr string) ([]Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	query = url.QueryEscape(query)
	c := colly.NewCollector(colly.MaxDepth(1))
	c.UserAgent = ua
	q, _ := queue.New(1, &queue.InMemoryQueueStorage{MaxSize: 10000})
	var results []Result
	var rErr error
	filteredRank := 1
	rank := 1
	c.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil {
			r.Abort()
			rErr = err
			return
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		rErr = err
	})
	// https://www.w3schools.com/cssref/css_selectors.asp
	c.OnHTML("#b_results > li", func(e *colly.HTMLElement) {

		sel := e.DOM
		linkHref, _ := sel.Find("h2 > a").Attr("href")
		linkText := strings.TrimSpace(linkHref)
		titleText := strings.TrimSpace(sel.Find("h2 > a").Text())
		prefixType := sel.Find("div > p > span").Text()

		descText := strings.TrimSpace(sel.Find("div > p").Text())
		descText = strings.TrimPrefix(descText, prefixType)
		descText = strings.TrimPrefix(descText, " · ")
		rank += 1
		if linkText != "" && linkText != "#" && titleText != "" {
			result := Result{
				Rank:        filteredRank,
				URL:         linkText,
				Title:       titleText,
				Description: descText,
			}
			results = append(results, result)
			filteredRank += 1
		}
	})

	url := "https://www.bing.com/search?q=" + query
	if proxyAddr != "" {
		rp, err := proxy.RoundRobinProxySwitcher(proxyAddr)
		if err != nil {
			return nil, err
		}
		c.SetProxyFunc(rp)
	}
	q.AddURL(url)
	q.Run(c)
	if rErr != nil {
		return nil, rErr
	}

	// Reduce results to max limit
	if limit != 0 && len(results) > limit {
		return results[:limit], nil
	}

	return results, nil
}

// Result represents a single result from Bing Search.
type Result struct {

	// Rank is the order number of the search result.
	Rank int `json:"rank"`

	// URL of result.
	URL string `json:"url"`

	// Title of result.
	Title string `json:"title"`

	// Description of the result.
	Description string `json:"description"`
}
