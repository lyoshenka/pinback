package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

func skip(p Post) bool {
	hostsToSkip := []string{
		"github.com",
		"soundcloud.com", // handle via scdl
		//"youtube.com",    // handle via youtube-dl or maybe ytmp3
	}

	for _, h := range hostsToSkip {
		if p.Domain() == h {
			return true
		}
	}
	return false
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s PINBOARD_API_TOKEN\n", path.Base(os.Args[0]))
		os.Exit(1)
	}
	c := NewClient(os.Args[1])
	recent, err := c.Since(time.Now().AddDate(0, -1, 0))
	//recent, err := c.Recent()
	if err != nil {
		panic(err)
	}

	var urls []string
	for _, p := range recent {
		if skip(p) {
			continue
		}
		urls = append(urls, "'"+p.URL+"'")
	}

	fmt.Println(strings.Join(urls, " "))
}
