package addrvrf

import (
	"net/http"
	"net/url"
)

type HTTPClient interface {
	Do(r *http.Request) (*http.Response, error)
}

type SmartyVerifier struct {
	HTTPClient HTTPClient
}

func NewSmartyVerifier(c HTTPClient) *SmartyVerifier {
	return &SmartyVerifier{
		HTTPClient: c,
	}
}

func (v SmartyVerifier) Verify(input AddressInput) {
	v.HTTPClient.Do(v.buildRequest(input))
}

func (v *SmartyVerifier) buildRequest(input AddressInput) *http.Request {
	r, _ := http.NewRequest(http.MethodGet, v.buildURL(input), nil)
	return r
}

func (v *SmartyVerifier) buildURL(input AddressInput) string {
	url, _ := url.Parse("https://us-street.api.smartystreets.com/street-address")
	q := url.Query()
	q.Set("street", input.Street)
	q.Set("city", input.City)
	q.Set("state", input.State)
	q.Set("zipcode", input.ZIPCode)
	url.RawQuery = q.Encode()
	return url.String()
}
