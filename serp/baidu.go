package serp

import (
	"context"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func SearchBaidu(ctx context.Context, query string, ua string, limit int) ([]Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	query = url.QueryEscape(query)
	c := colly.NewCollector(colly.MaxDepth(1))
	c.UserAgent = ua

	cookie, err := getCookie()
	c.SetCookies("https://www.baidu.com", cookie)
	if err != nil {
		return nil, err
	}
	q, _ := queue.New(1, &queue.InMemoryQueueStorage{MaxSize: 10000})
	var results []Result
	var rErr error
	filteredRank := 1
	rank := 1
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-XA;q=0.7,ja-JP;q=0.6,ja;q=0.5,zh-TW;q=0.4")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")

		if err := ctx.Err(); err != nil {
			r.Abort()
			rErr = err
			return
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		println(err.Error())
		rErr = err
	})
	// https://www.w3schools.com/cssref/css_selectors.asp
	c.OnHTML("#content_left > div", func(e *colly.HTMLElement) {

		sel := e.DOM

		linkHref, _ := sel.Find("div > div > h3 > a").Attr("href")
		linkText := strings.TrimSpace(linkHref)
		titleText := strings.TrimSpace(sel.Find("div > div > h3 > a").Text())
		descText := strings.TrimSpace(sel.Find("div > div > div:nth-of-type(2)").Text())
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

	url := "https://www.baidu.com/s?wd=" + query
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

func getCookie() ([]*http.Cookie, error) {
	client := &http.Client{}

	// 创建一个Cookie Jar来存储Cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println("创建Cookie Jar失败:", err)
		return nil, err
	}
	client.Jar = jar
	urlBaidu := "https://www.baidu.com/s?wd=1"
	url, err := url.Parse(urlBaidu)
	if err != nil {
		return nil, err
	}
	// 创建一个GET请求
	req, err := http.NewRequest("GET", urlBaidu, nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return nil, err
	}
	req.Header.Set("User-Agent", DefaultUA)
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return nil, err
	}
	defer resp.Body.Close()
	println(resp.StatusCode)
	// 读取响应内容
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return nil, err
	}
	return jar.Cookies(url), nil
}
