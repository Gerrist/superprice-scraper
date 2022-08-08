package store_handling

type Store struct {
	Id   string
	Name string
}

type StoreConfig struct {
	ProductIndexType
	ProductIndexUrl string
	Store
	URLBase string
}

type StoreConfigList []StoreConfig

type ProductIndexType int

const (
	Sitemap        ProductIndexType = 1
	SitemapGZipped ProductIndexType = 2
)

type ProductInfo struct {
	Id    string  `json:"id" csv:"id"`
	Store string  `json:"store" csv:"store"`
	Name  string  `json:"name" csv:"name"`
	Price float64 `json:"price" csv:"price"`
	URL   string  `json:"url" csv:"url"`
	Brand string  `json:"brand" csv:"brand"`
}

type ProductURL struct {
	URL   string
	Store Store
}
