package main

import (
	"log"
	"flag"
	"fmt"
	"os"
	"net/http"
)

func main() {
    flag.Usage = func() {
        fmt.Printf("Usage: ./gocrawlit [options] <website url>\n\nOptions:\n")
        flag.PrintDefaults()
    }
    reqLimit := flag.Int("l", 15, "Parallel Routines Limit")
    noLimit := flag.Bool("nl", false, "Disable Routines Limit")
    flag.Parse()
    if(len(flag.Args()) < 1) {
        flag.Usage()
        os.Exit(1)
    }
    startUrl := flag.Args()[0]
    _, err := http.Get(startUrl)
    if err != nil {
        log.Fatal("Invalid URL specified or URL is unreachable.")
    }

	output := make(chan []string, 100)
	go func() {
		for strs := range output {
			for _, str := range strs {
				fmt.Print(str)
			}
		}
	}()

    startCrawling(output, *reqLimit, *noLimit, startUrl)
    close(output)
}
