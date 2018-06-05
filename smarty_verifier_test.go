package addrvrf_test

import (
	"net/http"
	"testing"

	"github.com/geisonbiazus/addrvrf"
	"github.com/stretchr/testify/assert"
)

func TestSmartyVerifier(t *testing.T) {
	type fixture struct {
		client   *HTTPClientSpy
		verifier *addrvrf.SmartyVerifier
	}

	setup := func() *fixture {
		client := &HTTPClientSpy{}
		verifier := addrvrf.NewSmartyVerifier(client)

		return &fixture{
			client:   client,
			verifier: verifier,
		}
	}

	t.Run("Request is sent properly", func(t *testing.T) {
		f := setup()

		input := addrvrf.AddressInput{
			Street:  "street",
			City:    "city",
			State:   "state",
			ZIPCode: "zipcode",
		}

		f.verifier.Verify(input)

		assert.Equal(t, http.MethodGet, f.client.Request.Method)
		assert.Equal(t, "https", f.client.Request.URL.Scheme)
		assert.Equal(t, "us-street.api.smartystreets.com", f.client.Request.URL.Host)
		assert.Equal(t, "/street-address", f.client.Request.URL.Path)
		assert.Equal(t, "street", f.client.Request.URL.Query().Get("street"))
		assert.Equal(t, "city", f.client.Request.URL.Query().Get("city"))
		assert.Equal(t, "state", f.client.Request.URL.Query().Get("state"))
		assert.Equal(t, "zipcode", f.client.Request.URL.Query().Get("zipcode"))
	})
}

type HTTPClientSpy struct {
	Request *http.Request
}

func (c *HTTPClientSpy) Do(r *http.Request) (*http.Response, error) {
	c.Request = r
	return &http.Response{}, nil
}
