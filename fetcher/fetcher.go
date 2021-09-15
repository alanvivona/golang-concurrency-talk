package fetcher

import (
	"errors"
	"io"
	"net/http"
	"net/url"
)

var ErrStatusNotSuccess = errors.New("Response status code != 200")

type Fetcher struct {
	client   http.Client
	uniq     map[string]bool
	maxDepth int
}

func (f *Fetcher) isUniq(url string) bool {
	exists := f.uniq[url]
	if !exists {
		f.uniq[url] = true
	}
	return !exists
}

func NewFetcher() *Fetcher {
	return &Fetcher{
		client: http.Client{},
		uniq:   map[string]bool{},
	}
}

func (f *Fetcher) Fetch(target url.URL) (io.ReadCloser, error) {
	if !f.isUniq(target.String()) {
		return nil, nil
	}
	resp, err := f.client.Get(target.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, ErrStatusNotSuccess
	}
	return resp.Body, nil
}
