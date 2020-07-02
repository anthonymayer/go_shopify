package shopify

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Collect struct {
	CollectionId int64     `json:"collection_id"`
	CreatedAt    time.Time `json:"created_at"`
	Id           int64     `json:"id"`
	Position     int64     `json:"position"`
	ProductId    int64     `json:"product_id"`
	UpdatedAt    time.Time `json:"updated_at"`
	SortValue    string    `json:"sort_value"`

	api *API
}

type CollectOptions struct {
	Limit        int    `url:"limit,omitempty"`
	Page         int    `url:"page,omitempty"`
	CollectionID string `url:"collection_id,omitempty"`
	ProductID    string `url:"product_id,omitempty"`
}

func (api *API) NewCollect() *Collect {
	return &Collect{api: api}
}

func (api *API) Collects() ([]Collect, error) {
	return api.CollectsWithOptions(&CollectOptions{})
}

func (api *API) CollectsWithOptions(options *CollectOptions) ([]Collect, error) {
	qs := encodeOptions(options)
	endpoint := fmt.Sprintf("BASE_PATH/collects.json?%v", qs)
	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := &map[string][]Collect{}
	err = json.NewDecoder(res).Decode(r)
	result := (*r)["collects"]

	if err != nil {
		return nil, err
	}

	for _, v := range result {
		v.api = api
	}

	return result, nil
}

func (api *API) Collect(id int64) (*Collect, error) {
	endpoint := fmt.Sprintf("BASE_PATH/collects/%d.json", id)

	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := map[string]Collect{}
	err = json.NewDecoder(res).Decode(&r)
	result := r["collect"]

	if err != nil {
		return nil, err
	}

	result.api = api

	return &result, nil
}

type CollectsCountOptions struct {
	CollectionID string `url:"collection_id,omitempty"`
	ProductID    string `url:"product_id,omitempty"`
}

func (api *API) CollectsCount(options *CollectsCountOptions) (int, error) {

	qs := encodeOptions(options)
	endpoint := fmt.Sprintf("BASE_PATH/collects/count.json?%v", qs)

	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return 0, err
	}

	if status != 200 {
		return 0, fmt.Errorf("Status returned: %d", status)
	}

	r := map[string]interface{}{}
	err = json.NewDecoder(res).Decode(&r)

	result, _ := strconv.Atoi(fmt.Sprintf("%v", r["count"]))
	if err != nil {
		return 0, err
	}
	return result, nil
}
