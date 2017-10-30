package main

import (
	"sync/atomic"
	"fmt"
)

func crawl(thisUrl string, urls chan string, startUrl string) []string {
	curUrls, curImgs, curStatics := processPage(thisUrl, startUrl)

	retText := make([]string, 0, 20)
	retText = append(retText, fmt.Sprintf("Page URL: %s \n", thisUrl))
	retText = append(retText, fmt.Sprintf("Links to:\n"))
	for url, _ := range curUrls {
		retText = append(retText, fmt.Sprintf("%s\n", url))
		urls <- url
	}
	retText = append(retText, fmt.Sprintf("\nImages in page:\n"))
	for img, _ := range curImgs {
		retText = append(retText, fmt.Sprintf("%s\n", img))
	}
	retText = append(retText, fmt.Sprintf("\nStatic assets in page:\n"))
	for sta, _ := range curStatics {
		retText = append(retText, fmt.Sprintf("%s\n", sta))
	}
	retText = append(retText,
		fmt.Sprintf("------------------------------------------------------------\n"))
	return retText
}


func startCrawling(output chan []string, reqLimit int, noLimit bool, startUrl string) *map[string]bool {
    urls := make(chan string, 10)
    visited := make(map[string]bool)
    var launched uint64 = 0
	var done uint64 = 0
    urls <- startUrl
    sem := make(chan struct{}, reqLimit)

ForLoop:
	for {
		select {
		case url := <-urls:
			_, ok := visited[url]
			if !ok {
				visited[url] = true
				atomic.AddUint64(&launched, 1)
				go func(url string) {
                    if(!noLimit) {
                        sem <- struct{}{}
                    }
					output <- crawl(url, urls, startUrl)
					atomic.AddUint64(&done, 1)
                    if(!noLimit) {
                        <-sem
                    }
				}(url)
			}
		default:
			if (atomic.LoadUint64(&launched) == atomic.LoadUint64(&done)) {
				close(urls)
				break ForLoop
			}
		}
	}

    return &visited
}


