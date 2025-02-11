package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
)

const (
	displayContentNum = 5
	baseReadme        = `# Welcome to my playground 🏖

[![Blog Badge](https://img.shields.io/badge/-Blog-blue?style=flat&logo=hugo&logoColor=white)](https://tech.yyh-gl.dev/)
[![X Badge](https://img.shields.io/badge/-@yyh__gl-black?logo=x)](https://twitter.com/yyh_gl)
[![Bluesky Badge](https://img.shields.io/badge/-Bluesky｜@yyh__gl-1e90ff?style=flat)](https://bsky.app/profile/yyh-gl.bsky.social)
[![Speaker Deck Badge](https://img.shields.io/badge/-Speaker_Deck-009287?style=flat&logo=speaker-deck&logoColor=white)](https://speakerdeck.com/yyh_gl)
[![Crowdin Badge](https://img.shields.io/badge/-Crowdin-f2f2f2?style=flat&logo=crowdin&logoColor=black)](https://crowdin.com/profile/yyh-gl)

## Blog - Recent posts 📝

`
)

var (
	expTitle = regexp.MustCompile("<item><title>.*</title>")
	expLink  = regexp.MustCompile("<link>https://tech.yyh-gl.dev/blog.*</link>")
)

func main() {
	response, err := http.Get("https://tech.yyh-gl.dev/index.xml")
	if err != nil {
		panic(err)
	}
	defer func() { _ = response.Body.Close() }()
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	feed := string(bytes)

	titles := expTitle.FindAllString(feed, displayContentNum)
	links := expLink.FindAllString(feed, displayContentNum)

	readmeStr := baseReadme
	for i := 0; i < displayContentNum; i++ {
		t := titles[i]
		t = t[13 : len(t)-8]
		l := links[i]
		l = l[6 : len(l)-7]
		readmeStr += fmt.Sprintf("- [%s](%s)\n", t, l)
	}

	readmeFile, err := os.Create("README.md")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = readmeFile.Close() }()

	data := []byte(readmeStr)
	if _, err = readmeFile.Write(data); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
