package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strings"
	"sync"
)

func getAttr (token html.Token, attr string) string{
	for _, a := range token.Attr {
		if a.Key == attr {
			return a.Val
		}
	}
	return ""
}

func crawl(thisUrl string, visitedUrls *sync.Map, startUrl string, res chan string, depth int) {
	defer close(res)
	visitedUrls.Store(thisUrl, true)
	curUrls := make(map[string]bool)//TODO
	curScriptsLinks := make(map[string]bool)
	curImgs := make(map[string]bool)
	res <- fmt.Sprintf("Page URL: %s \n", thisUrl)
	res <- fmt.Sprintf("Depth: %d\n\n", depth)
	resp, err := http.Get(thisUrl)

	if err != nil {
		fmt.Println("Error:" + err.Error())
		return
	}

	body := resp.Body;
	tokenizer := html.NewTokenizer(body)
	for {
		nextToken := tokenizer.Next();

		if nextToken == html.ErrorToken {
			break
		}

		if nextToken == html.StartTagToken {
			token := tokenizer.Token()
			switch token.Data {
			case "a":
				nextUrl := getAttr(token, "href")
				if strings.Index(nextUrl, startUrl) == 0 {
					curUrls[nextUrl] = true
				}
			case "img":
				nextImg := getAttr(token, "src")
				curImgs[nextImg] = true
			case "link":
				nextLink := getAttr(token, "href")
				curScriptsLinks[nextLink] = true
			case "script":
				nextScript := getAttr(token, "src")
				curScriptsLinks[nextScript] = true
			}
		}
	}
	body.Close()

	res <- fmt.Sprintf("Links to:\n")
	for url, _ := range curUrls {
		res <- fmt.Sprintf("%s\n", url)
	}
	res <- fmt.Sprintf("\nImages in page:\n")
	for img, _ := range curImgs {
        res <- fmt.Sprintf("%s\n", img)
	}
	res <- fmt.Sprintf("\nScripts and other static assets in page:\n")
	for scrl, _ := range curScriptsLinks {
        res <- fmt.Sprintf("%s\n", scrl)
	}
	res <- fmt.Sprintf("------------------------------------------------------------\n")

	chRes := make([]chan string, len(curUrls))
	j := 0
	for nextUrl, _ := range curUrls {
		chRes[j] = make(chan string)
        _, loaded := visitedUrls.Load(nextUrl)
        if !loaded {
			go crawl(nextUrl, visitedUrls, startUrl, chRes[j], depth + 1)
		} else {
			close(chRes[j])
		}
		j++
	}

	for i := range chRes {
		for str := range chRes[i] {
			res <- str
		}
	}

	return
}

func main() {
	startUrl := os.Args[1];
	res := make(chan string)
	go crawl(startUrl, new(sync.Map), startUrl, res, 0) //TODO &
	for str := range res {
		fmt.Print(str)
	}
}


