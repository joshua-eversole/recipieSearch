package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gofiber/fiber/v2"
)

type NutritionFacts struct {
	Calories string `json:"calories"`
	Fat      string `json:"fat"`
	Carbs    string `json:"carbs"`
	Protein  string `json:"protein"`
}

type Recipe struct {
	Name           string         `json:"name"`
	PrepTime       string         `json:"prepTime"`
	CookTime       string         `json:"cookTime"`
	TotalTime      string         `json:"totalTime"`
	Ingredients    []string       `json:"ingredients"`
	Steps          []string       `json:"steps"`
	NutritionFacts NutritionFacts `json:"nutritionFacts"` // NitritionFacts should be its own sub-struct

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

	// Scrape the list of ingredients
	c.OnHTML("li.mm-recipes-structured-ingredients__list-item ", func(h *colly.HTMLElement) {
		ingredient := strings.TrimSpace(h.Text)
		recipe.Ingredients = append(recipe.Ingredients, ingredient)
	})

	// Scrape the list of steps
	c.OnHTML("li.comp.mntl-sc-block.mntl-sc-block-startgroup.mntl-sc-block-group--LI", func(h *colly.HTMLElement) {
		// Extract the text from the <p> tag inside the <li>
		step := h.ChildText("p.comp.mntl-sc-block.mntl-sc-block-html")
		step = strings.TrimSpace(step) // Trim any extra whitespace
		if step != "" {                // Only append if the step is not empty
			recipe.Steps = append(recipe.Steps, step)
		}
	})

	// Scrape the nutrition facts
	c.OnHTML("tbody.mm-recipes-nutrition-facts-summary__table-body", func(h *colly.HTMLElement) {
		h.ForEach("tr.mm-recipes-nutrition-facts-summary__table-row", func(_ int, row *colly.HTMLElement) {
			label := row.ChildText("td.mm-recipes-nutrition-facts-summary__table-cell.text-body-100")
			value := row.ChildText("td.mm-recipes-nutrition-facts-summary__table-cell.text-body-100-prominent")

			switch strings.TrimSpace(label) {
			case "Calories":
				recipe.NutritionFacts.Calories = value
			case "Fat":
				recipe.NutritionFacts.Fat = value
			case "Carbs":
				recipe.NutritionFacts.Carbs = value
			case "Protein":
				recipe.NutritionFacts.Protein = value
			}
		})
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

func printScrapedData(recipe Recipe) {
	fmt.Printf("Dish Name: %s\nPrep Time: %s\nCook Time: %s\nTotal Time: %s\n", recipe.Name, recipe.PrepTime, recipe.CookTime, recipe.TotalTime)
	// Print the ingredients
	fmt.Println("Ingredients:")
	for _, ingredient := range recipe.Ingredients {
		fmt.Println("-", ingredient)
	}
	// Print the steps
	fmt.Println("Steps:")
	for _, step := range recipe.Steps {
		fmt.Println("-", step)
	}

	// Print the nutrition facts
	fmt.Println("Nutrition Facts:")
	fmt.Println("- Calories:", recipe.NutritionFacts.Calories)
	fmt.Println("- Fat:", recipe.NutritionFacts.Fat)
	fmt.Println("- Carbs:", recipe.NutritionFacts.Carbs)
	fmt.Println("- Protein:", recipe.NutritionFacts.Protein)
}

func main() {
	// Scrape Allrecipes
	url := "https://www.allrecipes.com/recipe/223042/chicken-parmesan/"
	recipe, err := scrapeAllrecipes(url)
	if err != nil {
		log.Fatal("Error scraping Allrecipes:", err)
	}

	//Trying
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"msg": "Hello, World!"})
	})

	// Print the scraped data
	printScrapedData(recipe)
	// You can add more scraping functions for other sites here
}
