package shopify

type WeightBasedShippingRate struct {
	CountryId int64 `json:"country_id"`

	Id int64 `json:"id"`

	Name string `json:"name"`

	Price string `json:"price"`

	WeightHigh float64 `json:"weight_high"`

	WeightLow float64 `json:"weight_low"`

	Offsets []interface{} `json:"offsets"`
}
