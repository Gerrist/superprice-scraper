package main

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"log"
	"math"
	"os"
	"strconv"
	"superprice-scraper/store_handling"
	"superprice-scraper/stores/albertheijn"
	"time"
)

func main() {
	log.Println("SuperPrice Scraper")

	storeConfigs := store_handling.StoreConfigList{
		store_handling.StoreConfig{
			ProductIndexType: store_handling.Sitemap,
			ProductIndexUrl:  "https://ah.nl/sitemaps/entities/products/detail.xml",
			URLBase:          "https://ah.nl",
			Store: store_handling.Store{
				Id:   "albertheijn",
				Name: "Albert Heijn",
			},
		},
	}

	productUrls := store_handling.ProductIndexer(storeConfigs)
	products := make([]store_handling.ProductInfo, 0)

	for i, url := range productUrls {
		switch url.Store.Id {
		case "albertheijn":
			{
				productInfo, err := albertheijn.GetProductInfo(url.URL)
				if err != nil {
					log.Fatalln(err)
				}
				log.Println(productInfo.Store, productInfo.Name, productInfo.Price, productInfo.Id)
				log.Println(i, "of", len(productUrls))
				products = append(products, productInfo)
				if math.Mod(float64(i), 100) == 0 && i > 0 {
					{
						log.Println("Saving data and waiting 1 minute")
						productsFile, err := os.OpenFile(fmt.Sprintf("results/products-%s.csv", strconv.Itoa(int(time.Now().Unix()))), os.O_RDWR|os.O_CREATE, os.ModePerm)
						if err != nil {
							log.Fatalln(err)
						}
						err = gocsv.MarshalFile(&products, productsFile)
						if err != nil {
							log.Fatalln(err)
						}
						time.Sleep(1 * time.Minute)

					}
					break
				}
			}
		}
	}
}
