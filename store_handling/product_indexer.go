package store_handling

import (
	"fmt"
	"github.com/beevik/etree"
	"log"
	"net/http"
)

func ProductIndexer(storeConfigList StoreConfigList) (productUrls []ProductURL) {
	for _, storeConfig := range storeConfigList {
		log.Printf("Indexing " + storeConfig.Store.Name)

		switch storeConfig.ProductIndexType {
		case Sitemap:
			{
				sitemapUrls, err := IndexSitemap(storeConfig)
				if err != nil {
					log.Fatalln(err)
				}
				productUrls = append(productUrls, sitemapUrls...)
				break
			}
		default:
			{
				log.Println("Unknown index type", storeConfig.ProductIndexType)
			}
		}
	}

	return productUrls
}

func IndexSitemap(storeConfig StoreConfig) (
	productUrls []ProductURL,
	error error,
) {
	response, err := http.Get(storeConfig.ProductIndexUrl)
	if err != nil {
		log.Fatalln(err)
	}
	if response.StatusCode != 200 {
		log.Fatalln("Response status code is not 200")
	}

	doc := etree.NewDocument()
	if _, err := doc.ReadFrom(response.Body); err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parser error: %+v", r)
		}
	}()

	urlset := doc.SelectElement("urlset")
	urlsetUrls := urlset.SelectElements("url")

	for _, url := range urlsetUrls {
		productUrls = append(productUrls, ProductURL{
			URL:   url.SelectElement("loc").Text(),
			Store: storeConfig.Store,
		})
	}

	return productUrls, nil
}
