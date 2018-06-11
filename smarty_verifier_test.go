package addrvrf_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/geisonbiazus/addrvrf"
	"github.com/stretchr/testify/assert"
)

type smartyVerifierFixture struct {
	client   *HTTPClientSpy
	verifier *addrvrf.SmartyVerifier
	input    addrvrf.AddressInput
}

func TestSmartyVerifier(t *testing.T) {

	setup := func() *smartyVerifierFixture {
		client := &HTTPClientSpy{}
		verifier := addrvrf.NewSmartyVerifier(client)
		input := addrvrf.AddressInput{
			Street:  "street",
			City:    "city",
			State:   "state",
			ZIPCode: "zipcode",
		}

		client.Configure(http.StatusOK, SmartyAPISuccessResponse)

		return &smartyVerifierFixture{
			client:   client,
			verifier: verifier,
			input:    input,
		}
	}

	t.Run("Request is sent properly", func(t *testing.T) {
		f := setup()

		f.verifier.Verify(f.input)

		assert.Equal(t, http.MethodGet, f.client.Request.Method)
		assert.Equal(t, "https", f.client.Request.URL.Scheme)
		assert.Equal(t, "us-street.api.smartystreets.com", f.client.Request.URL.Host)
		assert.Equal(t, "/street-address", f.client.Request.URL.Path)
		assert.Equal(t, "street", f.client.Request.URL.Query().Get("street"))
		assert.Equal(t, "city", f.client.Request.URL.Query().Get("city"))
		assert.Equal(t, "state", f.client.Request.URL.Query().Get("state"))
		assert.Equal(t, "zipcode", f.client.Request.URL.Query().Get("zipcode"))
	})

	t.Run("Response is parsed and output is returned", func(t *testing.T) {
		f := setup()
		f.client.Configure(http.StatusOK, SmartyAPISuccessResponse)

		output := f.verifier.Verify(f.input)

		expected := addrvrf.AddressOutput{
			Status:        addrvrf.StatusDeliverable,
			DeliveryLine1: "delivery line 1",
			LastLine:      "last line",
			Street:        "street",
			City:          "city",
			State:         "state",
			ZIPCode:       "zip code",
		}

		assert.Equal(t, expected, output)
		assert.True(t, f.client.ReadCloser.Closed)
	})

	t.Run("Invalid JSON is returned", func(t *testing.T) {
		f := setup()
		f.client.Configure(http.StatusOK, SmartyAPIInvalidJSONResponse)

		output := f.verifier.Verify(f.input)

		expected := addrvrf.AddressOutput{
			Status: addrvrf.StatusInvalid,
		}

		assert.Equal(t, expected, output)
	})

	t.Run("Status is computed correctly", func(t *testing.T) {
		f := setup()

		var (
			deliverable      = buildSmartyAPIResponse("Y", "N", "Y")
			missingSecondary = buildSmartyAPIResponse("D", "N", "Y")
			droppedSecondary = buildSmartyAPIResponse("S", "N", "Y")
			vacant           = buildSmartyAPIResponse("Y", "Y", "Y")
			inactive         = buildSmartyAPIResponse("Y", "N", "?")
			invalid          = buildSmartyAPIResponse("N", "?", "?")
		)

		validateAndAssertStatus(t, f, deliverable, addrvrf.StatusDeliverable)
		validateAndAssertStatus(t, f, missingSecondary, addrvrf.StatusDeliverable)
		validateAndAssertStatus(t, f, droppedSecondary, addrvrf.StatusDeliverable)
		validateAndAssertStatus(t, f, vacant, addrvrf.StatusVacant)
		validateAndAssertStatus(t, f, inactive, addrvrf.StatusInactive)
		validateAndAssertStatus(t, f, invalid, addrvrf.StatusInvalid)
	})
}

var (
	SmartyAPISuccessResponse     = buildSmartyAPIResponse("Y", "N", "Y")
	SmartyAPIInvalidJSONResponse = "Invalid JSON"
)

func buildSmartyAPIResponse(match, vacant, active string) string {
	return fmt.Sprintf(`
	[
		{
	    "delivery_line_1": "delivery line 1",
	    "last_line": "last line",
	    "components": {
	      "street_name": "street",
	      "city_name": "city",
	      "state_abbreviation": "state",
	      "zipcode": "zip code"
	    },
			"analysis": {
				"dpv_match_code": "%s",
				"dpv_vacant": "%s",
				"active": "%s"
			}
	  }
	]
	`, match, vacant, active)
}

func validateAndAssertStatus(t *testing.T, f *smartyVerifierFixture, response, status string) {
	f.client.Configure(http.StatusOK, response)
	output := f.verifier.Verify(f.input)
	assert.Equal(t, status, output.Status)
}
