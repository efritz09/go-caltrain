package caltrain

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	// APILimitError is returned on a failed API request when the failure
	// reason is that the number of requests has exceeded the rate limit
	APILimitError = &apiLimitError{}

	// APIError is returned on a failed API request for any reason other
	// than too many requests
	APIError = &apiError{}
)

type apiLimitError struct{}

func (a *apiLimitError) Error() string {
	return "API call limit to 511.org has been reached"
}

type apiError struct {
	status string
}

func (a *apiError) Error() string {
	return fmt.Sprintf("API error: %s", a.status)
}

// APIClient is an interface for making requests
type APIClient interface {
	Get(ctx context.Context, url string, query map[string]string) ([]byte, error)
}

// APIClient511 implements APIClient with 511.org requests
type APIClient511 struct{}

// NewClient returns an instance of the APIClient511 struct
func NewClient() *APIClient511 {
	return &APIClient511{}
}

// Get makes a GET request to the 511.org API and returns the request body
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

	if resp.StatusCode != http.StatusOK {
		// return a specific error for too many requests
		if resp.StatusCode == http.StatusTooManyRequests {
			return nil, &apiLimitError{}
		} else {
			return nil, &apiError{status: resp.Status}
		}
	}

	// TODO: return the number of tries left? It exists in the header under
	// the Ratelimit-Limit and Ratelimit-Remaining keys. The api appears to
	// have be volatile on how many remaining calls can be made
	//
	// remain, err := strconv.Atoi(resp.Header["Ratelimit-Remaining"][0])

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	return body, nil
}

type apiClientMock struct {
	GetResult         []byte
	GetResultFilePath string
}

// Get returns either the value in a file or a defined byte array
func (a *apiClientMock) Get(ctx context.Context, url string, query map[string]string) ([]byte, error) {
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
