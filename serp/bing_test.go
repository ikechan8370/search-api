package serp

import "testing"

func TestSearchBing(t *testing.T) {
	bing, err := SearchBing("771,4865-4867-4866-49195-49199-52393-52392-49196-49200-49162-49161-49171-49172-51-57-47-53-10,0-23-65281-10-11-35-16-5-51-43-13-45-28-21,29-23-24-25-256-257,0", nil, "缇宝天赋攻略 米游社", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36", 10, "", nil)
	if err != nil {
		panic(err.Error())
		return
	}
	for _, result := range bing {
		println(result.Description)
	}
}
