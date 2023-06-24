package main

import (
	"github.com/gin-gonic/gin"
	googlesearch "github.com/rocketlaunchr/google-search"
	"net/http"
	"strconv"
)

func main() {
	r := gin.Default()
	r.GET("/google", func(c *gin.Context) {
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
		opts := googlesearch.SearchOptions{
			Limit: limitNum,
			//ProxyAddr:    "http://127.0.0.1:7890",
			LanguageCode: lang,
		}
		returnLinks, err := googlesearch.Search(nil, q, opts)
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
