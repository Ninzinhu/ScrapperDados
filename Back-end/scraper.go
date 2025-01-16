package main

import (
    "encoding/json"
    "net/http"
    "github.com/gocolly/colly/v2"
)

type Product struct {
    Name  string `json:"name"`
    Price string `json:"price"`
    URL   string `json:"url"`
}

var products []Product

func scrapeProducts() {
    c := colly.NewCollector()

    c.OnHTML("div.product-tile__content-wrapper", func(e *colly.HTMLElement) {
        product := Product{
            Name:  e.ChildText("a.product-tile__name__link"),
            Price: e.ChildText("em.value__price"),
            URL:   e.ChildAttr("a.product-tile__name__link", "href"),
        }
        products = append(products, product)
    })

    // Visita a página de produtos
    err := c.Visit("https://www.jackjones.com/nl/en/jj/shoes/")
    if err != nil {
        panic(err)
    }
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(products)
}

func main() {
    scrapeProducts()
    http.HandleFunc("/api/products", productsHandler)
    http.ListenAndServe(":8080", nil)
}