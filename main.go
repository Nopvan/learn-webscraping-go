package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gocolly/colly"
)

// buat struct untuk menampung elemen yang ingin diambil
type Item struct {
	Link    string `json:"link"`
	Name    string `json:"name"`
	Price   string `json:"price"`
	Instock string `json:"instock"`
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v", name, time.Since(start))
	}
}

func main() {
	defer timer("time")()

	c := colly.NewCollector()

	items := []Item{}

	// baris untuk mengambil dan menentukan elemen mana yang akan dipakai
	c.OnHTML("li.next a", func(h *colly.HTMLElement) {

		// baris yang akan mengatur apa yang akan dilakukan jika menemukan elemen css yang sudah diatur sebelumnya
		c.Visit(h.Request.AbsoluteURL(h.Attr("href")))
	})

	// baris untuk mengambil dan menentukan elemen mana yang akan dipakai
	c.OnHTML("article.product_pod", func(h *colly.HTMLElement) {
		i := Item{
			Link:    h.ChildAttr("a", "href"),
			Name:    h.ChildAttr("h3 a", "title"),
			Price:   h.ChildText("p.price_color"),
			Instock: h.ChildText("p.instock"),
		}
		items = append(items, i)
	})

	// baris untuk dapat mengambil semua data yang berada di page lain dengan cara request ke page nya
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL)
	})

	//set url website
	c.Visit("https://books.toscrape.com/catalogue/page-1.html")

	//make to json file
	data, err := json.Marshal(items)
	if err != nil {
		fmt.Println(err.Error())
	}
	os.WriteFile("books.json", data, 0644)
}
