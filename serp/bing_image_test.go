package serp

import "testing"

func TestSearchBingImage(t *testing.T) {
	bing, err := SearchBingImage(nil, "猫", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36", 10, "")
	if err != nil {
		return
	}
	for _, result := range bing {
		println(result.MUrl)
	}
}
