package remote_test

import "net/http"

const (
	MockCustomer = `{
		"id": 3,
		"name": "Customer Name",
		"cpf": "12345678910",
		"email": "email@gmail.com"
	}`
)

type MockRoundTripper struct {
	Response *http.Response
}

func (trip *MockRoundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	return trip.Response, nil
}
