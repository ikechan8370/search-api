package serp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
	"github.com/gocolly/colly/v2/queue"
	"net/url"
)

func SearchBingImage(ctx context.Context, query, ua string, limit int, proxyAddr string, verbose bool) ([]ImageResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	query = url.QueryEscape(query)
	c := colly.NewCollector(colly.MaxDepth(1))
	c.UserAgent = ua
	q, _ := queue.New(1, &queue.InMemoryQueueStorage{MaxSize: 10000})
	var results []ImageResult
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
	c.OnHTML("div.iuscp", func(e *colly.HTMLElement) {

		sel := e.DOM
		var m map[string]string
		mStr, _ := sel.Find("div.imgpt > a.iusc").Attr("m")
		err := json.Unmarshal([]byte(mStr), &m)
		if err != nil {
			fmt.Println("解析JSON字符串出错:", err)
			return
		}
		mUrl := m["murl"]
		PUrl := m["purl"]
		Cid := m["cid"]
		TUrl := m["turl"]
		t := m["t"]
		desc := m["desc"]
		src, _ := sel.Find("div.imgpt > a.iusc > div > img").Attr("src")
		label, _ := sel.Find("div.infopt > a.inflnk").Attr("aria-label")
		rank += 1
		if mUrl != "" && mUrl != "#" && src != "" {
			result := ImageResult{
				//Cid:   Cid,
				//PUrl:  PUrl,
				MUrl: mUrl,
				//Turl:  TUrl,
				//T:     t,
				//Desc:  desc,
				//Src:   src,
				Label: label,
			}
			if verbose {
				result.Cid = Cid
				result.PUrl = PUrl
				result.Turl = TUrl
				result.Desc = desc
				result.T = t
				result.Src = src
			}
			results = append(results, result)
			filteredRank += 1
		}
	})

	url := "https://www.bing.com/images/search?q=" + query
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

type ImageResult struct {
	Cid   string `json:"cid,omitempty"`
	PUrl  string `json:"purl,omitempty"`
	MUrl  string `json:"murl,omitempty"`
	Turl  string `json:"turl,omitempty"`
	T     string `json:"t,omitempty"`
	Desc  string `json:"desc,omitempty"`
	Src   string `json:"src,omitempty"`
	Label string `json:"label,omitempty"`
}
