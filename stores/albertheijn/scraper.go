package albertheijn

import (
	"fmt"
	"log"
	"net/http"
	"superprice-scraper/rich_snippet_parser"
	"superprice-scraper/store_handling"
)

func GetProductInfo(productUrl string) (productInfo store_handling.ProductInfo, err error) {
	response, err := http.Get(productUrl)
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		log.Fatalln("Response status code is not 200")
	}

	richProductSnippets, err := rich_snippet_parser.FindProductSnippetsInReader(response.Body)
	if err != nil {
		return store_handling.ProductInfo{}, err
	}
	if len(richProductSnippets) == 0 {
		return store_handling.ProductInfo{}, fmt.Errorf("No product snippets found")
	}
	return richProductSnippets[0].AsProductInfo(), nil

	//if productSnippetSelection.Length() == 1 {
	//
	//} else {
	//	log.Println("No product snippet found")
	//}
	//doc := etree.NewDocument()
	//
	//buf := new(bytes.Buffer)
	//_, err = buf.ReadFrom(response.Body)
	//if err != nil {
	//	return store_handling.ProductInfo{}, err
	//}
	////escapedString := strings.ReplaceAll(buf.String(), "&", "!!!AMP!!!")
	////
	////log.Println(escapedString)
	//
	////if _, err := doc.ReadFrom(strings.NewReader(escapedString)); err != nil {
	//if _, err := doc.ReadFrom(response.Body); err != nil {
	//	return store_handling.ProductInfo{}, err
	//}
	//
	//defer func() {
	//	if r := recover(); r != nil {
	//		err = fmt.Errorf("parser error: %+v", r)
	//	}
	//}()
	//
	//snippetElement := FindSnippetElement(doc.SelectElement("html").SelectElement("head"))
	//if snippetElement == nil {
	//	log.Fatalln("Could not find snippet element on", productUrl)
	//}
	//log.Println(snippetElement.Text())

}
