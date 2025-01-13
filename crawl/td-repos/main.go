package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

var url = "https://github.com/GoesToEleven?tab=repositories&q=&type=&language=&sort=stargazers"

type Repo struct {
    Name string
    Language string
    Stars int
    Forks int
    LastUpdated string
}

func main() {
    c := colly.NewCollector()
	
    repos := []Repo{}
    c.OnHTML("li > div:nth-child(1)", func(e *colly.HTMLElement) {
	repo := Repo{}
	repo.Name = e.ChildText("div:nth-child(1) > h3 > a")
	repo.Language = e.ChildText("div:nth-child(3) > span > span:nth-child(2)")


	htmlStars := e.ChildText("div:nth-child(3) > a:nth-child(2)")
	if htmlStars != "" {
	    textStars := strings.Replace(htmlStars, ",", "", -1)
	    stars, err := strconv.Atoi(textStars)
	    if err != nil {
		panic("couldn't parse stars")
	    }
	    repo.Stars = stars
	}


	htmlForks := e.ChildText("div:nth-child(3) > a:nth-child(3)")
	if htmlForks != "" {
	    textForks := strings.Replace(htmlForks, ",", "", -1)
	    forks, err := strconv.Atoi(textForks)
	    if err != nil {
		panic("couldn't parse forks")
	    }
	    repo.Forks = forks
	}
	
	repo.LastUpdated = e.ChildText("div:nth-child(3) > relative-time")

	repos = append(repos, repo)
    }) 
    
    c.Visit(url)

    printRepos(repos)

    c.Wait()
}

func printRepos(repos []Repo) {
    for _, repo := range repos {
	fmt.Println("Name:", repo.Name)
	fmt.Println("Language:", repo.Language)
	fmt.Println("Stars:", repo.Stars)
	fmt.Println("Forks:", repo.Forks)
	fmt.Println("Last updated:", repo.LastUpdated, "\n")
    }
}

// Within previously selected element
// Repo name: wb-break-all > a
// Lang: span > span:ntn-child(2)
// Stars: "f6 color-fg-muted mt-2" > a:ntn-child(2) 
// Forks: "f6 color-fg-muted mt-2" > a:ntn-child(3)
// Time: relative-time attr=datetime
