package shopify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type Order struct {
	BuyerAcceptsMarketing bool           `json:"buyer_accepts_marketing,omitempty"`
	CancelReason          string         `json:"cancel_reason,omitempty"`
	CancelledAt           string         `json:"cancelled_at,omitempty"`
	CartToken             string         `json:"cart_token,omitempty"`
	CheckoutToken         string         `json:"checkout_token,omitempty"`
	ClosedAt              string         `json:"closed_at,omitempty"`
	Confirmed             bool           `json:"confirmed,omitempty"`
	CreatedAt             time.Time      `json:"created_at,omitempty"`
	Currency              string         `json:"currency,omitempty"`
	Email                 string         `json:"email,omitempty"`
	FinancialStatus       string         `json:"financial_status,omitempty"`
	FulfillmentStatus     string         `json:"fulfillment_status,omitempty"`
	Gateway               string         `json:"gateway,omitempty"`
	Id                    int64          `json:"id,omitempty"`
	LandingSite           string         `json:"landing_site,omitempty"`
	LocationId            string         `json:"location_id,omitempty"`
	Name                  string         `json:"name,omitempty"`
	Note                  string         `json:"note,omitempty"`
	Number                int64          `json:"number,omitempty"`
	ProcessedAt           time.Time      `json:"processed_at,omitempty"`
	Reference             string         `json:"reference,omitempty"`
	ReferringSite         string         `json:"referring_site,omitempty"`
	SourceIdentifier      string         `json:"source_identifier,omitempty"`
	SourceName            string         `json:"source_name,omitempty"`
	SourceUrl             string         `json:"source_url,omitempty"`
	SubtotalPrice         string         `json:"subtotal_price,omitempty"`
	TaxesIncluded         bool           `json:"taxes_included,omitempty"`
	Test                  bool           `json:"test,omitempty"`
	Token                 string         `json:"token,omitempty"`
	TotalDiscounts        string         `json:"total_discounts,omitempty"`
	TotalLineItemsPrice   string         `json:"total_line_items_price,omitempty"`
	TotalPrice            string         `json:"total_price,omitempty"`
	TotalPriceUsd         string         `json:"total_price_usd,omitempty"`
	TotalTax              string         `json:"total_tax,omitempty"`
	TotalWeight           int64          `json:"total_weight,omitempty"`
	UpdatedAt             time.Time      `json:"updated_at,omitempty"`
	UserId                string         `json:"user_id,omitempty"`
	BrowserIp             string         `json:"browser_ip,omitempty"`
	LandingSiteRef        string         `json:"landing_site_ref,omitempty"`
	OrderNumber           int64          `json:"order_number,omitempty"`
	DiscountCodes         []interface{}  `json:"discount_codes,omitempty"`
	NoteAttributes        []interface{}  `json:"note_attributes,omitempty"`
	ProcessingMethod      string         `json:"processing_method,omitempty"`
	Source                string         `json:"source,omitempty"`
	CheckoutId            int64          `json:"checkout_id,omitempty"`
	TaxLines              []interface{}  `json:"tax_lines,omitempty"`
	Tags                  string         `json:"tags,omitempty"`
	LineItems             []LineItem     `json:"line_items,omitempty"`
	ShippingLines         []ShippingLine `json:"shipping_lines,omitempty"`
	BillingAddress        BillingAddress `json:"billing_address,omitempty"`
	ShippingAddress       BillingAddress `json:"shipping_address,omitempty"`
	Fulfillments          []interface{}  `json:"fulfillments,omitempty"`
	ClientDetails         ClientDetail   `json:"client_details,omitempty"`
	Refunds               []interface{}  `json:"refunds,omitempty"`
	Customer              Customer       `json:"customer,omitempty"`

	api *API
}

type OrderOptions struct {
	IDs          string `url:"ids,omitempty"`
	Limit        int    `url:"limit,omitempty"`
	Page         int    `url:"page,omitempty"`
	SinceID      int64  `url:"since_id,omitempty"`
	CollectionID string `url:"collection_id,omitempty"`
	CreatedAtMin string `url:"created_at_min,omitempty"`
	CreatedAtMax string `url:"created_at_max,omitempty"`
	UpdatedAtMin string `url:"updated_at_min,omitempty"`
	UpdatedAtMax string `url:"updated_at_max,omitempty"`
	Fields       string `url:"fields,omitempty"`
}

func (api *API) Orders() ([]Order, error) {
	return api.OrdersWithOptions(&OrderOptions{})
}

func (api *API) OrdersWithOptions(options *OrderOptions) ([]Order, error) {
	qs := encodeOptions(options)
	endpoint := fmt.Sprintf("/admin/orders.json?%v", qs)
	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := &map[string][]Order{}
	err = json.NewDecoder(res).Decode(r)

	result := (*r)["orders"]

	if err != nil {
		return nil, err
	}

	for _, v := range result {
		v.api = api
	}

	return result, nil
}

func (api *API) Order(id int64) (*Order, error) {
	endpoint := fmt.Sprintf("/admin/orders/%d.json", id)

	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := map[string]Order{}
	err = json.NewDecoder(res).Decode(&r)
	result := r["order"]

	if err != nil {
		return nil, err
	}

	result.api = api

	return &result, nil
}

func (api *API) NewOrder() *Order {
	return &Order{api: api}
}

func (obj *Order) Save() error {
	endpoint := fmt.Sprintf("/admin/orders/%d.json", obj.Id)
	method := "PUT"
	expectedStatus := 201

	if obj.Id == 0 {
		endpoint = fmt.Sprintf("/admin/orders.json")
		method = "POST"
		expectedStatus = 201
	}

	body := map[string]*Order{}
	body["order"] = obj

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(body)

	if err != nil {
		return err
	}

	res, status, err := obj.api.request(endpoint, method, nil, buf)

	if err != nil {
		return err
	}

	if status != expectedStatus {
		r := errorResponse{}
		err = json.NewDecoder(res).Decode(&r)
		if err == nil {
			return fmt.Errorf("Status %d: %v", status, r.Errors)
		} else {
			return fmt.Errorf("Status %d, and error parsing body: %s", status, err)
		}
	}

	r := map[string]Order{}
	err = json.NewDecoder(res).Decode(&r)

	if err != nil {
		return err
	}

	*obj = r["order"]

	return nil
}
