package clients_test

import (
	"bytes"
	"github.com/matthewrobinsondev/lazyjira/pkg/clients"
	"github.com/matthewrobinsondev/lazyjira/pkg/config"
	"io"
	"net/http"
	"net/url"
	"testing"
)

// Mock Http Client
type fakeClient func(*http.Request) (*http.Response, error)

func (s fakeClient) RoundTrip(req *http.Request) (*http.Response, error) {
	return s(req)
}

func TestCanMakeRequestAndReturnResponse(t *testing.T) {
	fakeResponse := `{
		"startAt": 0,
		"maxResults": 50,
		"total": 0,
		"issues": []
	}`

	fakeHttpClient := &http.Client{
		Transport: fakeClient(func(request *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: io.NopCloser(bytes.NewBuffer([]byte(fakeResponse))),
			}, nil
		}),
	}

	fakeConfig := &config.Config{}

	jiraClient := clients.NewJiraClient(
		fakeConfig,
		fakeHttpClient,
	)

	response, err := jiraClient.NewRequest("GET", "test", url.Values{}, nil, clients.VERSION_3)
	if err != nil {
		t.Error(err)
	}

	if string(response) != fakeResponse {
		t.Error("Client did not return expected response")
	}

}
