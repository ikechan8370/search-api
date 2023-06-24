package main

import "testing"

func TestSearchBaidu(t *testing.T) {
	bing, err := SearchBaidu(nil, "憨憨博客", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36", 10)
	if err != nil {
		panic(err)
	}
	for _, result := range bing {
		println(result.Description)
	}
}
