package serp

import "testing"

func TestSearchBing(t *testing.T) {
	bing, err := SearchBing(nil, "憨憨博客", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36", 10, "")
	if err != nil {
		return
	}
	for _, result := range bing {
		println(result.Description)
	}
}
