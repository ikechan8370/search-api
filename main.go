package main

import (
	"awesomeProject/serp"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const proxyAddr = ""
const ja3 = ""
const ja3Ua = ""
const geminiBaseUrl = "https://generativelanguage.googleapis.com"

var geminiKeys = []string{}

type Response struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"content"`
}

func main() {
	r := gin.Default()
	r.GET("/google", func(c *gin.Context) {
		if len(geminiKeys) == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "google is not available",
				"code":    500,
			})
			return
		}
		q := c.Query("q")
		if q == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "params q shouldn't be blank",
				"code":    400,
			})
			return
		}
		lang := c.Query("lang")
		if lang == "" {
			lang = "zh-CN"
		}
		limit := c.Query("limit")
		if limit == "" {
			limit = "10"
		}
		verbose := c.Query("verbose")
		if verbose == "" {
			verbose = "false"
		}
		_, err := strconv.Atoi(limit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "params limit should be a number",
				"code":    400,
			})
			return
		}
		url := geminiBaseUrl + "/v1beta/models/gemini-2.0-flash-001:generateContent"
		body := "{\n    \"contents\": [\n        {\n            \"parts\": [\n                {\n                    \"text\": \"search and tell me everything about 'ikechan8370'. Prefer to use LANGUAGE. Keep the original search results in the format of a json array. Each result should have these fields: title, description<URL>. Keep the description with the same language as the original search results. You only need to return at most the first LIMIT results.\"\n                }\n            ]\n        }\n    ],\n    \"tools\": [\n        {\n            \"google_search\": {}\n        }\n    ]\n}"
		body = strings.ReplaceAll(body, "ikechan8370", q)
		body = strings.ReplaceAll(body, "LANGUAGE", lang)
		body = strings.ReplaceAll(body, "LIMIT", limit)
		if verbose == "true" {
			body = strings.ReplaceAll(body, "<URL>", " and url. url should be the original url instead of the one with 'vertexaisearch.cloud.google.com'.")
		} else {
			body = strings.ReplaceAll(body, "<URL>", "")
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rand.Seed(time.Now().UnixNano())

		// 随机选择一个索引
		randomIndex := rand.Intn(len(geminiKeys))

		// 获取随机选中的元素
		randomKey := geminiKeys[randomIndex]
		req.Header.Set("x-goog-api-key", randomKey)
		client := &http.Client{}

		// 执行请求
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		// 打印响应状态码
		//println("Response Status:", resp.Status)

		// 读取并打印响应体
		rspBdy, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		// 打印响应状态码
		//println("Response Status:", resp.Status)
		//println(string(rspBdy))
		// 解析 JSON 响应体
		var response Response
		if err := json.Unmarshal(rspBdy, &response); err != nil {
			log.Fatal(err)
		}

		// 找到需要的部分
		var jsonStr string
		for _, part := range response.Candidates[0].Content.Parts {
			if strings.Contains(part.Text, "```json") {
				// 使用正则表达式去除 ```json 和 ``` 部分
				jsonStr = strings.TrimPrefix(part.Text, "```json")
				jsonStr = strings.TrimPrefix(jsonStr, "```")
				jsonStr = strings.TrimSuffix(jsonStr, "```")
				jsonStr = strings.TrimSpace(jsonStr) // 去除可能的换行或空格
				jsonStr = strings.ReplaceAll(jsonStr, "\\n", "")
				fmt.Println("Extracted JSON:", jsonStr)
				var result []map[string]interface{}
				if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
					log.Fatal(err)
				}
				c.JSON(http.StatusOK, gin.H{
					"message": "ok",
					"code":    200,
					"data":    result,
					"source":  "gemini",
				})
				break
			}
		}
		if jsonStr == "" {
			log.Fatal("No JSON content found")
			c.JSON(http.StatusOK, gin.H{
				"message": "ok",
				"code":    200,
				"data":    nil,
			})
		}
	})
	r.GET("/bing", func(c *gin.Context) {
		q := c.Query("q")
		if q == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "params q shouldn't be blank",
				"code":    400,
			})
			return
		}
		lang := c.Query("lang")
		if lang == "" {
			lang = "zh-CN"
		}
		limit := c.Query("limit")
		if limit == "" {
			limit = "20"
		}
		limitNum, err := strconv.Atoi(limit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "params limit should be a number",
				"code":    400,
			})
			return
		}
		returnLinks, err := serp.SearchBing(nil, q, serp.DefaultUA, limitNum, proxyAddr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"code":    500,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"code":    200,
			"data":    returnLinks,
		})
	})
	r.GET("/baidu", func(c *gin.Context) {
		q := c.Query("q")
		if q == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "params q shouldn't be blank",
				"code":    400,
			})
			return
		}
		lang := c.Query("lang")
		if lang == "" {
			lang = "zh-CN"
		}
		limit := c.Query("limit")
		if limit == "" {
			limit = "20"
		}
		limitNum, err := strconv.Atoi(limit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "params limit should be a number",
				"code":    400,
			})
			return
		}
		returnLinks, err := serp.SearchBaidu(nil, q, serp.DefaultUA, limitNum)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"code":    500,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"code":    200,
			"data":    returnLinks,
		})
	})
	r.GET("/image/bing", func(c *gin.Context) {
		q := c.Query("q")
		if q == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "params q shouldn't be blank",
				"code":    400,
			})
			return
		}
		limit := c.Query("limit")
		if limit == "" {
			limit = "20"
		}
		verbose := false
		verboseS := c.Query("verbose")
		if verboseS == "true" {
			verbose = true
		}
		limitNum, err := strconv.Atoi(limit)
		returnLinks, err := serp.SearchBingImage(nil, q, serp.DefaultUA, limitNum, proxyAddr, verbose)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"code":    500,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"code":    200,
			"data":    returnLinks,
		})
	})
	r.GET("/image/yandex", func(c *gin.Context) {
		q := c.Query("q")
		if q == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "params q shouldn't be blank",
				"code":    400,
			})
			return
		}
		limit := c.Query("limit")
		if limit == "" {
			limit = "20"
		}
		verbose := false
		verboseS := c.Query("verbose")
		if verboseS == "true" {
			verbose = true
		}
		limitNum, err := strconv.Atoi(limit)
		returnLinks, err := serp.SearchYandexImage(ja3, q, ja3Ua, limitNum, proxyAddr, verbose)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"code":    500,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"code":    200,
			"data":    returnLinks,
		})
	})
	r.Run("0.0.0.0:28080")
}
