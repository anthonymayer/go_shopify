package shopify

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Variant struct {
	Barcode              string      `json:"barcode,omitempty"`
	CompareAtPrice       string      `json:"compare_at_price,omitempty"`
	CreatedAt            string      `json:"created_at,omitempty"`
	FulfillmentService   string      `json:"fulfillment_service,omitempty"`
	Grams                float64     `json:"grams,omitempty"`
	Weight               float64     `json:"weight,omitempty"`
	WeightUnit           string      `json:"weight_unit,omitempty"`
	ID                   int64       `json:"id,omitempty"`
	InventoryManagement  string      `json:"inventory_management,omitempty"`
	InventoryPolicy      string      `json:"inventory_policy,omitempty"`
	InventoryQuantity    int64       `json:"inventory_quantity"`
	InventoryItemID      int64       `json:"inventory_item_id"`
	OldInventoryQuantity int64       `json:"old_inventory_quantity,omitempty"`
	Metafield            interface{} `json:"metafield,omitempty"`
	Option1              *string     `json:"option1,omitempty"`
	Option2              *string     `json:"option2,omitempty"`
	Option3              *string     `json:"option3,omitempty"`
	Position             int64       `json:"position,omitempty"`
	Price                string      `json:"price,omitempty"`
	ProductID            int64       `json:"product_id,omitempty"`
	RequiresShipping     bool        `json:"requires_shipping,omitempty"`
	SKU                  *string     `json:"sku,omitempty"`
	Taxable              bool        `json:"taxable,omitempty"`
	Title                *string     `json:"title,omitempty"`
	UpdatedAt            string      `json:"updated_at,omitempty"`
	ImageID              int64       `json:"image_id,omitempty"`

	api *API
}

func (api *API) NewVariant() *Variant {
	return &Variant{api: api}
}

func (obj *Variant) Save() error {
	endpoint := fmt.Sprintf("/admin/variants/%d.json", obj.ID)
	method := "PUT"
	expectedStatus := 200

	var buf bytes.Buffer
	body := map[string]*Variant{
		"variant": obj,
	}
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return err
	}
	reqBody := buf.Bytes()

	res, status, err := obj.api.request(endpoint, method, nil, &buf)
	if err != nil {
		return err
	}

	if status != expectedStatus {
		return newErrorResponse(status, reqBody, res)
	}

	r := map[string]Variant{}
	err = json.NewDecoder(res).Decode(&r)
	if err != nil {
		return err
	}

	api := obj.api
	*obj = r["variant"]
	obj.api = api

	return nil
}

func (obj *Variant) Metafields(options *MetafieldsOptions) ([]*Metafield, error) {
	if obj == nil || obj.api == nil {
		return nil, errors.New("Variant is nil")
	}
	qs := encodeOptions(options)
	endpoint := fmt.Sprintf("/admin/variants/%d/metafields.json?%v", obj.ID, qs)
	res, status, err := obj.api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := map[string][]*Metafield{}
	err = json.NewDecoder(res).Decode(&r)

	result := r["metafields"]

	if err != nil {
		return nil, err
	}

	for _, v := range result {
		v.api = obj.api
	}

	return result, nil

}
