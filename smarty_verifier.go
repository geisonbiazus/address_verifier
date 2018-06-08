package addrvrf

import (
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	StatusSuccess         = "SUCCESS"
	StatusInvalidResponse = "INVALID_RESPONSE"
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

func (v SmartyVerifier) Verify(input AddressInput) AddressOutput {
	response, _ := v.HTTPClient.Do(v.buildRequest(input))
	return v.buildOutput(response)
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

func (v *SmartyVerifier) buildOutput(r *http.Response) AddressOutput {
	parsedResponse := []smartyAPIResponse{}

	if err := json.NewDecoder(r.Body).Decode(&parsedResponse); err != nil {
		return v.buildInvalidOutput()
	}

	return v.buildSuccessOutput(parsedResponse[0])
}

func (v *SmartyVerifier) buildInvalidOutput() AddressOutput {
	return AddressOutput{Status: StatusInvalidResponse}
}

func (v *SmartyVerifier) buildSuccessOutput(r smartyAPIResponse) AddressOutput {
	return AddressOutput{
		Status:        StatusSuccess,
		DeliveryLine1: r.DeliveryLine1,
		LastLine:      r.LastLine,
		Street:        r.Components.StreetName,
		City:          r.Components.CityName,
		State:         r.Components.StateAbbreviation,
		ZIPCode:       r.Components.ZIPCode,
	}
}

type smartyAPIResponse struct {
	DeliveryLine1 string `json:"delivery_line_1"`
	LastLine      string `json:"last_line"`
	Components    struct {
		StreetName        string `json:"street_name"`
		CityName          string `json:"city_name"`
		StateAbbreviation string `json:"state_abbreviation"`
		ZIPCode           string `json:"zipcode"`
	} `json:"components"`
}
