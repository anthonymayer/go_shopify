package shopify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type Blog struct {
	Commentable        string    `json:"commentable"`
	CreatedAt          time.Time `json:"created_at"`
	Feedburner         string    `json:"feedburner"`
	FeedburnerLocation string    `json:"feedburner_location"`
	Handle             string    `json:"handle"`
	ID                 int64     `json:"id"`
	TemplateSuffix     string    `json:"template_suffix"`
	Title              string    `json:"title"`
	UpdatedAt          time.Time `json:"updated_at"`
	Tags               string    `json:"tags"`

	api *API
}

type BlogOptions struct {
	Handle  string `url:"handle,omitempty"`
	Limit   int    `url:"limit,omitempty"`
	Page    int    `url:"page,omitempty"`
	SinceID string `url:"since_id,omitempty"`
}

func (api *API) Blogs() ([]Blog, error) {
	return api.BlogsWithOptions(&BlogOptions{})
}

func (api *API) BlogsWithOptions(options *BlogOptions) ([]Blog, error) {
	qs := encodeOptions(options)
	endpoint := fmt.Sprintf("BASE_PATH/blogs.json?%v", qs)
	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := &map[string][]Blog{}
	err = json.NewDecoder(res).Decode(r)

	result := (*r)["blogs"]

	if err != nil {
		return nil, err
	}

	for _, v := range result {
		v.api = api
	}

	return result, nil
}

func (api *API) Blog(id int64) (*Blog, error) {
	endpoint := fmt.Sprintf("BASE_PATH/blogs/%d.json", id)

	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := map[string]Blog{}
	err = json.NewDecoder(res).Decode(&r)

	result := r["blog"]

	if err != nil {
		return nil, err
	}

	result.api = api

	return &result, nil
}

func (api *API) NewBlog() *Blog {
	return &Blog{api: api}
}

func (obj *Blog) Save() error {
	endpoint := fmt.Sprintf("BASE_PATH/blogs/%d.json", obj.ID)
	method := "PUT"
	expectedStatus := 201

	if obj.ID == 0 {
		endpoint = fmt.Sprintf("BASE_PATH/blogs.json")
		method = "POST"
		expectedStatus = 201
	}

	body := map[string]*Blog{}
	body["blog"] = obj

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

	r := map[string]Blog{}
	err = json.NewDecoder(res).Decode(&r)

	if err != nil {
		return err
	}

	*obj = r["blog"]

	return nil
}
