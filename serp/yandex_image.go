package serp

import (
  "encoding/json"
  "fmt"
  "github.com/Danny-Dasilva/CycleTLS/cycletls"
  "github.com/PuerkitoBio/goquery"
  "log"
  "net/url"
  "sort"

  //"strconv"
  "strings"
)

func SearchYandexImage(myJa3 string, query, ua string, limit int, proxyAddr string, verbose bool) ([]ImageResult, error) {
  client := cycletls.Init()
  url := "https://yandex.ru/images/search?lr=213&text=" + url.QueryEscape(query)
  response, err := client.Do(url, cycletls.Options{
    Body:      "",
    Ja3:       myJa3,
    UserAgent: ua,
    Headers: map[string]string{
      "accept-language":           "zh-CN,zh;q=0.9,en;q=0.8,en-XA;q=0.7,ja-JP;q=0.6,ja;q=0.5,zh-TW;q=0.4",
      "cache-control":             "no-cache",
      "connection":                "keep-alive",
      "dnt":                       "1",
      "pragma":                    "no-cache",
      "sec-ch-ua":                 "\"Not(A:Brand\";v=\"99\", \"Google Chrome\";v=\"133\", \"Chromium\";v=\"133\"",
      "sec-ch-ua-mobile":          "?0",
      "sec-ch-ua-platform":        "\"macOS\"",
      "sec-fetch-dest":            "document",
      "sec-fetch-mode":            "navigate",
      "sec-fetch-site":            "same-origin",
      "sec-fetch-user":            "?1",
      "upgrade-insecure-requests": "1",
    },
  }, "GET")
  if err != nil {
    log.Print("Request Failed: " + err.Error())
  }
  doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.Body))
  if err != nil {
    log.Fatal(err)
  }
  var results []ImageResult
  filteredRank := 1
  rank := 1
  //println(response.Body)
  // 遍历每个图片项
  //doc.Find("div.JustifierRowLayout-Item").Each(func(i int, sel *goquery.Selection) {
  //	// 解析 mUrl
  //	mStr, _ := sel.Find("div.SerpItem > a").Attr("href")
  //	mStr = strings.Replace(mStr, "/images/search?pos=", "", 1)
  //	splitRes := strings.SplitN(mStr, "&", 2)
  //	mUrl := splitRes[1]
  //
  //	// 解析图片宽度和高度
  //	width, _ := sel.Find("div.SerpItem > img").Attr("width")
  //	widthVal, _ := strconv.Atoi(width)
  //	height, _ := sel.Find("div.SerpItem > img").Attr("height")
  //	heightVal, _ := strconv.Atoi(height)
  //
  //	// 解析图片描述
  //	alt, _ := sel.Find("div.SerpItem > img").Attr("alt")
  //	desc := alt
  //
  //	// 解析图片源地址
  //	src, _ := sel.Find("div.imgpt > a.iusc > div > img").Attr("src")
  //
  //	// 解析标签
  //	label, _ := sel.Find("div.infopt > a.inflnk").Attr("aria-label")
  //
  //	// 如果 mUrl 和 src 有效，则添加到结果中
  //	if mUrl != "" && mUrl != "#" && src != "" {
  //		result := ImageResult{
  //			MUrl:  mUrl,
  //			Label: label,
  //			Size: Size{
  //				Width:  widthVal,
  //				Height: heightVal,
  //			},
  //		}
  //		if verbose {
  //			result.Desc = desc
  //			result.Src = src
  //		}
  //		results = append(results, result)
  //		filteredRank += 1
  //	}
  //	rank += 1
  //})

  doc.Find(".page-layout__content-wrapper > div.Root").Each(func(i int, sel *goquery.Selection) {
    data, _ := sel.Attr("data-state")

    var result map[string]interface{}
    err := json.Unmarshal([]byte(data), &result)
    if err != nil {
      fmt.Println("Error:", err)
      return
    }
    //println(result)
    initialState, ok := result["initialState"].(map[string]interface{})
    if !ok {
      fmt.Println("Error: initialState is not a map")
      return
    }

    serpList, ok := initialState["serpList"].(map[string]interface{})
    if !ok {
      fmt.Println("Error: serpList is not a map")
      return
    }

    items, ok := serpList["items"].(map[string]interface{})
    if !ok {
      fmt.Println("Error: items is not a map")
      return
    }

    entities, ok := items["entities"].(map[string]interface{})
    if !ok {
      fmt.Println("Error: entities is not a map")
      return
    }

    var entityList []Entity
    for key, value := range entities {
      entityList = append(entityList, Entity{Key: key, Value: value.(map[string]interface{})})
    }
    sort.Slice(entityList, func(i, j int) bool {
      posI := entityList[i].Value["pos"].(float64)
      posJ := entityList[j].Value["pos"].(float64)
      return posI < posJ // 从小到大排序
    })
    // 输出 entities
    for _, entity := range entityList {
      val := entity.Value
      if !ok {
        fmt.Println("Error: entity is not a map")
        continue
      }
      origUrl := val["origUrl"].(string)
      width := val["width"].(float64)
      height := val["height"].(float64)
      snippet := val["snippet"].(map[string]interface{})
      title := val["alt"].(string)
      if origUrl != "#" && origUrl != "" {
        result := ImageResult{
          MUrl:  origUrl,
          Label: title,
          Size: Size{
            Width:  int(width),
            Height: int(height),
          },
        }
        if verbose {
          //result.Desc = desc
          result.Src = snippet["url"].(string)
        }
        results = append(results, result)
        filteredRank += 1
      }
      rank += 1
    }
  })

  // 限制结果数量
  if limit != 0 && len(results) > limit {
    return results[:limit], nil
  }

  return results, nil
}

type Entity struct {
  Key   string
  Value map[string]interface{}
}
