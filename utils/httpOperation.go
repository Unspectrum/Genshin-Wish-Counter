package utils

import "net/url"

func GenerateGetParameter(params map[string]string) url.Values {
	queryParams := url.Values{}
	for key, val := range params {
		queryParams.Add(key, val)
	}
	return queryParams
}
