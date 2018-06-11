package addrvrf_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geisonbiazus/addrvrf"
	"github.com/stretchr/testify/assert"
)

func TestAuthorizerClient(t *testing.T) {
	type fixture struct {
		authID    string
		authToken string
		inner     *HTTPClientSpy
		client    *addrvrf.AuthorizerClient
		request   *http.Request
	}

	setup := func() *fixture {
		authID := "auth ID"
		authToken := "auth token"
		inner := &HTTPClientSpy{}
		client := addrvrf.NewAuthorizerClient(inner, authID, authToken)
		request := httptest.NewRequest(http.MethodGet, "/?some=param", nil)

		return &fixture{
			authID:    authID,
			authToken: authToken,
			inner:     inner,
			client:    client,
			request:   request,
		}
	}

	t.Run("Auth ID and Token are inserted in the request query", func(t *testing.T) {
		f := setup()
		f.client.Do(f.request)
		assert.Equal(t, f.authID, f.request.URL.Query().Get("auth-id"))
		assert.Equal(t, f.authToken, f.request.URL.Query().Get("auth-token"))
		assert.Equal(t, "param", f.request.URL.Query().Get("some"))
	})

	t.Run("Request is sent to inner client and response is returned", func(t *testing.T) {
		f := setup()
		f.inner.Configure(http.StatusOK, "body")
		f.inner.ConfigureError(errors.New("error"))

		resp, err := f.client.Do(f.request)

		assert.Equal(t, f.request, f.inner.Request)
		assert.Equal(t, f.inner.Response, resp)
		assert.Equal(t, f.inner.Error, err)
	})
}
