package shopify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type Collection struct {
	BodyHTML       string      `json:"body_html"`
	Disjunctive    bool        `json:"disjunctive"`
	Handle         string      `json:"handle"`
	ID             int64       `json:"id"`
	Image          interface{} `json:"image,omitempty"`
	ProductsCount  int         `json:"products_count"`
	PublishedAt    *time.Time  `json:"published_at,omitempty"`
	PublishedScope string      `json:"published_scope"`
	SortOrder      string      `json:"sort_order"`
	TemplateSuffix string      `json:"template_suffix"`
	Title          string      `json:"title"`
	UpdatedAt      *time.Time  `json:"updated_at,omitempty"`
	Rules          []Rule      `json:"rules"`

	api *API
}

type CollectionOptions struct {
	Handle          string `url:"handle,omitempty"`
	IDs             string `url:"ids,omitempty"`
	Limit           int    `url:"limit,omitempty"`
	Page            int    `url:"page,omitempty"`
	ProductID       string `url:"product_id,omitempty"`
	UpdatedAtMin    string `url:"updated_at_min,omitempty"`
	UpdatedAtMax    string `url:"updated_at_max,omitempty"`
	PublishedAtMin  string `url:"published_at_min,omitempty"`
	PublishedAtMax  string `url:"published_at_max,omitempty"`
	PublishedStatus string `url:"published_status,omitempty"`
	Title           string `url:"title,omitempty"`
}

func (api *API) Collections() ([]Collection, error) {
	return api.CollectionsWithOptions(&CollectionOptions{})
}

func (api *API) CollectionsWithOptions(options *CollectionOptions) ([]Collection, error) {
	qs := encodeOptions(options)
	endpoint := fmt.Sprintf("BASE_PATH/collection.json?%v", qs)
	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := &map[string][]Collection{}
	err = json.NewDecoder(res).Decode(r)

	result := (*r)["collection"]

	if err != nil {
		return nil, err
	}

	for _, v := range result {
		v.api = api
	}

	return result, nil
}

func (api *API) Collection(id int64) (*Collection, error) {
	endpoint := fmt.Sprintf("BASE_PATH/collections/%d.json", id)

	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := map[string]Collection{}
	err = json.NewDecoder(res).Decode(&r)
	result := r["collection"]

	if err != nil {
		return nil, err
	}

	result.api = api

	return &result, nil
}

func (api *API) CollectionProducts(collectionID int64, limit int) ([]*Product, *Pages, error) {
	endpoint := fmt.Sprintf("BASE_PATH/collections/%d/products.json?limit=%d", collectionID, limit)

	res, status, pages, err := api.requestWithPagination(endpoint, "GET", nil, nil)

	products, err := api.processProductsResponse(res, status, err)
	return products, pages, err
}

func (api *API) NewCollection() *Collection {
	return &Collection{api: api}
}

func (obj *Collection) Save() error {
	endpoint := fmt.Sprintf("BASE_PATH/collections/%d.json", obj.ID)
	method := "PUT"
	expectedStatus := 201

	if obj.ID == 0 {
		endpoint = fmt.Sprintf("BASE_PATH/collection.json")
		method = "POST"
		expectedStatus = 201
	}

	body := map[string]*Collection{}
	body["collection"] = obj

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
		}

		return fmt.Errorf("Status %d, and error parsing body: %s", status, err)
	}

	r := map[string]Collection{}
	err = json.NewDecoder(res).Decode(&r)

	if err != nil {
		return err
	}

	*obj = r["collection"]

	return nil
}
