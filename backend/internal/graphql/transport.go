package graphql

import "net/http"

// BearerTransport modifies the request to include a access token
type BearerTransport struct {
	Token string
}

// RoundTrip implements the roundtripper interface adding a bearer token to the request.
func (t *BearerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+t.Token)

	return http.DefaultTransport.RoundTrip(req)
}

// KeyTransport adds a keyed header to the request
type KeyTransport struct {
	Key string
	Value string
}

// RoundTrip implements the roundtripper interface adding a key value to the request.
func (t *KeyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add(t.Key, t.Value)

	return http.DefaultTransport.RoundTrip(req)
}
