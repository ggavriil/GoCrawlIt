package main

import (
	"log"
	"net/http"
	"golang.org/x/net/html"
	"strings"
)
	

func getAttr(token html.Token, attr string) string {
	for _, a := range token.Attr {
		if a.Key == attr {
			return a.Val
		}
	}
	return ""
}

func processPage(url string, startUrl string) (map[string]bool, map[string]bool, map[string]bool) {
	urls := make(map[string]bool)
	statics := make(map[string]bool)
	imgs := make(map[string]bool)
	resp, err := http.Get(url)

	if err != nil {
        log.Fatal("Error: " + err.Error())
	}

	body := resp.Body
	tokenizer := html.NewTokenizer(body)

	for nextToken := tokenizer.Next(); nextToken != html.ErrorToken; nextToken = tokenizer.Next() {
		if nextToken == html.StartTagToken {
			token := tokenizer.Token()
			switch token.Data {
			case "a":
				nextUrl := getAttr(token, "href")
				if strings.Index(nextUrl, startUrl) == 0 {
					urls[nextUrl] = true
				}
			case "img":
				nextImg := getAttr(token, "src")
				imgs[nextImg] = true
			case "link":
				nextLink := getAttr(token, "href")
				statics[nextLink] = true
			case "script":
				nextScript := getAttr(token, "src")
				statics[nextScript] = true
			}
		}
	}

	defer body.Close()

	return urls, imgs, statics
}
