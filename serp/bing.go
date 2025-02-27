package serp

import (
	"context"
	"errors"
	"github.com/Danny-Dasilva/CycleTLS/cycletls"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/url"
	"strings"
	"time"
)

var BingCookies []cycletls.Cookie

func SearchBing(myJa3 string, ctx context.Context, query, ua string, limit int, proxyAddr string, cookies []cycletls.Cookie, retry int) ([]Result, error) {
	if retry < 0 {
		return nil, errors.New("retry is invalid")
	}
	if myJa3 == "" {
		myJa3 = "771,4865-4867-4866-49195-49199-52393-52392-49196-49200-49162-49161-49171-49172-51-57-47-53-10,0-23-65281-10-11-35-16-5-51-43-13-45-28-21,29-23-24-25-256-257,0"
	}
	if cookies == nil {
		cookies = BingCookies
	}
	client := cycletls.Init()
	url := "https://www.bing.com/search?q=" + url.QueryEscape(query) + "&mkt=zh-CN"
	headers := map[string]string{
		"accept-language":           "zh-CN,zh;q=0.9,en;q=0.8,en-XA;q=0.7,ja-JP;q=0.6,ja;q=0.5,zh-TW;q=0.4",
		"cache-control":             "no-cache",
		"connection":                "keep-alive",
		"dnt":                       "1",
		"pragma":                    "no-cache",
		"sec-ch-ua":                 "\"Not(A:Brand\";v=\"99\", \"Google Chrome\";v=\"133\", \"Chromium\";v=\"133\"",
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        "\"Windows\"",
		"sec-fetch-dest":            "document",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-site":            "same-origin",
		"sec-fetch-user":            "?1",
		"upgrade-insecure-requests": "1",
		"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"accept-encoding":           "gzip, deflate, br, zstd",
		"ect":                       "4g",
	}
	response, err := client.Do(url, cycletls.Options{
		Body:      "",
		Ja3:       myJa3,
		UserAgent: ua,
		Headers:   headers,
		Cookies:   cookies,
	}, "GET")
	if err != nil {
		log.Print("Request Failed: " + err.Error())
	}
	for s := range response.Headers {
		if strings.ToLower(s) == "set-cookie" {
			setCookie := response.Headers[s]
			cookies, _ = ParseSetCookies(setCookie)
			BingCookies = cookies
		}
		//println(s, response.Headers[s])
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.Body))
	if err != nil {
		log.Fatal(err)
	}
	var results []Result
	doc.Find("#b_results > li").Each(func(i int, sel *goquery.Selection) {
		//println(sel.Text())
		linkHref, _ := sel.Find("h2 > a").Attr("href")
		linkText := strings.TrimSpace(linkHref)
		titleText := strings.TrimSpace(sel.Find("h2 > a").Text())
		prefixType := sel.Find("div > p > span").Text()

		descText := strings.TrimSpace(sel.Find("div > p").Text())
		descText = strings.TrimPrefix(descText, prefixType)
		descText = strings.TrimPrefix(descText, " · ")
		if linkText != "" && linkText != "#" && titleText != "" {
			result := Result{
				URL:         linkText,
				Title:       titleText,
				Description: descText,
			}
			results = append(results, result)
		}
	})
	if results == nil {
		retry = retry - 1
		return SearchBing(myJa3, ctx, query, ua, limit, proxyAddr, cookies, retry)
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

func ParseSetCookies(setCookies string) ([]cycletls.Cookie, error) {
	var cookies []cycletls.Cookie
	cookieParts := strings.Split(setCookies, "/,")
	for _, part := range cookieParts {
		var cookie cycletls.Cookie

		// Split each cookie part by semicolon
		attributes := strings.Split(part, ";")
		for _, attr := range attributes {
			attr = strings.TrimSpace(attr)
			switch {
			case strings.HasPrefix(attr, "domain="):
				cookie.Domain = strings.TrimPrefix(attr, "domain=")
			case strings.HasPrefix(attr, "expires="):
				expiryTime, err := time.Parse(time.RFC1123, strings.TrimPrefix(attr, "expires="))
				if err == nil {
					cookie.Expires = expiryTime
				}
			case strings.HasPrefix(attr, "path="):
				cookie.Path = strings.TrimPrefix(attr, "path=")
			case strings.HasPrefix(attr, "secure"):
				cookie.Secure = true
			case strings.HasPrefix(attr, "HttpOnly"):
				cookie.HTTPOnly = true
			default:
				// The name and value should be at the beginning
				if cookie.Name == "" {
					parts := strings.SplitN(attr, "=", 2)
					if len(parts) == 2 {
						cookie.Name = parts[0]
						cookie.Value = parts[1]
					}
				}
			}
		}
		cookies = append(cookies, cookie)
	}
	return cookies, nil
}
