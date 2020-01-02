package caltrain

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type APIClient interface {
	Get(ctx context.Context, url string, query map[string]string) ([]byte, error)
}

type APIClient511 struct{}

func NewClient() *APIClient511 {
	return &APIClient511{}
}

// Get returns the body of the request
func (a *APIClient511) Get(ctx context.Context, url string, query map[string]string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// update the url with the required query parameters
	q := req.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// TODO: check status codes first?

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	return body, nil
}

type MockAPIClient struct {
	GetResult         []byte
	GetResultFilePath string
}

// Get returns either the value in a file or a defined byte array
func (a *MockAPIClient) Get(ctx context.Context, url string, query map[string]string) ([]byte, error) {
	if a.GetResultFilePath != "" {
		f, err := os.Open(a.GetResultFilePath)
		if err != nil {
			return a.GetResult, err
		}
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return a.GetResult, err
		}
		return data, nil
	} else {
		return a.GetResult, nil
	}
}
