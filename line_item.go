package shopify

import (
	"time"
)

type LineItem struct {
	AppliedDiscounts   []interface{} `json:"applied_discounts,omitempty"`
	CompareAtPrice     string        `json:"compare_at_price,omitempty"`
	FulfillmentService string        `json:"fulfillment_service,omitempty"`
	GiftCard           bool          `json:"gift_card,omitempty"`
	Grams              int64         `json:"grams,omitempty"`
	LinePrice          time.Time     `json:"line_price,omitempty"`
	Price              string        `json:"price,omitempty"`
	ProductId          int64         `json:"product_id,omitempty"`
	Properties         string        `json:"properties,omitempty"`
	Quantity           int64         `json:"quantity,omitempty"`
	RequiresShipping   bool          `json:"requires_shipping,omitempty"`
	Sku                string        `json:"sku,omitempty"`
	TaxLines           []interface{} `json:"tax_lines,omitempty"`
	Taxable            bool          `json:"taxable,omitempty"`
	Title              string        `json:"title,omitempty"`
	VariantId          int64         `json:"variant_id,omitempty"`
	VariantTitle       string        `json:"variant_title,omitempty"`
	Vendor             string        `json:"vendor,omitempty"`
}
