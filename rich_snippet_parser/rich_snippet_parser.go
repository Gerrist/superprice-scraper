package rich_snippet_parser

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/gosimple/slug"
	"io"
	"log"
	"strconv"
	"superprice-scraper/store_handling"
)

type BasicSnippet struct {
	Context string `json:"@context"`
	Type    string `json:"@type"`
}

type ProductSnippet struct {
	Context     string `json:"@context"`
	Type        string `json:"@type"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	Description string `json:"description"`
	SKU         string `json:"sku"`
	Weight      string `json:"weight"`
	GTIN13      string `json:"gtin13"`
	Image       string `json:"image"`
	Brand       struct {
		Type string `json:"@type"`
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	Offers struct {
		Type          string `json:"@type"`
		Url           string `json:"url"`
		PriceCurrency string `json:"priceCurrency"`
		Price         string `json:"price"`
		ItemCondition string `json:"itemCondition"`
		Availability  string `json:"availability"`
		Seller        struct {
			Type string `json:"@type"`
			Name string `json:"name"`
		}
		PotentialAction struct {
			Type   string `json:"@type"`
			Target string `json:"target"`
		}
	}
}

func (s ProductSnippet) AsProductInfo() store_handling.ProductInfo {
	price, err := strconv.ParseFloat(s.Offers.Price, 8)
	if err != nil {
		price = -1.0
	}
	return store_handling.ProductInfo{
		Store: slug.Make(s.Offers.Seller.Name),
		Name:  s.Name,
		Price: price,
		URL:   s.Url,
		Brand: slug.Make(s.Brand.Name),
		Id:    slug.Make(s.Name + "-" + s.Offers.Seller.Name),
	}
}

func GetSnippetType(snippetString string) string {
	basicSnippet := BasicSnippet{}
	err := json.Unmarshal([]byte(snippetString), &basicSnippet)
	if err != nil {
		return "unknown"
	}
	return basicSnippet.Type
}

func ParseProductSnippet(snippetString string) (snippet ProductSnippet, error error) {
	productSnippet := ProductSnippet{}
	err := json.Unmarshal([]byte(snippetString), &productSnippet)
	if err != nil {
		return ProductSnippet{}, err
	}
	return productSnippet, nil
}

func FindProductSnippetsInReader(reader io.Reader) (snippets []ProductSnippet, error error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}

	productSnippetSelection := doc.Find("script[type=\"application/ld+json\"]")
	productSnippetSelection.Each(func(i int, selection *goquery.Selection) {
		snippetType := GetSnippetType(selection.Text())
		if snippetType == "Product" {
			productSnippet, err := ParseProductSnippet(selection.Text())
			if err != nil {
				log.Panicln(err)
			}
			snippets = append(snippets, productSnippet)
		}
	})
	return snippets, nil
}
