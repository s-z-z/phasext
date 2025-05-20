package util

import (
	"fmt"
	"net/url"
)

func UrlEncode(rawURL string) (string, error) {
	// 解析 URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("解析 URL 失败: %v", err)
	}
	query := parsedURL.Query()
	encodedQuery := query.Encode()

	// 重新构建 URL
	parsedURL.RawQuery = encodedQuery
	encodedURL := parsedURL.String()
	return encodedURL, nil
}
