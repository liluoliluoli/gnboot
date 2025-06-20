package httpclient_util

import (
	"bytes"
	"context"
	"encoding/json"
	errors "github.com/go-kratos/kratos/v2/errors"
	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/liluoliluoli/gnboot/internal/common/utils/json_util"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var (
	instance *httpclient.Client
	once     sync.Once
)

func GetHttpClient() *httpclient.Client {
	once.Do(func() {
		timeout := 10000 * time.Millisecond
		instance = httpclient.NewClient(
			httpclient.WithHTTPTimeout(timeout),
			httpclient.WithRetryCount(2),
			httpclient.WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(10*time.Millisecond, 50*time.Millisecond))),
		)
	})
	return instance
}

func DoPost[R, T any](ctx context.Context, url string, token string, body *R) (*T, error) {
	httpClient := GetHttpClient()
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Add("Authorization", token)
	marshalString, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	response, err := httpClient.Post(url, bytes.NewBuffer(marshalString), headers)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	rBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= 401 && response.StatusCode <= 403 {
		return nil, errors.New(response.StatusCode, "token expired", "")
	}
	if response.StatusCode != 200 {
		return nil, errors.New(response.StatusCode, "request failed", "")
	}
	unmarshal, err := json_util.Unmarshal[*T](rBody)
	if err != nil {
		return nil, err
	}
	return unmarshal, nil
}

func DoGet[T any](ctx context.Context, url string, headerMap map[string]string) (*T, error) {
	httpClient := GetHttpClient()
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	for key, value := range headerMap {
		headers.Set(key, value)
	}
	response, err := httpClient.Get(url, headers)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	rBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= 401 && response.StatusCode <= 403 {
		return nil, errors.New(response.StatusCode, "token expired", "")
	}
	if response.StatusCode != 200 {
		return nil, errors.New(response.StatusCode, "request failed", "")
	}
	unmarshal, err := json_util.Unmarshal[*T](rBody)
	if err != nil {
		return nil, err
	}
	return unmarshal, nil
}
