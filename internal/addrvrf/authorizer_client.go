package addrvrf

import "net/http"

type AuthorizerClient struct {
	inner     HTTPClient
	authID    string
	authToken string
}

func NewAuthorizerClient(inner HTTPClient, authID, authToken string) *AuthorizerClient {
	return &AuthorizerClient{
		inner:     inner,
		authID:    authID,
		authToken: authToken,
	}
}

func (c *AuthorizerClient) Do(r *http.Request) (*http.Response, error) {
	q := r.URL.Query()
	q.Set("auth-id", c.authID)
	q.Set("auth-token", c.authToken)
	r.URL.RawQuery = q.Encode()

	return c.inner.Do(r)
}
