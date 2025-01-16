package main

import (
    "encoding/json"
    "fmt"
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

    // Define o que fazer quando um elemento HTML correspondente é encontrado
    c.OnHTML("div.product-tile__content-wrapper", func(e *colly.HTMLElement) {
        name := e.ChildText("a.product-tile__name__link")
        price := e.ChildText("em.value__price")
        url := e.ChildAttr("a.product-tile__name__link", "href")

        // Verifica se os dados foram extraídos corretamente
        if name != "" && price != "" && url != "" {
            product := Product{
                Name:  name,
                Price: price,
                URL:   url,
            }
            products = append(products, product)
        } else {
            fmt.Println("Produto não encontrado ou incompleto.")
        }
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
    fmt.Println("Servidor rodando na porta 8080...")
    http.ListenAndServe(":8080", nil)
}