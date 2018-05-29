package shopify

type ClientDetail struct {
	AcceptLanguage string `json:"accept_language"`

	BrowserHeight int32 `json:"browser_height"`

	BrowserIp string `json:"browser_ip"`

	BrowserWidth int32 `json:"browser_width"`

	SessionHash string `json:"session_hash"`

	UserAgent string `json:"user_agent"`
}
