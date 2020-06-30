package shopify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
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

func (api *API) Articles() ([]Article, *Pages, error) {
	return api.ArticlesWithOptions(&ArticleOptions{})
}

func (api *API) ArticlesWithOptions(options *ArticleOptions) ([]Article, *Pages, error) {
	return api.getArticlesWithOptions("BASE_PATH/articles.json", options)
}

func (api *API) BlogArticlesWithOptions(blogID int64, options *ArticleOptions) ([]Article, *Pages, error) {
	return api.getArticlesWithOptions(fmt.Sprintf("BASE_PATH/blogs/%d/articles.json", blogID), options)
}

func (api *API) getArticlesWithOptions(path string, options *ArticleOptions) ([]Article, *Pages, error) {
	qs := encodeOptions(options)
	endpoint := fmt.Sprintf("%v?%v", path, qs)

	return api.processArticlesResponse(api.requestWithPagination(endpoint, "GET", nil, nil))
}

func (api *API) ArticlesFromPages(pages *Pages) ([]Article, *Pages, error) {
	if pages.HasNextPage() {
		return api.processArticlesResponse(api.getNextPage(pages))
	}
	return nil, &Pages{}, fmt.Errorf("No next page")
}

func (api *API) processArticlesResponse(res *bytes.Buffer, status int, pages *Pages, err error) ([]Article, *Pages, error) {
	if err != nil {
		return nil, pages, err
	}

	if status != 200 {
		return nil, pages, fmt.Errorf("Status returned: %d", status)
	}

	r := &map[string][]Article{}
	err = json.NewDecoder(res).Decode(r)

	result := (*r)["articles"]

	if err != nil {
		return nil, pages, err
	}

	for _, v := range result {
		v.api = api
	}

	return result, pages, nil
}

type BlogArticlesCountOptions struct {
	CreatedAtMin    string `url:"created_at_min,omitempty"`
	CreatedAtMax    string `url:"created_at_max,omitempty"`
	UpdatedAtMin    string `url:"updated_at_min,omitempty"`
	UpdatedAtMax    string `url:"updated_at_max,omitempty"`
	PublishedAtMin  string `url:"published_at_min,omitempty"`
	PublishedAtMax  string `url:"published_at_max,omitempty"`
	PublishedStatus string `url:"published_status,omitempty"`
}

func (api *API) BlogArticlesCount(blogID int64, options *BlogArticlesCountOptions) (int, error) {

	qs := encodeOptions(options)
	endpoint := fmt.Sprintf("BASE_PATH/blogs/%d/articles/count.json?%v", blogID, qs)

	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return 0, err
	}

	if status != 200 {
		return 0, fmt.Errorf("Status returned: %d", status)
	}

	r := map[string]interface{}{}
	err = jsoniter.ConfigFastest.NewDecoder(res).Decode(&r)

	result, _ := strconv.Atoi(fmt.Sprintf("%v", r["count"]))
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (api *API) Article(id int64) (*Article, error) {
	endpoint := fmt.Sprintf("BASE_PATH/articles/%d.json", id)

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
	endpoint := fmt.Sprintf("BASE_PATH/articles/%d.json", obj.ID)
	method := "PUT"
	expectedStatus := 201

	if obj.ID == 0 {
		endpoint = fmt.Sprintf("BASE_PATH/articles.json")
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
