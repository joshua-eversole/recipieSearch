package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type Recipe struct {
	Name      string `json:"name"`
	PrepTime  string `json:"prepTime"`
	CookTime  string `json:"cookTime"`
	TotalTime string `json:"totalTime"`
}

func scrapeAllrecipes(url string) (Recipe, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("www.allrecipes.com"),
	)

	var recipe Recipe

	// Scrape the dish name
	c.OnHTML("h1.article-heading", func(h *colly.HTMLElement) {
		recipe.Name = h.Text
	})

	// Scrape the times
	c.OnHTML("div.mm-recipes-details__item", func(h *colly.HTMLElement) {
		label := h.ChildText("div.mm-recipes-details__label")
		value := h.ChildText("div.mm-recipes-details__value")

		switch strings.TrimSpace(label) {
		case "Prep Time:":
			recipe.PrepTime = value
		case "Cook Time:":
			recipe.CookTime = value
		case "Total Time:":
			recipe.TotalTime = value
		}
	})

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Visit the URL
	err := c.Visit(url)
	if err != nil {
		return Recipe{}, err
	}

	return recipe, nil
}

func main() {
	// Scrape Allrecipes
	url := "https://www.allrecipes.com/recipe/223042/chicken-parmesan/"
	recipe, err := scrapeAllrecipes(url)
	if err != nil {
		log.Fatal("Error scraping Allrecipes:", err)
	}

	// Print the scraped data
	fmt.Printf("Dish Name: %s\nPrep Time: %s\nCook Time: %s\nTotal Time: %s\n", recipe.Name, recipe.PrepTime, recipe.CookTime, recipe.TotalTime)

	// You can add more scraping functions for other sites here
}
