package shopify

type ShippingLine struct {
	Code string `json:"code"`

	Price string `json:"price"`

	Source string `json:"source"`

	Title string `json:"title"`

	TaxLines []interface{} `json:"tax_lines"`
}
