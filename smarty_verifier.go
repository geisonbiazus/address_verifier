package addrvrf

import (
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	StatusDeliverable = "Deliverable"
	StatusVacant      = "Vacant"
	StatusInactive    = "Inactive"
	StatusInvalid     = "Invalid"
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
	return v.decodeAndBuildOutput(response)
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

func (v *SmartyVerifier) decodeAndBuildOutput(r *http.Response) AddressOutput {
	parsedResponse := []smartyAPIResponse{}

	if err := json.NewDecoder(r.Body).Decode(&parsedResponse); err != nil {
		return v.buildInvalidOutput()
	}

	return v.buildOutput(parsedResponse[0])
}

func (v *SmartyVerifier) buildInvalidOutput() AddressOutput {
	return AddressOutput{Status: StatusInvalid}
}

func (v *SmartyVerifier) buildOutput(r smartyAPIResponse) AddressOutput {
	return AddressOutput{
		Status:        v.computeStatus(r),
		DeliveryLine1: r.DeliveryLine1,
		LastLine:      r.LastLine,
		Street:        r.Components.StreetName,
		City:          r.Components.CityName,
		State:         r.Components.StateAbbreviation,
		ZIPCode:       r.Components.ZIPCode,
	}
}

func (v *SmartyVerifier) computeStatus(r smartyAPIResponse) string {
	if r.Analysis.Match == "Y" || r.Analysis.Match == "D" || r.Analysis.Match == "S" {
		if r.Analysis.Active != "Y" {
			return StatusInactive
		}

		if r.Analysis.Vacant == "Y" {
			return StatusVacant
		}

		return StatusDeliverable
	}

	return StatusInvalid
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
	Analysis struct {
		Match  string `json:"dpv_match_code"`
		Vacant string `json:"dpv_vacant"`
		Active string `json:"active"`
	} `json:"analysis"`
}
