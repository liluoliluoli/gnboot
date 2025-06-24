package httpclient_util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	errors "github.com/go-kratos/kratos/v2/errors"
	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/liluoliluoli/gnboot/internal/common/utils/json_util"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
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

func DoPost[R, T any](ctx context.Context, url string, body *R, headerMap map[string]string) (*T, error) {
	httpClient := GetHttpClient()
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	for key, value := range headerMap {
		headers.Set(key, value)
	}
	var bodyBytes []byte
	var err error
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	response, err := httpClient.Post(url, bytes.NewBuffer(bodyBytes), headers)
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

func DoHtml(ctx context.Context, url string) (string, error) {
	httpClient := GetHttpClient()
	headers := http.Header{}
	//headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	//headers.Set("Referer", "https://www.douban.com/")
	//headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
	response, err := httpClient.Get(url, headers)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	rBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	if response.StatusCode >= 401 && response.StatusCode <= 403 {
		return "", errors.New(response.StatusCode, "token expired", "")
	}
	if response.StatusCode != 200 {
		return "", errors.New(response.StatusCode, "request failed", "")
	}
	if err != nil {
		return "", err
	}
	return string(rBody), nil
}

func CheckImageUrl(url string) (bool, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return false, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	// 只读取前512字节用于验证
	buf := make([]byte, 512)
	_, err = io.ReadFull(resp.Body, buf)
	if err != nil {
		return false, fmt.Errorf("读取图片失败: %v", err)
	}

	// 检测内容类型
	contentType := http.DetectContentType(buf)
	if !strings.HasPrefix(contentType, "image/") {
		return false, fmt.Errorf("无效的图片格式: %s", contentType)
	}

	return true, nil
}
