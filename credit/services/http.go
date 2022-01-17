package services

import (
	"bytes"
	"context"
	"credit/authorizer"
	"credit/settlement"
	"io/ioutil"
	"net/http"
)

type httpService struct{}

func (h *httpService) GetWithContext(ctx context.Context, url string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	bt, err := ioutil.ReadAll(resp.Body)

	return bt, resp.StatusCode, nil
}

func (h *httpService) PostWithContext(ctx context.Context, url string, payload []byte) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	bt, err := ioutil.ReadAll(resp.Body)

	return bt, resp.StatusCode, nil
}

func NewHttp() (authorizer.Http, settlement.Http) {
	return &httpService{}, &httpService{}
}
