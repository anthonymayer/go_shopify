package shopify

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jpillora/backoff"
)

const REFILL_RATE = float64(0.5) // 2 per second
const BUCKET_LIMIT = 40
const MAX_RETRIES = 3
const BASE_PATH = "/admin/api/%v"

type API struct {
	Shop         string // for e.g. demo-3.myshopify.com
	AccessToken  string // permanent store access token
	Token        string // API client token
	Secret       string // API client secret for this shop
	RetryLimit   int
	LogRetryFail bool
	APIVersion   string

	// map[endpoint:method]response
	RequestCache RequestCache

	client *http.Client

	callLimit  int
	callsMade  int
	backoff    *backoff.Backoff
	retryCount int
}

type errorResponse struct {
	Errors map[string]interface{} `json:"errors"`
}

type RequestCache interface {
	Contains(string) bool
	Get(string) *bytes.Buffer
	Set(string, *bytes.Buffer)
}

type Pages struct {
	prevPage string
	nextPage string
}

func (pages *Pages) HasNextPage() bool {
	return pages.nextPage != ""
}

func (pages *Pages) HasPrevPage() bool {
	return pages.prevPage != ""
}

func (api *API) getNextPage(pages *Pages) (result *bytes.Buffer, status int, p *Pages, err error) {
	return api.baseRequest(pages.nextPage, "GET", nil, nil)
}

func NewPages(linkHeader string) (pages *Pages) {
	pages = &Pages{}
	if linkHeader == "" {
		return
	}
	rels := strings.Split(linkHeader, ",")
	for _, rel := range rels {
		pieces := strings.Split(rel, ";")
		link := strings.Trim(pieces[0], " <>")
		if strings.TrimSpace(pieces[1]) == "rel=\"next\"" {
			pages.nextPage = link
		} else {
			pages.prevPage = link
		}
	}
	return pages
}

func (api *API) request(endpoint string, method string, params map[string]interface{}, body io.Reader) (result *bytes.Buffer, status int, err error) {
	result, status, _, err = api.requestWithPagination(endpoint, method, params, body)
	return
}

func (api *API) requestWithPagination(endpoint string, method string, params map[string]interface{}, body io.Reader) (result *bytes.Buffer, status int, pages *Pages, err error) {
	return api.baseRequest(fmt.Sprintf("https://%s%s", api.Shop, endpoint), method, params, body)
}

func (api *API) baseRequest(uri string, method string, params map[string]interface{}, body io.Reader) (result *bytes.Buffer, status int, pages *Pages, err error) {
	if api.RequestCache != nil && api.RequestCache.Contains(uri+":"+method) {
		// make a copy so that the original object doesn't get emptied
		cachedBuffer := api.RequestCache.Get(uri + ":" + method)
		if cachedBuffer.Len() > 0 {
			return bytes.NewBuffer(cachedBuffer.Bytes()), 200, &Pages{}, nil
		}
	}
	if api.client == nil {
		api.client = &http.Client{}
	}
	if api.backoff == nil {
		api.backoff = &backoff.Backoff{
			//These are the defaults
			Min:    100 * time.Millisecond,
			Max:    2 * time.Second,
			Jitter: true,
		}
	}
	if api.callLimit == 0 {
		api.callLimit = BUCKET_LIMIT
	}

	switch api.APIVersion {
	case "":
		api.APIVersion = "2021-07" // current stable release
	case "2021-07", "2021-04", "2021-01", "2020-10": // valid do nothing
	default:
		fmt.Println("unknown API version")
	}
	uri = strings.Replace(uri, "BASE_PATH", fmt.Sprintf(BASE_PATH, api.APIVersion), 1)

	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return
	}

	if api.AccessToken != "" {
		req.Header.Set("X-Shopify-Access-Token", api.AccessToken)
	} else {
		req.SetBasicAuth(api.Token, api.Secret)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := api.client.Do(req)
	if err != nil {
		return
	}

	pages = NewPages(resp.Header.Get("Link"))

	calls, total := parseAPICallLimit(resp.Header.Get("HTTP_X_SHOPIFY_SHOP_API_CALL_LIMIT"))
	api.callsMade = calls
	api.callLimit = total

	status = resp.StatusCode
	if status == 429 { // statusTooManyRequests
		if api.RetryLimit == 0 {
			api.RetryLimit = MAX_RETRIES
		}
		if api.retryCount < api.RetryLimit {
			api.retryCount = api.retryCount + 1
			b := api.backoff.Duration()
			time.Sleep(b)
			// try again
			return api.baseRequest(uri, method, params, body)
		}
		if api.LogRetryFail {
			fmt.Println(
				"shopify api retry failed",
				"shop:", api.Shop,
				"calls made:", calls,
				"call limit:", total,
			)
		}
		// else just return
	}

	result = &bytes.Buffer{}
	defer resp.Body.Close()
	if _, err = io.Copy(result, resp.Body); err != nil {
		return
	}
	if result.Len() > 0 &&
		status == 200 &&
		params == nil &&
		body == nil &&
		api.RequestCache != nil {
		api.RequestCache.Set(uri+":"+method, bytes.NewBuffer(result.Bytes()))
	}
	return
}

func parseAPICallLimit(str string) (int, int) {
	tokens := strings.Split(str, "/")
	if len(tokens) != 2 {
		return 0, 0
	}
	calls, err := strconv.Atoi(tokens[0])
	if err != nil {
		return 0, 0
	}
	total, err := strconv.Atoi(tokens[1])
	if err != nil {
		return 0, 0
	}
	return calls, total
}
