package addrvrf_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/geisonbiazus/addrvrf"
	"github.com/stretchr/testify/assert"
)

func TestSmartyVerifier(t *testing.T) {
	type fixture struct {
		client   *HTTPClientSpy
		verifier *addrvrf.SmartyVerifier
		input    addrvrf.AddressInput
	}

	setup := func() *fixture {
		client := &HTTPClientSpy{}
		verifier := addrvrf.NewSmartyVerifier(client)
		input := addrvrf.AddressInput{
			Street:  "street",
			City:    "city",
			State:   "state",
			ZIPCode: "zipcode",
		}

		client.Configure(http.StatusOK, SmartyAPISuccessResponse)

		return &fixture{
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
			Status:        addrvrf.StatusSuccess,
			DeliveryLine1: "delivery line 1",
			LastLine:      "last line",
			Street:        "street",
			City:          "city",
			State:         "state",
			ZIPCode:       "zip code",
		}

		assert.Equal(t, expected, output)
	})

	t.Run("Invalid JSON is returned", func(t *testing.T) {
		f := setup()
		f.client.Configure(http.StatusOK, SmartyAPIInvalidJSONResponse)

		output := f.verifier.Verify(f.input)

		expected := addrvrf.AddressOutput{
			Status: addrvrf.StatusInvalidResponse,
		}

		assert.Equal(t, expected, output)
	})
}

const SmartyAPISuccessResponse = `
[
	{
    "delivery_line_1": "delivery line 1",
    "last_line": "last line",
    "components": {
      "street_name": "street",
      "city_name": "city",
      "state_abbreviation": "state",
      "zipcode": "zip code"
    }
  }
]
`

const SmartyAPIInvalidJSONResponse = "Invalid JSON"

type HTTPClientSpy struct {
	Request  *http.Request
	Response *http.Response
}

func (c *HTTPClientSpy) Do(r *http.Request) (*http.Response, error) {
	c.Request = r
	return c.Response, nil
}

func (c *HTTPClientSpy) Configure(status int, body string) {
	c.Response = &http.Response{
		StatusCode: status,
		Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
	}
}
