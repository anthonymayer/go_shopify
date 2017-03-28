package shopify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type Article struct {
	Author         string       `json:"author"`
	BlogID         int64        `json:"blog_id"`
	BodyHTML       string       `json:"body_html"`
	CreatedAt      time.Time    `json:"created_at"`
	Handle         string       `json:"handle,omitempty"`
	ID             int64        `json:"id"`
	Image          articleImage `json:"image,omitempty"`
	PublishedAt    time.Time    `json:"published_at"`
	SummaryHTML    string       `json:"summary_html"`
	TemplateSuffix string       `json:"template_suffix"`
	Title          string       `json:"title"`
	UpdatedAt      time.Time    `json:"updated_at"`
	UserID         int64        `json:"user_id"`
	Tags           string       `json:"tags"`

	api *API
}

type ArticleOptions struct {
	Author          string `url:"author,omitempty"`
	Handle          string `url:"handle,omitempty"`
	Limit           int    `url:"limit,omitempty"`
	Page            int    `url:"page,omitempty"`
	CreatedAtMin    string `url:"created_at_min,omitempty"`
	CreatedAtMax    string `url:"created_at_max,omitempty"`
	UpdatedAtMin    string `url:"updated_at_min,omitempty"`
	UpdatedAtMax    string `url:"updated_at_max,omitempty"`
	PublishedAtMin  string `url:"published_at_min,omitempty"`
	PublishedAtMax  string `url:"published_at_max,omitempty"`
	PublishedStatus string `url:"published_status,omitempty"`
	Order           string `url:"order,omitempty"`
	SinceID         string `url:"since_id,omitempty"`
	Tag             string `url:"tag,omitempty"`
}

type articleImage struct {
	Src       string `json:"src,omitempty"`
	CreatedAt string `json:"created_at,omitempty`
}

func (api *API) Articles() ([]Article, error) {
	return api.ArticlesWithOptions(&ArticleOptions{})
}

func (api *API) ArticlesWithOptions(options *ArticleOptions) ([]Article, error) {
	qs := encodeOptions(options)
	endpoint := fmt.Sprintf("/admin/articles.json?%v", qs)
	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := &map[string][]Article{}
	err = json.NewDecoder(res).Decode(r)

	result := (*r)["articles"]

	if err != nil {
		return nil, err
	}

	for _, v := range result {
		v.api = api
	}

	return result, nil
}

func (api *API) BlogArticlesWithOptions(blogID int64, options *ArticleOptions) ([]Article, error) {
	qs := encodeOptions(options)
	endpoint := fmt.Sprintf("/admin/blogs/%d/articles.json?%v", blogID, qs)
	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := &map[string][]Article{}
	err = json.NewDecoder(res).Decode(r)

	result := (*r)["articles"]

	if err != nil {
		return nil, err
	}

	for _, v := range result {
		v.api = api
	}

	return result, nil
}

func (api *API) Article(id int64) (*Article, error) {
	endpoint := fmt.Sprintf("/admin/articles/%d.json", id)

	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := map[string]Article{}
	err = json.NewDecoder(res).Decode(&r)

	result := r["article"]

	if err != nil {
		return nil, err
	}

	result.api = api

	return &result, nil
}

func (api *API) NewArticle() *Article {
	return &Article{api: api}
}

func (obj *Article) Save() error {
	endpoint := fmt.Sprintf("/admin/articles/%d.json", obj.ID)
	method := "PUT"
	expectedStatus := 201

	if obj.ID == 0 {
		endpoint = fmt.Sprintf("/admin/articles.json")
		method = "POST"
		expectedStatus = 201
	}

	body := map[string]*Article{}
	body["article"] = obj

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

	r := map[string]Article{}
	err = json.NewDecoder(res).Decode(&r)

	if err != nil {
		return err
	}

	*obj = r["article"]

	return nil
}
