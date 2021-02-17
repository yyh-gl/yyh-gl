package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

const (
	displayContentNum = 5
	baseReadme        = `# Welcome to my playground 🏖

[![Blog Badge](https://img.shields.io/badge/-Blog-blue?style=flat&logo=hugo&logoColor=white)](https://yyh-gl.github.io/tech-blog/)
[![Speaker Deck Badge](https://img.shields.io/badge/-Speaker_Deck-009287?style=flat&logo=speaker-deck&logoColor=white)](https://speakerdeck.com/yyh_gl)
[![Twitter Badge](https://img.shields.io/badge/-@yyh__gl-1ca0f1?style=flat&logo=twitter&logoColor=white)](https://twitter.com/yyh_gl)

## Recent posts - Blog 📝

`
)

var (
	expTitle = regexp.MustCompile("<item><title>.*</title>")
	expLink  = regexp.MustCompile("<link>https://yyh-gl.github.io/tech-blog/blog.*</link>")
)

func main() {
	response, err := http.Get("https://yyh-gl.github.io/tech-blog/index.xml")
	if err != nil {
		panic(err)
	}
	defer func() { _ = response.Body.Close() }()
	bytes, err := ioutil.ReadAll(response.Body)
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
